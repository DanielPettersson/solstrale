// Package spec provides data structures for controlling the ray tracing operation
package spec

import (
	"image"

	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
)

// Scene contains all information needed to render an image
type Scene struct {
	World           hittable.Hittable
	Lights          hittable.HittableList
	Cam             camera.Camera
	BackgroundColor geo.Vec3
	Spec            TraceSpecification
}

// TraceSpecification is input to the ray tracer for how the image should be rendered
type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	SamplesPerPixel int
	MaxDepth        int
}

// TraceProgress is progress reported back to the caller of the raytrace function
type TraceProgress struct {
	Progress    float64
	RenderImage image.Image
}
