package trace

import (
	"math/rand"
	"time"
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

func (s scene) render() {

	spec := s.spec
	pixels := make([]vec3, spec.DrawWidth*spec.DrawHeight)

	for sample := 0; sample < spec.SamplesPerPixel; sample++ {

		yStart := spec.ImageHeight - spec.DrawOffsetY - spec.DrawHeight
		for y := yStart; y < yStart+spec.DrawHeight; y++ {

			for x := spec.DrawOffsetX; x < spec.DrawOffsetX+spec.DrawWidth; x++ {
				dx := x - spec.DrawOffsetX
				dy := y - yStart
				i := (((spec.DrawHeight-1)-dy)*spec.DrawWidth + dx)

				u := (float64(x) + rand.Float64()) / float64(spec.ImageWidth-1)
				v := (float64(y) + rand.Float64()) / float64(spec.ImageHeight-1)
				r := s.cam.getRay(u, v)
				pixelColor := s.rayColor(r, 0)

				pixels[i] = pixels[i].add(pixelColor)
			}

		}

		ret := make([]byte, len(pixels)*4)

		for y := 0; y < spec.DrawHeight; y++ {
			for x := 0; x < spec.DrawWidth; x++ {

				i := (((spec.DrawHeight-1)-y)*spec.DrawWidth + x)
				ri := i * 4
				col := pixels[i].toRgba(sample)

				ret[ri] = col.r
				ret[ri+1] = col.g
				ret[ri+2] = col.b
				ret[ri+3] = col.a

			}
		}

		s.output <- TraceProgress{
			float64(sample+1) / float64(spec.SamplesPerPixel),
			spec,
			ret,
		}

		// A bit of a hack to let the running web worker context interrupt and do it's callback
		time.Sleep(1 * time.Millisecond)
	}

	close(s.output)
}
