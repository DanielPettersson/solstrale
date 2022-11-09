// Package post provides post processor for the rendered image
package post

import (
	"image"

	"github.com/DanielPettersson/solstrale/geo"
)

// PostProcessor is responsible for taking the rendered image and transforming it
type PostProcessor interface {
	PostProcess(
		pixelColors []geo.Vec3,
		albedoColors []geo.Vec3,
		normalColors []geo.Vec3,
		width, height, numSamples int,
	) (image.Image, error)
}
