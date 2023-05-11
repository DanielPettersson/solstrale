package main

import (
	"fmt"
	"image/jpeg"
	"os"

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
		ApertureSize:       0,
		FocusDistance:      10,
		LookFrom:           geo.NewVec3(1, .05, 0),
		LookAt:             geo.NewVec3(0, .05, 0),
	}

	world := hittable.NewHittableList()
	light := material.NewLight(15, 15, 15)

	world.Add(hittable.NewSphere(geo.NewVec3(100, 100, 100), 35, light))

	dragon, _ := hittable.NewObjModel("", "dragon.obj", 1)
	world.Add(dragon)

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
		SamplesPerPixel: 1,
		Shader:          renderer.SimpleShader{},
	}
	scene := createObjScene(renderConfig)

	renderProgress := make(chan renderer.RenderProgress, 1)
	go solstrale.RayTrace(400, 200, scene, renderProgress, make(<-chan bool))

	for p := range renderProgress {
		if p.Error != nil {
			fmt.Println(p.Error.Error())
			return
		}
		f, err := os.Create("out.jpg")
		if err != nil {
			fmt.Println(p.Error.Error())
			return
		}
		jpeg.Encode(f, p.RenderImage, nil)
	}
}
