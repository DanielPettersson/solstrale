// Package solstrale is the main package for the ray tracer and contains the functions for starting the raytracing
package solstrale

import (
	"github.com/DanielPettersson/solstrale/renderer"
)

// RayTrace executes the ray tracing with the given scene and reports progress on
// the output channel. Listens to abort channel for aborting a started ray trace operation
func RayTrace(scene *renderer.Scene, output chan renderer.RenderProgress, abort chan bool) {
	aborted, pixels := renderer.NewRenderer(scene, output, abort).Render()

	if scene.RenderConfig.PostProcessor != nil && !aborted {
		postProcessed, image := scene.RenderConfig.PostProcessor.PostProcess(pixels, scene.RenderConfig.ImageWidth, scene.RenderConfig.ImageHeight)

		if postProcessed {

			output <- renderer.RenderProgress{
				Progress:    1,
				RenderImage: *image,
			}
		}
	}

	close(output)
}
