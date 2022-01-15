package trace

import (
	"math/rand"
	"time"
)

type TraceSpecification struct {
	ImageWidth  int
	ImageHeight int
	DrawOffsetX int
	DrawOffsetY int
	DrawWidth   int
	DrawHeight  int
}

type TraceProgress struct {
	Progress      float64
	Specification TraceSpecification
	ImageData     []byte
}

var (
	world           hittableList
	samplesPerPixel int = 50
	maxDepth        int = 50
)

func rayColor(r ray, depth int) vec3 {
	if depth >= maxDepth {
		return black
	}

	hit, rec := world.hit(r, interval{0.001, infinity})
	if hit {
		target := rec.p.add(rec.normal).add(randomUnitVector())
		ray := ray{rec.p, target.sub(rec.p)}
		return rayColor(ray, depth+1).mulS(0.5)
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
	world.add(sphere{vec3{0, -100.5, -1}, 100})
	world.add(sphere{vec3{0, 0, -1}, 0.5})

	// camera

	camera := createCamera(spec.ImageWidth, spec.ImageHeight)

	pixels := make([]vec3, spec.DrawWidth*spec.DrawHeight)

	for s := 0; s < samplesPerPixel; s++ {

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
			float64(s+1) / float64(samplesPerPixel),
			spec,
			ret,
		}

		// A bit of a hack to let the running web worker context interrupt and do it's callback
		time.Sleep(1 * time.Millisecond)
	}

	close(output)
}
