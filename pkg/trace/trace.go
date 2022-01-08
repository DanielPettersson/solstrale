package trace

import (
	"bytes"
)

func RayTrace(width int, height int, progress chan float32, byteBuffer *bytes.Buffer) {

	ret := make([]byte, width*height*4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4

			r := float32(x) / float32(width)
			g := float32(y) / float32(height)
			b := 0.25
			a := 1.0

			ret[i] = byte(r * 255.9)
			ret[i+1] = byte(g * 255.9)
			ret[i+2] = byte(b * 255.9)
			ret[i+3] = byte(a * 255.9)
		}
		progress <- float32(y) / float32(height)
	}

	byteBuffer.Write(ret)
	close(progress)
}
