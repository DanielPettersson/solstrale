package trace

import (
	"image"
	"image/color"
)

type renderImage struct {
	spec TraceSpecification
	data []color.RGBA
}

func (ri renderImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (ri renderImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, ri.spec.ImageWidth, ri.spec.ImageHeight)
}

func (ri renderImage) At(x, y int) color.Color {
	return ri.data[y*ri.spec.ImageWidth+x]
}
