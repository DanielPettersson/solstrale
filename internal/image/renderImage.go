package image

import (
	"image"
	"image/color"

	"github.com/DanielPettersson/solstrale/spec"
)

type RenderImage struct {
	Spec spec.TraceSpecification
	Data []color.RGBA
}

func (ri RenderImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (ri RenderImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, ri.Spec.ImageWidth, ri.Spec.ImageHeight)
}

func (ri RenderImage) At(x, y int) color.Color {
	return ri.Data[y*ri.Spec.ImageWidth+x]
}
