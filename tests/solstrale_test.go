package tests

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"

	"github.com/DanielPettersson/solstrale"
	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/post"
	"github.com/DanielPettersson/solstrale/renderer"
	"github.com/stretchr/testify/assert"
	"github.com/vitali-fedulov/images3"
)

func createTestScene(renderConfig renderer.RenderConfig) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 20,
		ApertureSize:       .1,
		FocusDistance:      10,
		LookFrom:           geo.NewVec3(-5, 3, 6),
		LookAt:             geo.NewVec3(.25, 1, 0),
	}

	world := hittable.NewHittableList()

	imageTex, err := material.LoadImageTexture("textures/tex.jpg")
	if err != nil {
		panic(err)
	}

	groundMaterial := material.NewLambertian(imageTex)
	glassMat := material.NewDielectric(material.NewSolidColor(1, 1, 1), 1.5)
	lightMat := material.NewLight(10, 10, 10)
	redMat := material.NewLambertian(material.NewSolidColor(1, 0, 0))

	world.Add(hittable.NewQuad(
		geo.NewVec3(-5, 0, -15), geo.NewVec3(20, 0, 0), geo.NewVec3(0, 0, 20),
		groundMaterial,
	))
	world.Add(hittable.NewSphere(geo.NewVec3(-1, 1, 0), 1, glassMat))
	world.Add(hittable.NewRotationY(
		hittable.NewBox(geo.NewVec3(0, 0, -.5), geo.NewVec3(1, 2, .5), redMat),
		15,
	))
	world.Add(hittable.NewConstantMedium(
		hittable.NewTranslation(
			hittable.NewBox(geo.NewVec3(0, 0, -.5), geo.NewVec3(1, 2, .5), nil),
			geo.NewVec3(0, 0, 1),
		),
		0.1,
		geo.NewVec3(1, 1, 1),
	))
	world.Add(hittable.NewMotionBlur(
		hittable.NewBox(geo.NewVec3(-1, 2, 0), geo.NewVec3(-.5, 2.5, .5), redMat),
		geo.NewVec3(0, 1, 0),
	))

	balls := []hittable.Hittable{}
	for i := 0.; i < 1; i += .2 {
		for j := 0.; j < 1; j += .2 {
			for k := 0.; k < 1; k += .2 {
				balls = append(balls, hittable.NewSphere(geo.NewVec3(i, j+.05, k+.8), .05, redMat))
			}
		}
	}
	world.Add(hittable.NewBoundingVolumeHierarchy(balls))

	world.Add(hittable.NewTriangle(geo.NewVec3(1, .1, 2), geo.NewVec3(3, .1, 2), geo.NewVec3(2, .1, 1), redMat))

	// Lights

	world.Add(hittable.NewSphere(geo.NewVec3(10, 5, 10), 10, lightMat))
	world.Add(
		hittable.NewTranslation(
			hittable.NewRotationY(
				hittable.NewQuad(geo.NewVec3(0, 0, 0), geo.NewVec3(2, 0, 0), geo.NewVec3(0, 0, 2), lightMat),
				45,
			),
			geo.NewVec3(-1, 10, -1),
		),
	)
	world.Add(hittable.NewTriangle(geo.NewVec3(-2, 1, -3), geo.NewVec3(0, 1, -3), geo.NewVec3(-1, 2, -3), lightMat))

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}

}

func createBvhTestScene(renderConfig renderer.RenderConfig, useBvh bool, numSpheres int) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 20,
		ApertureSize:       .1,
		FocusDistance:      10,
		LookFrom:           geo.NewVec3(-.5, 0, 4),
		LookAt:             geo.NewVec3(-.5, 0, 0),
	}

	world := hittable.NewHittableList()
	yellow := material.NewLambertian(material.NewSolidColor(1, 1, 0))
	light := material.NewLight(10, 10, 10)
	world.Add(hittable.NewSphere(geo.NewVec3(0, 4, 10), 4, light))

	triangles := []hittable.Hittable{}
	for x := 0.; x < float64(numSpheres); x += 1 {
		cx := x - float64(numSpheres)/2
		s := hittable.NewTriangle(geo.NewVec3(cx, -.5, 0), geo.NewVec3(cx+1, -.5, 0), geo.NewVec3(cx+.5, .5, 0), yellow)
		if useBvh {
			triangles = append(triangles, s)
		} else {
			world.Add(s)
		}
	}

	if useBvh {
		world.Add(hittable.NewBoundingVolumeHierarchy(triangles))
	}

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

