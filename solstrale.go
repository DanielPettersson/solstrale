// Package solstrale is the main package for the ray tracer and contains the functions for starting the raytracing
package solstrale

import (
	"github.com/DanielPettersson/solstrale/internal/renderer"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/spec"
)

// RayTrace executes the ray tracing with the given scene and reports progress on
// the output channel. Listens to abort channel for aborting a started ray trace operation
func RayTrace(scene *spec.Scene, output chan spec.TraceProgress, abort chan bool) {
	util.SetRandomSeed(uint32(scene.Spec.RandomSeed))
	renderer.Render(scene, output, abort)
}
