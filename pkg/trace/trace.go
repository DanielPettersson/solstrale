package trace

import (
	"bytes"
	"math/rand"
	"time"
)

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

func RayTrace(imageWidth int, imageHeight int, progress chan float32, byteBuffer *bytes.Buffer) {

	rand.Seed(time.Now().UTC().UnixNano())
	ret := make([]byte, imageWidth*imageHeight*4)

	// world

	world = hittableList{}
	world.add(sphere{vec3{0, -100.5, -1}, 100})
	world.add(sphere{vec3{0, 0, -1}, 0.5})

	// camera

	camera := createCamera(imageWidth, imageHeight)

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			i := (((imageHeight-1)-y)*imageWidth + x) * 4

			pixelColor := vec3{}
			for s := 0; s < samplesPerPixel; s++ {
				u := (float64(x) + rand.Float64()) / float64(imageWidth-1)
				v := (float64(y) + rand.Float64()) / float64(imageHeight-1)
				r := camera.getRay(u, v)
				pixelColor = pixelColor.add(rayColor(r, 0))
			}
			col := pixelColor.toRgba(samplesPerPixel)

			ret[i] = col.r
			ret[i+1] = col.g
			ret[i+2] = col.b
			ret[i+3] = col.a
		}
		reportProgress(y, imageHeight, progress)

	}

	byteBuffer.Write(ret)
	close(progress)
}

func reportProgress(num int, total int, progress chan float32) {

	progressInterval := total / 100
	if progressInterval == 0 {
		progressInterval = 1
	}

	if num%progressInterval == 0 {
		progress <- float32(num) / float32(total)
		time.Sleep(5 * time.Millisecond)
	}
}
