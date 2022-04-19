package material

import (
	im "image"
	"math"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/image"
	"github.com/ojrac/opensimplex-go"
)

var (
	noise opensimplex.Noise = opensimplex.NewNormalized(123456)
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
	}
	return ct.Odd.Color(rec)
}

type imageTexture struct {
	image      im.Image
	mirror     bool
	maxX, maxY float64
}

// NewImageTexture creates a texture that uses image data for color
func NewImageTexture(image im.Image, mirror bool) Texture {
	max := image.Bounds().Max
	return imageTexture{
		image:  image,
		mirror: mirror,
		maxX:   float64(max.X - 1),
		maxY:   float64(max.Y - 1),
	}
}

// Color returns the color in the image data that corresponds to the UV coordinate of the hittable
// If UV coordinates from hit record is <0 or >1 texture wraps
func (it imageTexture) Color(rec *HitRecord) geo.Vec3 {
	u := math.Mod(math.Abs(rec.U), 1)
	if it.mirror {
		u = 1 - u
	}
	v := 1 - math.Mod(math.Abs(rec.V), 1)

	x := int(u * it.maxX)
	y := int(v * it.maxY)

	r, g, b, _ := it.image.At(x, y).RGBA()
	return image.RgbToVec3(r, g, b)
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
