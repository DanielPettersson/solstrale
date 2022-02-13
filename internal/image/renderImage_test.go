package image

import (
	"image"
	"image/color"
	"testing"

	"github.com/DanielPettersson/solstrale/spec"
	"github.com/stretchr/testify/assert"
)

func TestColorModel(t *testing.T) {
	assert.Equal(t, color.RGBAModel, RenderImage{}.ColorModel())
}

func TestBounds(t *testing.T) {
	assert.Equal(t, image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 1,
			Y: 2,
		},
	}, RenderImage{
		Spec: spec.TraceSpecification{
			ImageWidth:      1,
			ImageHeight:     2,
			SamplesPerPixel: 0,
			RandomSeed:      0,
		},
		Data: []color.RGBA{},
	}.Bounds())
}

func TestAt(t *testing.T) {

	color1 := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

	color2 := color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}

	image := RenderImage{
		Spec: spec.TraceSpecification{
			ImageWidth:      2,
			ImageHeight:     1,
			SamplesPerPixel: 0,
			RandomSeed:      0,
		},
		Data: []color.RGBA{color1, color2},
	}

	assert.Equal(t, color1, image.At(0, 0))
	assert.Equal(t, color2, image.At(1, 0))
	assert.Panics(t, func() { image.At(0, 1) })
}
