package trace

import "math"

var (
	colorScale float64 = 1.0 / 255
)

type rgbaColor struct {
	r byte
	g byte
	b byte
	a byte
}

func (c vec3) toRgba(samplesPerPixel int) rgbaColor {
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

	return rgbaColor{
		byte(256 * intensity.clamp(r)),
		byte(256 * intensity.clamp(g)),
		byte(256 * intensity.clamp(b)),
		255,
	}
}

func (rgba rgbaColor) toVec3() vec3 {
	return vec3{
		float64(rgba.r) * colorScale,
		float64(rgba.g) * colorScale,
		float64(rgba.b) * colorScale,
	}
}