func createSimpleTestScene(renderConfig renderer.RenderConfig, addLight bool) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 20,
		ApertureSize:       .1,
		FocusDistance:      10,
		LookFrom:           geo.NewVec3(0, 0, 4),
		LookAt:             geo.NewVec3(0, 0, 0),
	}

	world := hittable.NewHittableList()
	yellow := material.NewLambertian(material.NewSolidColor(1, 1, 0))
	light := material.NewLight(10, 10, 10)
	if addLight {
		world.Add(hittable.NewSphere(geo.NewVec3(0, 100, 0), 20, light))
	}
	world.Add(hittable.NewSphere(geo.NewVec3(0, 0, 0), .5, yellow))

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

func createUvScene(renderConfig renderer.RenderConfig) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 20,
		ApertureSize:       0,
		FocusDistance:      1,
		LookFrom:           geo.NewVec3(0, 1, 5),
		LookAt:             geo.NewVec3(0, 1, 0),
	}

	world := hittable.NewHittableList()
	light := material.NewLight(10, 10, 10)

	world.Add(hittable.NewSphere(geo.NewVec3(50, 50, 50), 20, light))

	tex, err := material.LoadImageTexture("textures/checker.jpg")
	if err != nil {
		panic(err)
	}
	checkerMat := material.NewLambertian(tex)

	world.Add(hittable.NewTriangleWithTexCoords(
		geo.NewVec3(-1, 0, 0),
		geo.NewVec3(1, 0, 0),
		geo.NewVec3(0, 2, 0),
		-1, -1,
		2, -1,
		0, 2,
		checkerMat,
	))

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

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
	model, err := hittable.NewObjModel("spider/", "spider.obj", 1)
	if err != nil {
		panic(err)
	}
	world.Add(model)

	imageTex, err := material.LoadImageTexture("textures/tex.jpg")
	if err != nil {
		panic(err)
	}

	groundMaterial := material.NewLambertian(imageTex)
	world.Add(hittable.NewQuad(geo.NewVec3(-200, -30, -200), geo.NewVec3(400, 0, 0), geo.NewVec3(0, 0, 400), groundMaterial))

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

func createObjWithBox(renderConfig renderer.RenderConfig, path, filename string) *renderer.Scene {
	camera := camera.CameraConfig{
		VerticalFovDegrees: 30,
		ApertureSize:       0,
		FocusDistance:      1,
		LookFrom:           geo.NewVec3(2, 1, 3),
		LookAt:             geo.NewVec3(0, 0, 0),
	}

	world := hittable.NewHittableList()
	light := material.NewLight(15, 15, 15)
	red := material.NewLambertian(material.NewSolidColor(1, 0, 0))

	world.Add(hittable.NewSphere(geo.NewVec3(-100, 100, 40), 35, light))
	model, err := hittable.NewObjModelWithDefaultMaterial(path, filename, 1, red)
	if err != nil {
		panic(err)
	}
	world.Add(model)

	return &renderer.Scene{
		World:           &world,
		Camera:          camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		RenderConfig:    renderConfig,
	}
}

func TestRenderScene(t *testing.T) {

	shaders := map[string]renderer.Shader{
		"pathTracing": renderer.PathTracingShader{MaxDepth: 50},
		"simple":      renderer.SimpleShader{},
	}

	for shaderName, shader := range shaders {

		t.Run(shaderName, func(t *testing.T) {

			traceSpec := renderer.RenderConfig{
				SamplesPerPixel: 25,
				Shader:          shader,
			}
			scene := createTestScene(traceSpec)

			renderAndCompareOutput(t, scene, shaderName, 200, 100)
		})
	}
}

func TestRenderSceneWithOidn(t *testing.T) {

	oidnPost, err := post.NewOidn("./mock_oidn.sh")
	assert.Nil(t, err)

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 10,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
		PostProcessor:   oidnPost,
	}
	scene := createSimpleTestScene(traceSpec, true)

	renderAndCompareOutput(t, scene, "oidn", 200, 100)
}

