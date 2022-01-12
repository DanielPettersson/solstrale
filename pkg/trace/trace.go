package trace

import (
	"bytes"
	"math/rand"
	"time"
)

var (
	world           hittableList
	samplesPerPixel int = 50
)

func rayColor(r ray) vec3 {
	hit, hitRecord := world.hit(r, interval{0, infinity})
	if hit {
		return hitRecord.normal.add(vec3{1, 1, 1}).mulS(0.5)
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
				pixelColor = pixelColor.add(rayColor(r))
			}
			col := pixelColor.toRgba(samplesPerPixel)

			ret[i] = col.r
			ret[i+1] = col.g
			ret[i+2] = col.b
			ret[i+3] = col.a
		}
		progress <- float32(y) / float32(imageHeight)
		if y%(imageHeight/100) == 0 {
			time.Sleep(2 * time.Millisecond)
		}
	}

	byteBuffer.Write(ret)
	close(progress)
}
