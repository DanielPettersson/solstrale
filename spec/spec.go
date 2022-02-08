// Package spec provides data structures for controlling the ray tracing operation
package spec

import "image"

// TraceSpecification is input to the ray tracer for how the image should be rendered
type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	SamplesPerPixel int
	RandomSeed      int
}

// TraceProgress is progress reported back to the caller of the raytrace function
type TraceProgress struct {
	Progress    float64
	RenderImage image.Image
}
