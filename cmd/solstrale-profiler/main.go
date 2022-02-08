// solstrale-profiler is a debug cmd for generating profiling information.
// This Profiling info cna be inspected using pprof
package main

import (
	"fmt"

	"github.com/DanielPettersson/solstrale"
	"github.com/DanielPettersson/solstrale/spec"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start(profile.ProfilePath(".")).Stop()

	progress := make(chan spec.TraceProgress)
	go solstrale.RayTrace(spec.TraceSpecification{
		ImageWidth:      100,
		ImageHeight:     100,
		SamplesPerPixel: 100,
		RandomSeed:      123456,
	}, progress, make(chan bool))

	for p := range progress {
		fmt.Println(p.Progress)
	}
}
