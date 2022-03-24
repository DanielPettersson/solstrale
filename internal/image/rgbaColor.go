package image

import (
	"image/color"
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
)

var (
	colorScale float64 = 1.0 / 255
)

// ToRgba converts a color in a Vec3 that is the sum of a given of amounts of samples
// to a RGBA color. Applies gamma correction to the output color.
func ToRgba(col geo.Vec3, samplesPerPixel int) color.RGBA {
	c := ToFloat(col, samplesPerPixel)

	return color.RGBA{
		byte(256 * c.X),
		byte(256 * c.Y),
		byte(256 * c.Z),
		255,
	}
}

// ToFloat converts a color in a Vec3 that is the sum of a given of amounts of samples
// to a float color. Applies gamma correction to the output color.
func ToFloat(col geo.Vec3, samplesPerPixel int) geo.Vec3 {
	r := col.X
	g := col.Y
	b := col.Z

	// Divide the color by the number of samples
	// and gamma-correct for gamma=2.0
	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	intensity := util.Interval{Min: 0, Max: 0.999}

	return geo.NewVec3(intensity.Clamp(r), intensity.Clamp(g), intensity.Clamp(b))
}

// RgbToVec3 converts rgb bytes to a Vec3 color
func RgbToVec3(r, g, b uint32) geo.Vec3 {
	return geo.Vec3{
		X: float64(r>>8) * colorScale,
		Y: float64(g>>8) * colorScale,
		Z: float64(b>>8) * colorScale,
	}
}
