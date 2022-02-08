package material

import (
	"math"
	"math/rand"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/ojrac/opensimplex-go"
)

var (
	noise opensimplex.Noise = opensimplex.NewNormalized(rand.Int63())
)

// Texture describes the color of a material.
// The color can vary by the uv coordinates of the hittable
type Texture interface {
	Color(rec *HitRecord) geo.Vec3
}

// SolidColor is a texture with just the same color everywhere
type SolidColor struct {
	ColorValue geo.Vec3
}

// Color returns the solid color
func (sc SolidColor) Color(rec *HitRecord) geo.Vec3 {
	return sc.ColorValue
}

// CheckerTexture is a checkered texture
type CheckerTexture struct {
	Scale float64
	Even  Texture
	Odd   Texture
}

// Color returns either Even of Odd color depending on the UV coordinates of the hittable
func (ct CheckerTexture) Color(rec *HitRecord) geo.Vec3 {
	invScale := 1 / ct.Scale
	uInt := math.Floor(rec.U * invScale)
	vInt := math.Floor(rec.V * invScale)

	if int(uInt+vInt)%2 == 0 {
		return ct.Even.Color(rec)
	} else {
		return ct.Odd.Color(rec)
	}

}

// ImageTexture is a texture that uses image data for color
type ImageTexture struct {
	Bytes  []byte
	Width  int
	Height int
	Mirror bool
}

// Color returns the color in the image data that corresponds to the UV coordinate of the hittable
func (it ImageTexture) Color(rec *HitRecord) geo.Vec3 {
	u := util.Interval{Min: 0, Max: 1}.Clamp(rec.U)
	if it.Mirror {
		u = 1 - u
	}
	v := 1 - util.Interval{Min: 0, Max: 1}.Clamp(rec.V)

	x := int(u * float64(it.Width))
	y := int(v * float64(it.Height))
	i := (y*it.Width + x) * 4

	return image.RgbToVec3(
		it.Bytes[i],
		it.Bytes[i+1],
		it.Bytes[i+2],
	)
}

// NoiseTexture is a "random" noise texture
type NoiseTexture struct {
	ColorValue geo.Vec3
	Scale      float64
}

// Color returns the "random" color at the given UV coordinate given the simplex noise algorithm
func (nt NoiseTexture) Color(rec *HitRecord) geo.Vec3 {
	p := rec.HitPoint.MulS(1 / nt.Scale)
	val := noise.Eval3(p.X, p.Y, p.Z)
	return nt.ColorValue.MulS(val)
}
