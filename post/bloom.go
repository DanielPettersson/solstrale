package post

import (
	"image"

	"github.com/DanielPettersson/solstrale/geo"
)

type bloomPostProcessor struct {
}

// PostProcess applies a bloom filter to the renderered image
func (b bloomPostProcessor) PostProcess(
	pixelColors []geo.Vec3,
	albedoColors []geo.Vec3,
	normalColors []geo.Vec3,
	width, height, numSamples int,
) (image.Image, error) {
	return nil, nil
}
