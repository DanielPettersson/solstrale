package trace

import (
	"math/rand"
	"time"
)

type TraceProgress struct {
	Progress  float64
	ImageData []byte
}

var (
	world           hittableList
	samplesPerPixel int = 100
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

func RayTrace(imageWidth int, imageHeight int, output chan TraceProgress) {

	rand.Seed(time.Now().UTC().UnixNano())

	// world

	world = hittableList{}
	world.add(sphere{vec3{0, -100.5, -1}, 100})
	world.add(sphere{vec3{0, 0, -1}, 0.5})

	// camera

	camera := createCamera(imageWidth, imageHeight)

	pixels := make([]vec3, imageWidth*imageHeight*4)

	for s := 0; s < samplesPerPixel; s++ {

		i := 0

		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {

				u := (float64(x) + rand.Float64()) / float64(imageWidth-1)
				v := (float64(y) + rand.Float64()) / float64(imageHeight-1)
				r := camera.getRay(u, v)
				pixelColor := rayColor(r, 0)

				pixels[i] = pixels[i].add(pixelColor)
				i += 1
			}

		}

		ret := make([]byte, imageWidth*imageHeight*4)
		i = 0

		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				ri := (((imageHeight-1)-y)*imageWidth + x) * 4

				col := pixels[i].toRgba(s)

				ret[ri] = col.r
				ret[ri+1] = col.g
				ret[ri+2] = col.b
				ret[ri+3] = col.a

				i += 1
			}
		}

		output <- TraceProgress{float64(s+1) / float64(samplesPerPixel), ret}
		time.Sleep(10 * time.Millisecond)
	}

	close(output)
}
