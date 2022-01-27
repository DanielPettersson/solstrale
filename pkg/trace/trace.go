package trace

import (
	"github.com/ojrac/opensimplex-go"
)

var (
	textureData map[string]imageTexture = map[string]imageTexture{}
)

type TraceSpecification struct {
	ImageWidth      int
	ImageHeight     int
	DrawOffsetX     int
	DrawOffsetY     int
	DrawWidth       int
	DrawHeight      int
	SamplesPerPixel int
	RandomSeed      int
}

type TraceProgress struct {
	Progress      float64
	Specification TraceSpecification
	ImageData     []byte
}

func AddTexture(name string, width, height int, bytes []byte) {
	textureData[name] = imageTexture{
		bytes:  bytes,
		width:  width,
		height: height,
		mirror: false,
	}
}

func RayTrace(spec TraceSpecification, output chan TraceProgress) {
	setRandomSeed(uint32(spec.RandomSeed))

	world, camera, background := cornellBox(spec)

	scene{
		world:           world,
		cam:             camera,
		backgroundColor: background,
		spec:            spec,
		output:          output,
	}.render()

}

func finalScene(spec TraceSpecification) (hittableList, camera, vec3) {
	camera := createCamera(
		spec,
		40,
		0,
		100,
		vec3{478, 278, -600},
		vec3{278, 278, 0},
		vec3{0, 1, 0},
	)

	boxes1 := emptyHittableList()
	ground := lambertian{solidColor{vec3{0.48, 0.83, 0.53}}}

	boxesPerSide := 20.
	for i := .0; i < boxesPerSide; i++ {
		for j := .0; j < boxesPerSide; j++ {
			w := 100.0
			x0 := -1000 + i*w
			z0 := -1000 + j*w
			y0 := .0
			x1 := x0 + w
			y1 := randomFloat(1, 101)
			z1 := z0 + w

			boxes1.add(createBox(vec3{x0, y0, z0}, vec3{x1, y1, z1}, ground))
		}
	}

	world := emptyHittableList()
	world.add(createBvh(boxes1))

	light := diffuseLight{solidColor{vec3{7, 7, 7}}}
	world.add(createQuad(vec3{123, 554, 147}, vec3{300, 0, 0}, vec3{0, 0, 265}, light))

	movingSphereMaterial := lambertian{solidColor{vec3{0.7, 0.3, 0.1}}}
	world.add(createMotionBlur(createSphere(vec3{400, 400, 200}, 50, movingSphereMaterial), vec3{30, 0, 0}))

	glass := dielectric{solidColor{vec3{1, 1, 1}}, 1.5}

	world.add(createSphere(vec3{260, 150, 45}, 50, glass))
	world.add(createSphere(vec3{0, 150, 145}, 50, metal{solidColor{vec3{0.8, 0.8, 0.9}}, 1}))

	boundary := createSphere(vec3{360, 150, 145}, 70, glass)
	world.add(boundary)
	world.add(createConstantMedium(boundary, 0.2, solidColor{vec3{0.2, 0.4, 0.9}}))
	boundary = createSphere(vec3{0, 0, 0}, 5000, glass)
	world.add(createConstantMedium(boundary, 0.0001, solidColor{vec3{1, 1, 1}}))

	world.add(createSphere(vec3{400, 200, 400}, 100, lambertian{textureData["earth"]}))
	noiseTexture := noiseTexture{opensimplex.NewNormalized(int64(spec.RandomSeed)), vec3{1, 1, 1}, .1}
	world.add(createSphere(vec3{220, 280, 300}, 80, lambertian{noiseTexture}))

	boxes2 := emptyHittableList()
	white := lambertian{solidColor{vec3{0.73, 0.73, 0.73}}}
	for j := 0; j < 1000; j++ {
		boxes2.add(createSphere(randomVec3(0, 165), 10, white))
	}

	world.add(createTranslation(createRotationY(createBvh(boxes2), 15), vec3{-100, 270, 395}))

	return world, camera, vec3{}
}

func cornellBox(spec TraceSpecification) (hittableList, camera, vec3) {
	camera := createCamera(
		spec,
		40,
		20,
		1070,
		vec3{278, 278, -800},
		vec3{278, 278, 0},
		vec3{0, 1, 0},
	)

	red := lambertian{solidColor{vec3{.65, .05, .05}}}
	white := lambertian{solidColor{vec3{.73, .73, .73}}}
	green := lambertian{solidColor{vec3{.12, .45, .15}}}
	light := diffuseLight{solidColor{vec3{15, 15, 15}}}

	world := emptyHittableList()

	world.add(createQuad(vec3{555, 0, 0}, vec3{0, 555, 0}, vec3{0, 0, 555}, green))
	world.add(createQuad(vec3{0, 0, 0}, vec3{0, 555, 0}, vec3{0, 0, 555}, red))
	world.add(createQuad(vec3{408, 554, 383}, vec3{-260, 0, 0}, vec3{0, 0, -210}, light))
	world.add(createQuad(vec3{0, 0, 0}, vec3{555, 0, 0}, vec3{0, 0, 555}, white))
	world.add(createQuad(vec3{555, 555, 555}, vec3{-555, 0, 0}, vec3{0, 0, -555}, white))
	world.add(createQuad(vec3{0, 0, 555}, vec3{555, 0, 0}, vec3{0, 555, 0}, white))

	return world, camera, vec3{}
}

func randomSpheres(spec TraceSpecification) (hittableList, camera, vec3) {

	camera := createCamera(
		spec,
		20,
		0.1,
		10,
		vec3{13, 2, 3},
		vec3{0, 0, 0},
		vec3{0, 1, 0},
	)

	world := emptyHittableList()

	groundMaterial := lambertian{solidColor{vec3{0.2, 0.3, 0.1}}}
	world.add(createSphere(vec3{0, -1000, 0}, 1000, groundMaterial))

	spheres := emptyHittableList()
	for a := -7.0; a < 7; a++ {
		for b := -5.0; b < 5; b++ {
			chooseMat := randomNormalFloat()
			center := vec3{a + 0.9*randomNormalFloat(), 0.2, b + 0.9*randomNormalFloat()}

			if center.sub(vec3{4, 0.2, 0}).length() > 0.9 {

				if chooseMat < 0.8 {
					material := lambertian{solidColor{randomVec3(0, 1).mul(randomVec3(0, 1))}}
					sphere := createSphere(center, 0.2, material)
					blur := createMotionBlur(sphere, vec3{0, randomFloat(0, 0.5), 0})
					spheres.add(blur)
				} else if chooseMat < 0.95 {
					material := metal{solidColor{randomVec3(0.5, 1)}, randomFloat(0, 0.5)}
					spheres.add(createSphere(center, 0.2, material))
				} else {
					material := dielectric{solidColor{vec3{1, 1, 1}}, 1.5}
					spheres.add(createSphere(center, 0.2, material))
				}

			}
		}
	}

	spheres.add(createSphere(vec3{0, 1, 0}, 1.0, dielectric{solidColor{white}, 1.5}))
	spheres.add(createSphere(vec3{-4, 1, 0}, 1, lambertian{solidColor{vec3{0.4, 0.2, 0.1}}}))
	spheres.add(createSphere(vec3{4, 1, 0}, 1.0, metal{solidColor{vec3{0.7, 0.6, 0.5}}, 0}))
	world.add(createBvh(spheres))

	return world, camera, vec3{0.70, 0.80, 1.00}
}
