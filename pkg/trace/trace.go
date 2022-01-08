package trace

import (
	"bytes"
)

var (
	white     vec3 = vec3{1, 1, 1}
	lightBlue vec3 = vec3{0.5, 0.7, 1}
)

func rayColor(r ray) vec3 {
	t := 0.5 * (r.dir.unit().y + 1)

	whiteness := white.mulS(1 - t)
	blueness := lightBlue.mulS(t)

	return whiteness.add(blueness)
}

func RayTrace(imageWidth int, imageHeight int, progress chan float32, byteBuffer *bytes.Buffer) {

	ret := make([]byte, imageWidth*imageHeight*4)

	aspectRatio := float64(imageWidth) / float64(imageHeight)

	viewPortHeight := 2.0
	viewPortWidth := aspectRatio * viewPortHeight
	focalLength := 1.0

	origin := vec3{0, 0, 0}
	horizontal := vec3{viewPortWidth, 0, 0}
	vertical := vec3{0, viewPortHeight, 0}
	lowerLeftCorner := origin.sub(horizontal.divS(2))
	lowerLeftCorner = lowerLeftCorner.sub(vertical.divS(2))
	lowerLeftCorner = lowerLeftCorner.sub(vec3{0, 0, focalLength})

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			i := (((imageHeight-1)-y)*imageWidth + x) * 4

			u := float64(x) / float64(imageWidth-1)
			v := float64(y) / float64(imageHeight-1)

			rDir := lowerLeftCorner.add(horizontal.mulS(u))
			rDir = rDir.add(vertical.mulS(v))
			rDir = rDir.sub(origin)
			r := ray{origin, rDir}
			col := rayColor(r).toRgba()

			ret[i] = col.r
			ret[i+1] = col.g
			ret[i+2] = col.b
			ret[i+3] = col.a
		}
		progress <- float32(y) / float32(imageHeight)
	}

	byteBuffer.Write(ret)
	close(progress)
}
