package renderer

import (
	"image/color"
	"sync"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/spec"
)

func rayColor(s *spec.Scene, r geo.Ray, depth int) geo.Vec3 {
	if depth >= s.Spec.MaxDepth {
		return geo.ZeroVector
	}

	hit, rec := s.World.Hit(r, util.Interval{Min: 0.001, Max: util.Infinity})
	if hit {

		emitted := rec.Material.Emitted(rec)
		scatter, attenuation, scatterRay := rec.Material.Scatter(r, rec)
		if scatter {
			return emitted.Add(attenuation.Mul(rayColor(s, scatterRay, depth+1)))
		}
		return emitted
	}

	return s.BackgroundColor
}

// Render executes the rendering of the image
func Render(s *spec.Scene, abort chan bool) {

	pixels := make([]geo.Vec3, s.Spec.ImageWidth*s.Spec.ImageHeight)

	for sample := 0; sample < s.Spec.SamplesPerPixel; sample++ {

		select {
		case <-abort:
			close(s.Output)
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

					u := (float64(x) + util.RandomNormalFloat()) / float64(s.Spec.ImageWidth-1)
					v := (float64(yy) + util.RandomNormalFloat()) / float64(s.Spec.ImageHeight-1)
					r := s.Cam.GetRay(u, v)
					pixelColor := rayColor(s, r, 0)

					pixels[i] = pixels[i].Add(pixelColor)
				}
			}(y, &waitGroup)

		}
		waitGroup.Wait()

		ret := make([]color.RGBA, len(pixels))

		for y := 0; y < s.Spec.ImageHeight; y++ {
			for x := 0; x < s.Spec.ImageWidth; x++ {

				i := (((s.Spec.ImageHeight-1)-y)*s.Spec.ImageWidth + x)
				ret[i] = image.ToRgba(pixels[i], sample)
			}
		}

		s.Output <- spec.TraceProgress{
			Progress: float64(sample+1) / float64(s.Spec.SamplesPerPixel),
			RenderImage: image.RenderImage{
				ImageWidth:  s.Spec.ImageWidth,
				ImageHeight: s.Spec.ImageHeight,
				Data:        ret,
			},
		}
	}

	close(s.Output)
}
