package image

import (
	"image"
	"image/color"
)

// RenderImage is an image output by the ray tracer.
// Implements the image.Image interface
type RenderImage struct {
	ImageWidth  int
	ImageHeight int
	Data        []color.RGBA
}

// ColorModel returns the RGBA color model
func (ri RenderImage) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the dimensions of the image
func (ri RenderImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, ri.ImageWidth, ri.ImageHeight)
}

// At returns the color at a given point in the image
func (ri RenderImage) At(x, y int) color.Color {
	return ri.Data[y*ri.ImageWidth+x]
}