func TestRenderObjWithTextures(t *testing.T) {

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 20,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createObjScene(traceSpec)

	renderAndCompareOutput(t, scene, "obj", 200, 100)
}

func TestRenderObjWithDefaultMaterial(t *testing.T) {

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 50,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createObjWithBox(traceSpec, "obj/", "box.obj")

	renderAndCompareOutput(t, scene, "obj_default", 200, 100)
}

func TestRenderObjWithDiffuseMaterial(t *testing.T) {

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 50,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createObjWithBox(traceSpec, "obj/", "boxWithMat.obj")

	renderAndCompareOutput(t, scene, "obj_diffuse", 200, 100)
}

func TestRenderUvMapping(t *testing.T) {
	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 5,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createUvScene(traceSpec)

	renderAndCompareOutput(t, scene, "uv", 200, 200)
}

func TestAbortRenderScene(t *testing.T) {

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 100,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createTestScene(traceSpec)

	renderProgress := make(chan renderer.RenderProgress, 1)
	abort := make(chan bool, 1)
	go solstrale.RayTrace(10, 10, scene, renderProgress, abort)

	progressCount := 0
	for range renderProgress {
		progressCount++
		abort <- true
	}

	assert.Equal(t, 1, progressCount)
}

func TestRenderSceneWithoutLight(t *testing.T) {

	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: 100,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createSimpleTestScene(traceSpec, false)

	renderProgress := make(chan renderer.RenderProgress)
	go solstrale.RayTrace(10, 10, scene, renderProgress, make(chan bool))

	for p := range renderProgress {
		assert.Equal(t, "Scene should have at least one light", p.Error.Error())
	}

}

func BenchmarkBvh(b *testing.B) {
	bvh := map[string]bool{
		"with bvh":    true,
		"without bvh": false,
	}

	spheres := []int{10, 100, 1000, 10000}

	for bhvLabel, useBvh := range bvh {
		for _, numSpheres := range spheres {

			b.Run(fmt.Sprintf("%v spheres %v", numSpheres, bhvLabel), func(b *testing.B) {
				b.StopTimer()
				traceSpec := renderer.RenderConfig{
					SamplesPerPixel: b.N,
					Shader:          renderer.PathTracingShader{MaxDepth: 50},
				}
				scene := createBvhTestScene(traceSpec, useBvh, numSpheres)
				b.StartTimer()

				renderProgress := make(chan renderer.RenderProgress)
				go solstrale.RayTrace(20, 10, scene, renderProgress, make(chan bool))
				for range renderProgress {
				}
			})
		}
	}
}

func BenchmarkScene(b *testing.B) {
	traceSpec := renderer.RenderConfig{
		SamplesPerPixel: b.N,
		Shader:          renderer.PathTracingShader{MaxDepth: 50},
	}
	scene := createTestScene(traceSpec)
	renderProgress := make(chan renderer.RenderProgress)

	b.ResetTimer()

	go solstrale.RayTrace(100, 50, scene, renderProgress, make(chan bool))
	for range renderProgress {
	}
}

func renderAndCompareOutput(t *testing.T, scene *renderer.Scene, name string, imageWidth, imageHeight int) {
	renderProgress := make(chan renderer.RenderProgress, 1)
	go solstrale.RayTrace(imageWidth, imageHeight, scene, renderProgress, make(chan bool))

	var im image.Image
	for p := range renderProgress {
		im = p.RenderImage
	}

	actualFileName := fmt.Sprintf("output/out_actual_%v.png", name)
	expectedFileName := fmt.Sprintf("output/out_expected_%v.png", name)

	actualFile, err := os.Create(actualFileName)
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(actualFile, im, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	actualFile.Close()

	actualImage, _ := images3.Open(actualFileName)
	expectedImage, _ := images3.Open(expectedFileName)
	actualIcon := images3.Icon(actualImage, actualFileName)
	expectedIcon := images3.Icon(expectedImage, expectedFileName)

	// Image comparison.
	assert.True(t, images3.Similar(actualIcon, expectedIcon))
}
