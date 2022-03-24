package renderer

import (
	"image"

	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/post"
)

// RenderConfig is input to the ray tracer for how the image should be rendered
type RenderConfig struct {
	ImageWidth      int
	ImageHeight     int
	SamplesPerPixel int
	Shader          Shader
	PostProcessor   post.PostProcessor
}

// Scene contains all information needed to render an image
type Scene struct {
	World           hittable.Hittable
	Cam             camera.Camera
	BackgroundColor geo.Vec3
	RenderConfig    RenderConfig
}

// RenderProgress is progress reported back to the caller of the raytrace function
type RenderProgress struct {
	Progress    float64
	RenderImage image.Image
}
