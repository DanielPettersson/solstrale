package renderer

import (
	"image/color"
	"math"
	"sync"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/pdf"
	"github.com/DanielPettersson/solstrale/random"
	"github.com/DanielPettersson/solstrale/spec"
)

type Renderer struct {
	scene  *spec.Scene
	lights *hittable.HittableList
	output chan spec.TraceProgress
	abort  chan bool
}

func NewRenderer(scene *spec.Scene, output chan spec.TraceProgress, abort chan bool) *Renderer {

	lights := hittable.NewHittableList()
	findLights(scene.World, &lights)

	return &Renderer{
		scene:  scene,
		lights: &lights,
		output: output,
		abort:  abort,
	}
}

func (r *Renderer) rayColor(ray geo.Ray, depth int) geo.Vec3 {

	s := r.scene

	if depth >= s.Spec.MaxDepth {
		return geo.ZeroVector
	}

	hit, rec := s.World.Hit(ray, util.Interval{Min: 0.001, Max: util.Infinity})
	if hit {

		emittedColor := rec.Material.Emitted(rec)
		scatter, scatterRecord := rec.Material.Scatter(ray, rec)
		if !scatter {
			return emittedColor
		}

		if scatterRecord.SkipPdf {
			return scatterRecord.Attenuation.Mul(r.rayColor(scatterRecord.SkipPdfRay, depth+1))
		}

		lightPtr := hittable.NewHittablePdf(r.lights, rec.HitPoint)
		mixturePdf := pdf.NewMixturePdf(&lightPtr, scatterRecord.PdfPtr)

		scattered := geo.Ray{
			Origin:    rec.HitPoint,
			Direction: mixturePdf.Generate(),
			Time:      ray.Time,
		}
		pdfVal := mixturePdf.Value(scattered.Direction)
		scatteringPdf := rec.Material.ScatteringPdf(ray, rec, scattered)
		scatterColor := scatterRecord.Attenuation.MulS(scatteringPdf).Mul(r.rayColor(scattered, depth+1)).DivS(pdfVal)

		return filterInvalidColorValues(emittedColor.Add(scatterColor))
	}

	return s.BackgroundColor
}

// Render executes the rendering of the image
func (r *Renderer) Render() {

	s := r.scene

	pixels := make([]geo.Vec3, s.Spec.ImageWidth*s.Spec.ImageHeight)

	for sample := 0; sample < s.Spec.SamplesPerPixel; sample++ {

		select {
		case <-r.abort:
			close(r.output)
			return
		default:
		}

		var waitGroup sync.WaitGroup
		for y := 0; y < s.Spec.ImageHeight; y++ {

			waitGroup.Add(1)
			go func(yy int, wg *sync.WaitGroup) {
				defer wg.Done()

				for x := 0; x < s.Spec.ImageWidth; x++ {
					i := (((s.Spec.ImageHeight-1)-yy)*s.Spec.ImageWidth + x)

					u := (float64(x) + random.RandomNormalFloat()) / float64(s.Spec.ImageWidth-1)
					v := (float64(yy) + random.RandomNormalFloat()) / float64(s.Spec.ImageHeight-1)
					ray := s.Cam.GetRay(u, v)
					pixelColor := r.rayColor(ray, 0)

					pixels[i] = pixels[i].Add(pixelColor)
				}
			}(y, &waitGroup)

		}
		waitGroup.Wait()

		ret := make([]color.RGBA, len(pixels))

		for y := 0; y < s.Spec.ImageHeight; y++ {
			for x := 0; x < s.Spec.ImageWidth; x++ {

				i := (((s.Spec.ImageHeight-1)-y)*s.Spec.ImageWidth + x)
				ret[i] = image.ToRgba(pixels[i], sample+1)
			}
		}

		r.output <- spec.TraceProgress{
			Progress: float64(sample+1) / float64(s.Spec.SamplesPerPixel),
			RenderImage: image.RenderImage{
				ImageWidth:  s.Spec.ImageWidth,
				ImageHeight: s.Spec.ImageHeight,
				Data:        ret,
			},
		}
	}

	close(r.output)
}

func filterInvalidColorValues(col geo.Vec3) geo.Vec3 {
	return geo.NewVec3(
		filterColorValue(col.X),
		filterColorValue(col.Y),
		filterColorValue(col.Z),
	)
}

func filterColorValue(val float64) float64 {
	if math.IsNaN(val) {
		return 0
	}
	if val > 100 {
		return 100
	}
	return val
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
