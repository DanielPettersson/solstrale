package trace

import (
	"math/rand"
	"time"
)

type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	DrawOffsetX     int
	DrawOffsetY     int
	DrawWidth       int
	DrawHeight      int
	SamplesPerPixel int
}

type TraceProgress struct {
	Progress      float64
	Specification TraceSpecification
	ImageData     []byte
}

func RayTrace(spec TraceSpecification, output chan TraceProgress) {
	rand.Seed(time.Now().UTC().UnixNano())

	world := hittableList{}

	materialGround := lambertian{vec3{0.8, 0.8, 0}}
	materialCenter := lambertian{vec3{0.1, 0.2, 0.5}}
	materialLeft := dielectric{vec3{1, 0.8, 0.8}, 1.5}
	materialRight := metal{vec3{0.8, 0.6, 0.2}, 0.1}

	world.add(sphere{vec3{0, -100.5, -1}, 100, materialGround})
	world.add(sphere{vec3{0, 0, -1}, 0.5, materialCenter})
	world.add(sphere{vec3{-1, 0, -1}, 0.5, materialLeft})
	world.add(sphere{vec3{-1, 0, -1}, -0.4, materialLeft})
	world.add(sphere{vec3{1, 0, -1}, 0.5, materialRight})

	scene{
		world:  world,
		cam:    createCamera(spec, 20, vec3{-2, 2, 1}, vec3{0, 0, -1}, vec3{0, 1, 0}),
		spec:   spec,
		output: output,
	}.render()

}
