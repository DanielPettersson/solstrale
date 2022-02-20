package renderer

import (
	"image/color"
	"sync"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/spec"
)

func rayColor(s *spec.Scene, r geo.Ray, depth int, rand util.Random) geo.Vec3 {
	if depth >= s.Spec.MaxDepth {
		return geo.ZeroVector
	}

	hit, rec := s.World.Hit(r, util.Interval{Min: 0.001, Max: util.Infinity}, rand)
	if hit {

		emitted := rec.Material.Emitted(rec)
		scatter, attenuation, scatterRay := rec.Material.Scatter(r, rec, rand)
		if scatter {
			return emitted.Add(attenuation.Mul(rayColor(s, scatterRay, depth+1, rand)))
		}
		return emitted
	}

	return s.BackgroundColor
}

func getRandom(seed uint32, sample, y int) util.Random {
	if seed == 0 {

		// If seed is 0, we will get random
		return util.NewRandom(0)

	} else {

		// Otherwise we want a fixed random, but with
		// non uniform values
		r := util.NewRandom(uint32(sample*y + 1))
		seedRandomizer := r.RandomUint32()
		return util.NewRandom(seed + seedRandomizer)
	}
}

// Render executes the rendering of the image
func Render(s *spec.Scene, output chan spec.TraceProgress, abort chan bool) {

	pixels := make([]geo.Vec3, s.Spec.ImageWidth*s.Spec.ImageHeight)

	for sample := 0; sample < s.Spec.SamplesPerPixel; sample++ {

		select {
		case <-abort:
			close(output)
			return
		default:
		}

		var waitGroup sync.WaitGroup
		for y := 0; y < s.Spec.ImageHeight; y++ {

			waitGroup.Add(1)
			go func(yy int, wg *sync.WaitGroup) {
				defer wg.Done()

				rand := getRandom(s.Spec.RandomSeed, sample, yy)

				for x := 0; x < s.Spec.ImageWidth; x++ {
					i := (((s.Spec.ImageHeight-1)-yy)*s.Spec.ImageWidth + x)

					u := (float64(x) + rand.RandomNormalFloat()) / float64(s.Spec.ImageWidth-1)
					v := (float64(yy) + rand.RandomNormalFloat()) / float64(s.Spec.ImageHeight-1)
					r := s.Cam.GetRay(u, v, rand)
					pixelColor := rayColor(s, r, 0, rand)

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

		output <- spec.TraceProgress{
			Progress: float64(sample+1) / float64(s.Spec.SamplesPerPixel),
			RenderImage: image.RenderImage{
				ImageWidth:  s.Spec.ImageWidth,
				ImageHeight: s.Spec.ImageHeight,
				Data:        ret,
			},
		}
	}

	close(output)
}
