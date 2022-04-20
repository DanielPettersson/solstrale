package material

import (
	"errors"
	"fmt"
	im "image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/image"
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

func NewSolidColor(r, g, b float64) Texture {
	return SolidColor{
		ColorValue: geo.NewVec3(r, g, b),
	}
}

// Color returns the solid color
func (sc SolidColor) Color(rec *HitRecord) geo.Vec3 {
	return sc.ColorValue
}

type imageTexture struct {
	image      im.Image
	mirror     bool
	maxX, maxY float64
}

// LoadImageTexture creates a texture that uses image data for color by loading the image from the path
func LoadImageTexture(path string) (Texture, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to load image texture %v. Got error: %v", path, err.Error()))
	}

	image, _, err := im.Decode(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to decode image texture %v. Got error: %v", path, err.Error()))
	}

	return NewImageTexture(image, false), nil
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
