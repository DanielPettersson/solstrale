package trace

import (
	"image/color"
	"sync"
)

var (
	maxDepth int = 50
)

type scene struct {
	world           hittableList
	cam             camera
	backgroundColor vec3
	spec            TraceSpecification
	output          chan TraceProgress
}

func (s scene) rayColor(r ray, depth int) vec3 {
	if depth >= maxDepth {
		return black
	}

	hit, rec := s.world.hit(r, interval{0.001, infinity})
	if hit {

		emitted := rec.material.emitted(rec)
		scatter, attenuation, scatterRay := rec.material.scatter(r, rec)
		if scatter {
			return emitted.add(attenuation.mul(s.rayColor(scatterRay, depth+1)))
		}
		return emitted
	}

	return s.backgroundColor
}

func (s scene) render(abort chan bool) {

	spec := s.spec
	pixels := make([]vec3, spec.ImageWidth*spec.ImageHeight)

	for sample := 0; sample < spec.SamplesPerPixel; sample++ {

		select {
		case <-abort:
			close(s.output)
			return
		default:
		}

		var waitGroup sync.WaitGroup
		for y := 0; y < spec.ImageHeight; y++ {

			waitGroup.Add(1)
			go func(yy int, wg *sync.WaitGroup) {
				defer wg.Done()

				for x := 0; x < spec.ImageWidth; x++ {
					i := (((spec.ImageHeight-1)-yy)*spec.ImageWidth + x)

					u := (float64(x) + randomNormalFloat()) / float64(spec.ImageWidth-1)
					v := (float64(yy) + randomNormalFloat()) / float64(spec.ImageHeight-1)
					r := s.cam.getRay(u, v)
					pixelColor := s.rayColor(r, 0)

					pixels[i] = pixels[i].add(pixelColor)
				}
			}(y, &waitGroup)

		}
		waitGroup.Wait()

		ret := make([]color.RGBA, len(pixels))

		for y := 0; y < spec.ImageHeight; y++ {
			for x := 0; x < spec.ImageWidth; x++ {

				i := (((spec.ImageHeight-1)-y)*spec.ImageWidth + x)
				ret[i] = pixels[i].toRgba(sample)
			}
		}

		s.output <- TraceProgress{
			float64(sample+1) / float64(spec.SamplesPerPixel),
			renderImage{spec, ret},
		}
	}

	close(s.output)
}
