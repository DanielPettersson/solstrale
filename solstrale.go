package solstrale

import (
	"github.com/DanielPettersson/solstrale/camera"
	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/DanielPettersson/solstrale/internal/renderer"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/DanielPettersson/solstrale/spec"
)

func RayTrace(spec spec.TraceSpecification, output chan spec.TraceProgress, abort chan bool) {
	util.SetRandomSeed(uint32(spec.RandomSeed))

	world, camera, background := CornellBox(spec)

	renderer.Renderer{
		World:           world,
		Cam:             camera,
		BackgroundColor: background,
		Spec:            spec,
		Output:          output,
	}.Render(abort)

}

func FinalScene(spec spec.TraceSpecification) (hittable.Hittable, camera.Camera, geo.Vec3) {
	camera := camera.New(
		spec,
		40,
		0,
		100,
		geo.NewVec3(478, 278, -600),
		geo.NewVec3(278, 278, 0),
		geo.NewVec3(0, 1, 0),
	)

	boxes1 := hittable.NewHittableList()
	ground := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.48, 0.83, 0.53)}}

	boxesPerSide := 20.
	for i := .0; i < boxesPerSide; i++ {
		for j := .0; j < boxesPerSide; j++ {
			w := 100.0
			x0 := -1000 + i*w
			z0 := -1000 + j*w
			y0 := .0
			x1 := x0 + w
			y1 := util.RandomFloat(1, 101)
			z1 := z0 + w

			boxes1.Add(hittable.NewBox(geo.NewVec3(x0, y0, z0), geo.NewVec3(x1, y1, z1), ground))
		}
	}

	world := hittable.NewHittableList()
	world.Add(hittable.NewBoundingVolumeHierarchy(boxes1))

	light := material.DiffuseLight{Emit: material.SolidColor{ColorValue: geo.NewVec3(7, 7, 7)}}
	world.Add(hittable.NewQuad(geo.NewVec3(123, 554, 147), geo.NewVec3(300, 0, 0), geo.NewVec3(0, 0, 265), light))

	movingSphereMaterial := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.7, 0.3, 0.1)}}
	world.Add(hittable.NewMotionBlur(hittable.NewSphere(geo.NewVec3(400, 400, 200), 50, movingSphereMaterial), geo.NewVec3(30, 0, 0)))

	glass := material.Dielectric{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}, IndexOfRefraction: 1.5}

	world.Add(hittable.NewSphere(geo.NewVec3(260, 150, 45), 50, glass))
	world.Add(hittable.NewSphere(geo.NewVec3(0, 150, 145), 50, material.Metal{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.8, 0.8, 0.9)}, Fuzz: 1}))

	boundary := hittable.NewSphere(geo.NewVec3(360, 150, 145), 70, glass)
	world.Add(boundary)
	world.Add(hittable.NewConstantMedium(boundary, 0.2, material.SolidColor{ColorValue: geo.NewVec3(0.2, 0.4, 0.9)}))
	boundary = hittable.NewSphere(geo.NewVec3(0, 0, 0), 5000, glass)
	world.Add(hittable.NewConstantMedium(boundary, 0.00013, material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}))

	// world.Add(hittable.NewSphere(geo.NewVec3(400, 200, 400), 100, material.Lambertian{imageTexture{}}))
	noiseTexture := material.NoiseTexture{ColorValue: geo.NewVec3(1, 1, 1), Scale: .1}
	world.Add(hittable.NewSphere(geo.NewVec3(220, 280, 300), 80, material.Lambertian{Tex: noiseTexture}))

	boxes2 := hittable.NewHittableList()
	white := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.73, 0.73, 0.73)}}
	for j := 0; j < 1000; j++ {
		boxes2.Add(hittable.NewSphere(geo.RandomVec3(0, 165), 10, white))
	}

	world.Add(hittable.NewTranslation(hittable.NewRotationY(hittable.NewBoundingVolumeHierarchy(boxes2), 15), geo.NewVec3(-100, 270, 395)))

	return &world, camera, geo.NewVec3(0, 0, 0)
}

