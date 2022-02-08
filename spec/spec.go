package spec

import "image"

type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	SamplesPerPixel int
	RandomSeed      int
}

type TraceProgress struct {
	Progress    float64
	RenderImage image.Image
}
