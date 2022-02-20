package tests

import (
	"image/color"
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
	"github.com/stretchr/testify/assert"
)

func TestRgbToVec3(t *testing.T) {
	assert.Equal(t, geo.NewVec3(0, 0.39215686274509803, 1), im.RgbToVec3(0, 25600, 65535))
}

func TestToRgba(t *testing.T) {
	rgba := im.ToRgba(geo.NewVec3(0, 0.3, 1), 1)
	assert.Equal(t, color.RGBA{R: 0x0, G: 0x8c, B: 0xff, A: 0xff}, rgba)

	rgba = im.ToRgba(geo.NewVec3(0, 0.3, 1), 2)
	assert.Equal(t, color.RGBA{R: 0x0, G: 0x63, B: 0xb5, A: 0xff}, rgba)
}