func CornellBox(spec spec.TraceSpecification) (hittable.Hittable, camera.Camera, geo.Vec3) {
	camera := camera.New(
		spec,
		40,
		20,
		1070,
		geo.NewVec3(278, 278, -800),
		geo.NewVec3(278, 278, 0),
		geo.NewVec3(0, 1, 0),
	)

	red := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(.65, .05, .05)}}
	white := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(.73, .73, .73)}}
	green := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(.12, .45, .15)}}
	light := material.DiffuseLight{Emit: material.SolidColor{ColorValue: geo.NewVec3(15, 15, 15)}}

	world := hittable.NewHittableList()

	world.Add(hittable.NewQuad(geo.NewVec3(555, 0, 0), geo.NewVec3(0, 555, 0), geo.NewVec3(0, 0, 555), green))
	world.Add(hittable.NewQuad(geo.NewVec3(0, 0, 0), geo.NewVec3(0, 555, 0), geo.NewVec3(0, 0, 555), red))
	world.Add(hittable.NewQuad(geo.NewVec3(408, 554, 383), geo.NewVec3(-260, 0, 0), geo.NewVec3(0, 0, -210), light))
	world.Add(hittable.NewQuad(geo.NewVec3(0, 0, 0), geo.NewVec3(555, 0, 0), geo.NewVec3(0, 0, 555), white))
	world.Add(hittable.NewQuad(geo.NewVec3(555, 555, 555), geo.NewVec3(-555, 0, 0), geo.NewVec3(0, 0, -555), white))
	world.Add(hittable.NewQuad(geo.NewVec3(0, 0, 555), geo.NewVec3(555, 0, 0), geo.NewVec3(0, 555, 0), white))

	box1 := hittable.NewBox(geo.NewVec3(0, 0, 0), geo.NewVec3(165, 330, 165), white)
	box1 = hittable.NewRotationY(box1, 15)
	box1 = hittable.NewTranslation(box1, geo.NewVec3(265, 0, 295))
	world.Add(box1)

	box2 := hittable.NewBox(geo.NewVec3(0, 0, 0), geo.NewVec3(165, 165, 165), white)
	box2 = hittable.NewRotationY(box2, -18)
	box2 = hittable.NewTranslation(box2, geo.NewVec3(130, 0, 65))
	world.Add(box2)

	return &world, camera, geo.NewVec3(0, 0, 0)
}

func RandomSpheres(spec spec.TraceSpecification) (hittable.Hittable, camera.Camera, geo.Vec3) {

	camera := camera.New(
		spec,
		20,
		0.1,
		10,
		geo.NewVec3(13, 2, 3),
		geo.NewVec3(0, 0, 0),
		geo.NewVec3(0, 1, 0),
	)

	world := hittable.NewHittableList()

	groundMaterial := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.2, 0.3, 0.1)}}
	world.Add(hittable.NewSphere(geo.NewVec3(0, -1000, 0), 1000, groundMaterial))

	spheres := hittable.NewHittableList()
	for a := -7.0; a < 7; a++ {
		for b := -5.0; b < 5; b++ {
			chooseMat := util.RandomNormalFloat()
			center := geo.NewVec3(a+0.9*util.RandomNormalFloat(), 0.2, b+0.9*util.RandomNormalFloat())

			if center.Sub(geo.NewVec3(4, 0.2, 0)).Length() > 0.9 {

				if chooseMat < 0.8 {
					material := material.Lambertian{Tex: material.SolidColor{ColorValue: geo.RandomVec3(0, 1).Mul(geo.RandomVec3(0, 1))}}
					sphere := hittable.NewSphere(center, 0.2, material)
					blur := hittable.NewMotionBlur(sphere, geo.NewVec3(0, util.RandomFloat(0, 0.5), 0))
					spheres.Add(blur)
				} else if chooseMat < 0.95 {
					material := material.Metal{Tex: material.SolidColor{ColorValue: geo.RandomVec3(0.5, 1)}, Fuzz: util.RandomFloat(0, 0.5)}
					spheres.Add(hittable.NewSphere(center, 0.2, material))
				} else {
					material := material.Dielectric{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}, IndexOfRefraction: 1.5}
					spheres.Add(hittable.NewSphere(center, 0.2, material))
				}

			}
		}
	}

	spheres.Add(hittable.NewSphere(geo.NewVec3(0, 1, 0), 1.0, material.Dielectric{Tex: material.SolidColor{ColorValue: geo.NewVec3(1, 1, 1)}, IndexOfRefraction: 1.5}))
	spheres.Add(hittable.NewSphere(geo.NewVec3(-4, 1, 0), 1, material.Lambertian{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.4, 0.2, 0.1)}}))
	spheres.Add(hittable.NewSphere(geo.NewVec3(4, 1, 0), 1.0, material.Metal{Tex: material.SolidColor{ColorValue: geo.NewVec3(0.7, 0.6, 0.5)}, Fuzz: 0}))
	world.Add(hittable.NewBoundingVolumeHierarchy(spheres))

	return &world, camera, geo.NewVec3(0.70, 0.80, 1.00)
}
