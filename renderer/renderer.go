package renderer

import (
	"errors"
	"image/color"
	"runtime"
	"time"

	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	im "github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
)

// Renderer is a central part of the raytracer responsible for controlling the
// process reporting back progress to the caller
type Renderer struct {
	scene                          *Scene
	lights                         *hittable.HittableList
	output                         chan<- RenderProgress
	abort                          <-chan bool
	albedoShader                   AlbedoShader
	normalShader                   NormalShader
	maxMillisBetweenProgressOutput int64
}

// NewRenderer creates a new renderer given a scene and channels for communicating with the caller
func NewRenderer(scene *Scene, output chan<- RenderProgress, abort <-chan bool) (*Renderer, error) {

	lights := hittable.NewHittableList()
	findLights(scene.World, &lights)

	if len(lights.List()) == 0 {
		return nil, errors.New("Scene should have at least one light")
	}

	return &Renderer{
		scene:                          scene,
		lights:                         &lights,
		output:                         output,
		abort:                          abort,
		albedoShader:                   AlbedoShader{},
		normalShader:                   NormalShader{},
		maxMillisBetweenProgressOutput: 500,
	}, nil
}

func (r *Renderer) rayColor(ray geo.Ray, depth int) (geo.Vec3, geo.Vec3, geo.Vec3) {

	s := r.scene

	hit, rec := s.World.Hit(ray, util.Interval{Min: 0.001, Max: util.Infinity})
	if hit {
		pixelColor := s.RenderConfig.Shader.Shade(r, rec, ray, depth)

		var albedoColor geo.Vec3
		var normalColor geo.Vec3

		if r.scene.RenderConfig.PostProcessor != nil && depth == 0 {
			albedoColor = r.albedoShader.Shade(r, rec, ray, depth)
			normalColor = r.normalShader.Shade(r, rec, ray, depth)
		}

		return pixelColor, albedoColor, normalColor
	}

	return s.BackgroundColor, s.BackgroundColor, geo.ZeroVector
}

// Render executes the rendering of the image
func (r *Renderer) Render(imageWidth, imageHeight int) {

	s := r.scene
	samplesPerPixel := s.RenderConfig.SamplesPerPixel
	postProcessor := s.RenderConfig.PostProcessor
	pixelCount := imageWidth * imageHeight

	pixelColors := make([]geo.Vec3, pixelCount)
	albedoColors := make([]geo.Vec3, pixelCount)
	normalColors := make([]geo.Vec3, pixelCount)

	workerJobChannel := make(chan int, imageHeight)
	workerDoneChannel := make(chan bool)
	aborted := false

	camera := camera.New(imageWidth, imageHeight, r.scene.Camera)

	// Setup the pool of worker goroutines responsible for rendering lines

	numWorkers := numWorkers()
	for i := 0; i < numWorkers; i++ {
		go func() {

			for y := range workerJobChannel {

				for x := 0; x < imageWidth; x++ {
					if aborted {
						break
					}

					i := (((imageHeight-1)-y)*imageWidth + x)

					u := (float64(x) + random.RandomNormalFloat()) / float64(imageWidth-1)
					v := (float64(y) + random.RandomNormalFloat()) / float64(imageHeight-1)
					ray := camera.GetRay(u, v)
					pixelColor, albedoColor, normalColor := r.rayColor(ray, 0)

					pixelColors[i] = pixelColors[i].Add(pixelColor)

					if postProcessor != nil {
						albedoColors[i] = albedoColors[i].Add(albedoColor)
						normalColors[i] = normalColors[i].Add(normalColor)
					}

				}
				workerDoneChannel <- true
			}
		}()
	}

	var lastProgressTime time.Time

	// Render loop, executes the workers and reports progress
RenderLoop:
	for sample := 1; sample <= samplesPerPixel; sample++ {

		lastProgressTime = time.Now()

		// Submit jobs to the workers

		for y := imageHeight - 1; y >= 0; y-- {
			workerJobChannel <- y
		}
		for y := 0; y < imageHeight; y++ {

			<-workerDoneChannel

			// Does caller want to abort the render?

			select {
			case <-r.abort:
				aborted = true
				break RenderLoop
			default:
			}

			// If it was sufficiently long ago we last reported progress, do it
			// even though a sample is not complete

			nowTime := time.Now()
			millisSinceLastProgress := nowTime.Sub(lastProgressTime).Milliseconds()
			if millisSinceLastProgress > r.maxMillisBetweenProgressOutput && !aborted {
				lastProgressTime = nowTime
				createProgress(pixelCount, imageWidth, imageHeight, sample, samplesPerPixel, pixelColors, r.output)
			}
		}

		createProgress(pixelCount, imageWidth, imageHeight, sample, samplesPerPixel, pixelColors, r.output)
	}

	// Apply post processing if applicable, and report a final progress

	if postProcessor != nil && !aborted {

		img := *postProcessor.PostProcess(
			pixelColors,
			albedoColors,
			normalColors,
			imageWidth,
			imageHeight,
			samplesPerPixel,
		)

		r.output <- RenderProgress{
			Progress:    1,
			RenderImage: img,
		}
	}

	close(r.output)
}

func numWorkers() int {
	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}
	return numWorkers
}

func createProgress(
	pixelCount, imageWidth, imageHeight, sample, samplesPerPixel int,
	pixelColors []geo.Vec3,
	output chan<- RenderProgress) {
	ret := make([]color.RGBA, pixelCount)
	for i := 0; i < pixelCount; i++ {
		ret[i] = im.ToRgba(pixelColors[i], sample)
	}
	img := im.RenderImage{
		ImageWidth:  imageWidth,
		ImageHeight: imageHeight,
		Data:        ret,
	}

	output <- RenderProgress{
		Progress:    float64(sample) / float64(samplesPerPixel),
		RenderImage: img,
	}
}

func findLights(s hittable.Hittable, list *hittable.HittableList) {

	switch v := s.(type) {
	case *hittable.HittableList:
		for _, child := range v.List() {
			findLights(child, list)
		}
	default:
		if v.IsLight() {
			list.Add(v)
		}
	}

}
