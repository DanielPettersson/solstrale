package trace

import (
	"bytes"
)

func RayTrace(width int, height int, progress chan float32, byteBuffer *bytes.Buffer) {

	ret := make([]byte, width*height*4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4

			col := color{
				float64(x) / float64(width),
				float64(y) / float64(height),
				0.25,
			}.toRgba()

			ret[i] = col.r
			ret[i+1] = col.g
			ret[i+2] = col.b
			ret[i+3] = col.a
		}
		progress <- float32(y) / float32(height)
	}

	byteBuffer.Write(ret)
	close(progress)
}
