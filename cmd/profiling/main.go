package main

import (
	"fmt"

	"github.com/DanielPettersson/solstrale"
	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/renderer"
	"github.com/pkg/profile"
)

func createObjScene(renderConfig renderer.RenderConfig) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 30,
		ApertureSize:       20,
		FocusDistance:      260,
		LookFrom:           geo.NewVec3(-250, 30, 150),
		LookAt:             geo.NewVec3(-50, 0, 0),
	}

	world := hittable.NewHittableList()
	light := material.NewLight(15, 15, 15)

	world.Add(hittable.NewSphere(geo.NewVec3(-100, 100, 40), 35, light))

	// Add scene to profile here...

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

// Usage tips:
// 1. run this
// 2. Display profile info by: pprof -http=localhost:8888 profiling cpu.pprof
func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	renderConfig := renderer.RenderConfig{
		SamplesPerPixel: 100,
		Shader: renderer.PathTracingShader{
			MaxDepth: 50,
		},
	}
	scene := createObjScene(renderConfig)

	renderProgress := make(chan renderer.RenderProgress, 1)
	go solstrale.RayTrace(400, 200, scene, renderProgress, make(<-chan bool))

	for p := range renderProgress {
		fmt.Println(p.Progress)
	}
}
