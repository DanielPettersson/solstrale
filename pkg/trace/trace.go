package trace

import (
	"math/rand"
	"time"
)

type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	InterlaceSize   int
	InterlaceOffset int
}

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

func RayTrace(specification TraceSpecification, output chan TraceProgress) {

	rand.Seed(time.Now().UTC().UnixNano())
	imageWidth := specification.ImageWidth
	imageHeight := specification.ImageHeight
	interlaceSize := specification.InterlaceSize
	interlaceOffset := specification.InterlaceOffset

	// world

	world = hittableList{}
	world.add(sphere{vec3{0, -100.5, -1}, 100})
	world.add(sphere{vec3{0, 0, -1}, 0.5})

	// camera

	camera := createCamera(imageWidth, imageHeight)

	pixels := make([]vec3, imageWidth*imageHeight*4)

	for s := 0; s < samplesPerPixel; s++ {

		for y := 0; y < imageHeight; y++ {

			if y%interlaceSize != interlaceOffset {
				continue
			}

			for x := 0; x < imageWidth; x++ {
				i := (((imageHeight-1)-y)*imageWidth + x)

				u := (float64(x) + rand.Float64()) / float64(imageWidth-1)
				v := (float64(y) + rand.Float64()) / float64(imageHeight-1)
				r := camera.getRay(u, v)
				pixelColor := rayColor(r, 0)

				pixels[i] = pixels[i].add(pixelColor)
			}

		}

		ret := make([]byte, imageWidth*imageHeight*4)

		for y := 0; y < imageHeight; y++ {

			if y%interlaceSize != interlaceOffset {
				continue
			}

			for x := 0; x < imageWidth; x++ {
				i := (((imageHeight-1)-y)*imageWidth + x)
				ri := i * 4
				col := pixels[i].toRgba(s)

				ret[ri] = col.r
				ret[ri+1] = col.g
				ret[ri+2] = col.b
				ret[ri+3] = col.a

			}
		}

		reportOutput := s%interlaceSize == interlaceOffset || s == samplesPerPixel-1

		if reportOutput {
			output <- TraceProgress{float64(s+1) / float64(samplesPerPixel), ret}

			// A bit of a hack to let the running web worker context interrupt and do it's callback
			time.Sleep(1 * time.Millisecond)
		}
	}

	close(output)
}
