package trace

import (
	"math/rand"
	"time"
)

type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	DrawOffsetX     int
	DrawOffsetY     int
	DrawWidth       int
	DrawHeight      int
	SamplesPerPixel int
}

type TraceProgress struct {
	Progress      float64
	Specification TraceSpecification
	ImageData     []byte
}

var (
	world    hittableList
	maxDepth int = 50
)

func rayColor(r ray, depth int) vec3 {
	if depth >= maxDepth {
		return black
	}

	hit, rec := world.hit(r, interval{0.001, infinity})
	if hit {

		scatter, attenuation, scatterRay := rec.mat.scatter(r, *rec)
		if scatter {
			return attenuation.mul(rayColor(scatterRay, depth+1))
		}
		return black
	}

	t := 0.5 * (r.dir.unit().y + 1)

	whiteness := white.mulS(1 - t)
	blueness := lightBlue.mulS(t)

	return whiteness.add(blueness)
}

func RayTrace(spec TraceSpecification, output chan TraceProgress) {

	rand.Seed(time.Now().UTC().UnixNano())

	// world

	world = hittableList{}

	materialGround := lambertian{vec3{0.8, 0.8, 0}}
	materialCenter := lambertian{vec3{0.1, 0.2, 0.5}}
	materialLeft := dielectric{vec3{1, 0.8, 0.8}, 1.5}
	materialRight := metal{vec3{0.8, 0.6, 0.2}, 0.1}

	world.add(sphere{vec3{0, -100.5, -1}, 100, materialGround})
	world.add(sphere{vec3{0, 0, -1}, 0.5, materialCenter})
	world.add(sphere{vec3{-1, 0, -1}, 0.5, materialLeft})
	world.add(sphere{vec3{-1, 0, -1}, -0.4, materialLeft})
	world.add(sphere{vec3{1, 0, -1}, 0.5, materialRight})

	// camera

	camera := createCamera(spec.ImageWidth, spec.ImageHeight)

	pixels := make([]vec3, spec.DrawWidth*spec.DrawHeight)

	for s := 0; s < spec.SamplesPerPixel; s++ {

		yStart := spec.ImageHeight - spec.DrawOffsetY - spec.DrawHeight
		for y := yStart; y < yStart+spec.DrawHeight; y++ {

			for x := spec.DrawOffsetX; x < spec.DrawOffsetX+spec.DrawWidth; x++ {
				dx := x - spec.DrawOffsetX
				dy := y - yStart
				i := (((spec.DrawHeight-1)-dy)*spec.DrawWidth + dx)

				u := (float64(x) + rand.Float64()) / float64(spec.ImageWidth-1)
				v := (float64(y) + rand.Float64()) / float64(spec.ImageHeight-1)
				r := camera.getRay(u, v)
				pixelColor := rayColor(r, 0)

				pixels[i] = pixels[i].add(pixelColor)
			}

		}

		ret := make([]byte, len(pixels)*4)

		for y := 0; y < spec.DrawHeight; y++ {
			for x := 0; x < spec.DrawWidth; x++ {

				i := (((spec.DrawHeight-1)-y)*spec.DrawWidth + x)
				ri := i * 4
				col := pixels[i].toRgba(s)

				ret[ri] = col.r
				ret[ri+1] = col.g
				ret[ri+2] = col.b
				ret[ri+3] = col.a

			}
		}

		output <- TraceProgress{
			float64(s+1) / float64(spec.SamplesPerPixel),
			spec,
			ret,
		}

		// A bit of a hack to let the running web worker context interrupt and do it's callback
		time.Sleep(1 * time.Millisecond)
	}

	close(output)
}
