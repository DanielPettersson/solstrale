// Package solstrale is the main package for the ray tracer and contains the functions for starting the raytracing
package solstrale

import (
	"github.com/DanielPettersson/solstrale/renderer"
)

// RayTrace executes the ray tracing with the given scene and reports progress on
// the output channel. Listens to abort channel for aborting a started ray trace operation
func RayTrace(width, height int, scene *renderer.Scene, output chan<- renderer.RenderProgress, abort <-chan bool) {
	r, err := renderer.NewRenderer(scene, output, abort)
	if err != nil {
		output <- renderer.RenderProgress{
			Error: err,
		}
		close(output)
		return
	}

	r.Render(width, height)
}
