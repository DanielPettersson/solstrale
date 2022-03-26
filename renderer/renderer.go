package renderer

import (
	"image"
	"image/color"
	"sync"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	im "github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/random"
)

// Renderer is a central part of the raytracer responsible for controlling the
// process reporting back progress to the caller
type Renderer struct {
	scene        *Scene
	lights       *hittable.HittableList
	output       chan RenderProgress
	abort        chan bool
	albedoShader AlbedoShader
	normalShader NormalShader
}

// NewRenderer creates a new renderer given a scene and channels for communicating with the caller
func NewRenderer(scene *Scene, output chan RenderProgress, abort chan bool) *Renderer {

	lights := hittable.NewHittableList()
	findLights(scene.World, &lights)

	if len(lights.List()) == 0 {
		panic("Scene should have at least one light")
	}

	return &Renderer{
		scene:        scene,
		lights:       &lights,
		output:       output,
		abort:        abort,
		albedoShader: AlbedoShader{},
		normalShader: NormalShader{},
	}
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
func (r *Renderer) Render() {

	s := r.scene
	imageWidth := s.RenderConfig.ImageWidth
	imageHeight := s.RenderConfig.ImageHeight

	pixelColors := make([]geo.Vec3, imageWidth*imageHeight)
	albedoColors := make([]geo.Vec3, imageWidth*imageHeight)
	normalColors := make([]geo.Vec3, imageWidth*imageHeight)

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
					pixelColor, albedoColor, normalColor := r.rayColor(ray, 0)

					pixelColors[i] = pixelColors[i].Add(pixelColor)

					if r.scene.RenderConfig.PostProcessor != nil {
						albedoColors[i] = albedoColors[i].Add(albedoColor)
						normalColors[i] = normalColors[i].Add(normalColor)
					}
				}
			}(y, &waitGroup)

		}
		waitGroup.Wait()

		var image image.Image

		if r.scene.RenderConfig.PostProcessor != nil {

			image = *r.scene.RenderConfig.PostProcessor.PostProcess(
				pixelColors,
				albedoColors,
				normalColors,
				imageWidth,
				imageHeight,
				sample+1,
			)
		} else {

			ret := make([]color.RGBA, len(pixelColors))

			for y := 0; y < imageHeight; y++ {
				for x := 0; x < imageWidth; x++ {

					i := (((imageHeight-1)-y)*imageWidth + x)
					ret[i] = im.ToRgba(pixelColors[i], sample+1)
				}
			}
			image = im.RenderImage{
				ImageWidth:  imageWidth,
				ImageHeight: imageHeight,
				Data:        ret,
			}
		}

		r.output <- RenderProgress{
			Progress:    float64(sample+1) / float64(s.RenderConfig.SamplesPerPixel),
			RenderImage: image,
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
