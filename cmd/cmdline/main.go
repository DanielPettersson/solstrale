package main

import (
	"fmt"
	"log"
	"time"

	"github.com/DanielPettersson/solstrale/pkg/trace"
)

func main() {

	start := time.Now()

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
	}, progress)

	for p := range progress {
		fmt.Println(p.Progress)
	}

	elapsed := time.Since(start)
	log.Printf("Raytracing took %s", elapsed)

}
