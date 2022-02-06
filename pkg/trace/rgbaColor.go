package trace

import (
	"image/color"
	"math"
)

var (
	colorScale float64 = 1.0 / 255
)

func (c vec3) toRgba(samplesPerPixel int) color.RGBA {
	r := c.x
	g := c.y
	b := c.z

	// Divide the color by the number of samples
	// and gamma-correct for gamma=2.0
	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	intensity := interval{0, 0.999}

	return color.RGBA{
		byte(256 * intensity.clamp(r)),
		byte(256 * intensity.clamp(g)),
		byte(256 * intensity.clamp(b)),
		255,
	}
}

func rgbToVec3(r, g, b byte) vec3 {
	return vec3{
		float64(r) * colorScale,
		float64(g) * colorScale,
		float64(b) * colorScale,
	}
}
