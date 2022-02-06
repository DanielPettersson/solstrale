package main

import (
	"fmt"

	"github.com/DanielPettersson/solstrale/pkg/trace"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start(profile.ProfilePath(".")).Stop()

	progress := make(chan trace.TraceProgress)
	go trace.RayTrace(trace.TraceSpecification{
		ImageWidth:      100,
		ImageHeight:     100,
		DrawOffsetX:     0,
		DrawOffsetY:     0,
		DrawWidth:       100,
		DrawHeight:      100,
		SamplesPerPixel: 100,
		RandomSeed:      123456,
	}, progress, make(chan bool))

	for p := range progress {
		fmt.Println(p.Progress)
	}
}
