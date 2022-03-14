package renderer

import (
	"image/color"
	"sync"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
)

// Renderer is a central part of the raytracer responsible for controlling the
// process reporting back progress to the caller
type Renderer struct {
	scene  *Scene
	lights *hittable.HittableList
	output chan RenderProgress
	abort  chan bool
}

// NewRenderer creates a new renderer given a scene and channels for communicating with the caller
func NewRenderer(scene *Scene, output chan RenderProgress, abort chan bool) *Renderer {

	lights := hittable.NewHittableList()
	findLights(scene.World, &lights)

	if len(lights.List()) == 0 {
		panic("Scene should have at least one light")
	}

	return &Renderer{
		scene:  scene,
		lights: &lights,
		output: output,
		abort:  abort,
	}
}

func (r *Renderer) rayColor(ray geo.Ray, depth int) geo.Vec3 {

	s := r.scene

	hit, rec := s.World.Hit(ray, util.Interval{Min: 0.001, Max: util.Infinity})
	if hit {
		return s.RenderConfig.Shader.Shade(r, rec, ray, depth)
	}

	return s.BackgroundColor
}

// Render executes the rendering of the image
func (r *Renderer) Render() {

	s := r.scene
	imageWidth := s.RenderConfig.ImageWidth
	imageHeight := s.RenderConfig.ImageHeight

	pixels := make([]geo.Vec3, imageWidth*imageHeight)

	for sample := 0; sample < s.RenderConfig.SamplesPerPixel; sample++ {

		select {
		case <-r.abort:
			close(r.output)
			return
		default:
		}

		var waitGroup sync.WaitGroup
		for y := 0; y < imageHeight; y++ {

			waitGroup.Add(1)
			go func(yy int, wg *sync.WaitGroup) {
				defer wg.Done()

				for x := 0; x < imageWidth; x++ {
					i := (((imageHeight-1)-yy)*imageWidth + x)

					u := (float64(x) + random.RandomNormalFloat()) / float64(imageWidth-1)
					v := (float64(yy) + random.RandomNormalFloat()) / float64(imageHeight-1)
					ray := s.Cam.GetRay(u, v)
					pixelColor := r.rayColor(ray, 0)

					pixels[i] = pixels[i].Add(pixelColor)
				}
			}(y, &waitGroup)

		}
		waitGroup.Wait()

		ret := make([]color.RGBA, len(pixels))

		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {

				i := (((imageHeight-1)-y)*imageWidth + x)
				ret[i] = image.ToRgba(pixels[i], sample+1)
			}
		}

		r.output <- RenderProgress{
			Progress: float64(sample+1) / float64(s.RenderConfig.SamplesPerPixel),
			RenderImage: image.RenderImage{
				ImageWidth:  imageWidth,
				ImageHeight: imageHeight,
				Data:        ret,
			},
		}
	}

	close(r.output)
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
