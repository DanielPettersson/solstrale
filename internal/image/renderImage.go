package image

import (
	"image"
	"image/color"

	"github.com/DanielPettersson/solstrale/spec"
)

// RenderImage is an image output by the ray tracer.
// Implements the image.Image interface
type RenderImage struct {
	Spec spec.TraceSpecification
	Data []color.RGBA
}

// ColorModel returns the RGBA color model
func (ri RenderImage) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the dimensions of the image
func (ri RenderImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, ri.Spec.ImageWidth, ri.Spec.ImageHeight)
}

// At returns the color at a given point in the image
func (ri RenderImage) At(x, y int) color.Color {
	return ri.Data[y*ri.Spec.ImageWidth+x]
}
