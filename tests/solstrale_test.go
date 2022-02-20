package tests

import (
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
	"github.com/DanielPettersson/solstrale/spec"
	"github.com/stretchr/testify/assert"
)

func createTestScene(traceSpec spec.TraceSpecification) *spec.Scene {
	camera := camera.New(
		traceSpec.ImageWidth,
		traceSpec.ImageHeight,
		20,
		0.1,
		10,
		geo.NewVec3(-5, 3, 6),
		geo.NewVec3(.25, 1, 0),
		geo.NewVec3(0, 1, 0),
	)

	world := hittable.NewHittableList()

	checkerTex := material.CheckerTexture{
		Scale: 0.1,
		Even:  material.SolidColor{ColorValue: geo.NewVec3(0.4, 0.2, 0.1)},
		Odd:   material.SolidColor{ColorValue: geo.NewVec3(0.8, 0.4, 0.2)},
	}
	noiseTex := material.NoiseTexture{
		ColorValue: geo.NewVec3(1, 1, 0),
		Scale:      .05,
	}

	f, _ := os.Open("tex.jpg")
	defer f.Close()
	image, _, _ := image.Decode(f)
	imageTex := material.ImageTexture{
		Image:  image,
		Mirror: false,
	}

	groundMaterial := material.Lambertian{Tex: imageTex}
	checkerMat := material.Lambertian{Tex: checkerTex}
	glassMat := material.Dielectric{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}, IndexOfRefraction: 1.5}
	goldMat := material.Metal{Tex: noiseTex, Fuzz: .2}
	lightMat := material.DiffuseLight{Emit: material.SolidColor{ColorValue: geo.NewVec3(5, 5, 5)}}
	fogMat := material.Isotropic{Albedo: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}}
	redMat := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 0, 0)}}

	world.Add(hittable.NewQuad(
		geo.NewVec3(-5, 0, -15), geo.NewVec3(20, 0, 0), geo.NewVec3(0, 0, 20),
		groundMaterial,
	))
	world.Add(hittable.NewSphere(geo.NewVec3(-1, 1, 0), 1, glassMat))
	world.Add(hittable.NewRotationY(
		hittable.NewBox(geo.NewVec3(0, 0, -.5), geo.NewVec3(1, 2, .5), checkerMat),
		15,
	))
	world.Add(hittable.NewConstantMedium(
		hittable.NewTranslation(
			hittable.NewBox(geo.NewVec3(0, 0, -.5), geo.NewVec3(1, 2, .5), fogMat),
			geo.NewVec3(0, 0, 1),
		),
		0.1,
		material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)},
	))
	world.Add(hittable.NewSphere(geo.NewVec3(2, 1, 0), 1, goldMat))
	world.Add(hittable.NewSphere(geo.NewVec3(10, 5, 10), 10, lightMat))

	world.Add(hittable.NewMotionBlur(
		hittable.NewBox(geo.NewVec3(-1, 2, 0), geo.NewVec3(-.5, 2.5, .5), redMat),
		geo.NewVec3(0, 1, 0),
	))

	balls := hittable.NewHittableList()
	for i := 0.; i < 1; i += .2 {
		for j := 0.; j < 1; j += .2 {
			for k := 0.; k < 1; k += .2 {
				balls.Add(hittable.NewSphere(geo.NewVec3(i, j+.05, k+.8), .05, redMat))
			}
		}
	}

	world.Add(hittable.NewBoundingVolumeHierarchy(balls))

	return &spec.Scene{
		World:           &world,
		Cam:             camera,
		BackgroundColor: geo.NewVec3(.2, .3, .5),
		Spec:            traceSpec,
	}

}

func TestRenderScene(t *testing.T) {

	traceSpec := spec.TraceSpecification{
		ImageWidth:      100,
		ImageHeight:     50,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		RandomSeed:      123,
	}
	scene := createTestScene(traceSpec)

	renderProgress := make(chan spec.TraceProgress, 1)
	go solstrale.RayTrace(scene, renderProgress, make(chan bool))

	var im image.Image
	for p := range renderProgress {
		im = p.RenderImage
	}

	actualFile, err := os.Create("out_actual.png")
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(actualFile, im, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	actualFile.Close()

	expectedFile, _ := os.Open("out_expected.png")
	defer expectedFile.Close()

	actualFile, _ = os.Open("out_actual.png")
	defer actualFile.Close()
	expected_im, _, _ := image.Decode(expectedFile)
	actual_im, _, _ := image.Decode(actualFile)

	assert.Equal(t, expected_im.Bounds(), im.Bounds())

	for x := 0; x < im.Bounds().Max.X; x++ {
		for y := 0; y < im.Bounds().Max.Y; y++ {
			er, eg, eb, ea := expected_im.At(x, y).RGBA()
			r, g, b, a := actual_im.At(x, y).RGBA()
			assert.Equal(t, er, r)
			assert.Equal(t, eg, g)
			assert.Equal(t, eb, b)
			assert.Equal(t, ea, a)
		}
	}

}
