package post

import (
	"image"
	"image/color"

	"github.com/DanielPettersson/solstrale/geo"
	im "github.com/DanielPettersson/solstrale/internal/image"
)

type bloomPostProcessor struct {
}

// NewBloom returns a new bloom post processor instance
func NewBloom() PostProcessor {
	return bloomPostProcessor{}
}

// PostProcess applies a bloom filter to the renderered image
func (b bloomPostProcessor) PostProcess(
	pixelColors []geo.Vec3,
	albedoColors []geo.Vec3,
	normalColors []geo.Vec3,
	width, height, numSamples int,
) (image.Image, error) {

	pixelCount := width * height
	ret := make([]color.RGBA, pixelCount)
	for i := 0; i < pixelCount; i++ {
		if albedoColors[i].Length() > float64(numSamples) {
			ret[i] = color.RGBA{
				R: 0,
				G: 255,
				B: 0,
				A: 255,
			}
		} else {
			ret[i] = im.ToRgba(pixelColors[i], numSamples)
		}
	}
	img := im.RenderImage{
		ImageWidth:  width,
		ImageHeight: height,
		Data:        ret,
	}

	return img, nil
}
