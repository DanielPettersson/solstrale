package trace

import (
	"math/rand"

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
	}
}

func RayTrace(spec TraceSpecification, output chan TraceProgress) {
	rand.Seed(int64(spec.RandomSeed))

	//world, camera := randomSpheres(spec)
	world, camera, background := cornellBox(spec)

	scene{
		world:           world,
		cam:             camera,
		backgroundColor: background,
		spec:            spec,
		output:          output,
	}.render()

}

func cornellBox(spec TraceSpecification) (hittableList, camera, vec3) {
	camera := createCamera(
		spec,
		40,
		0,
		20,
		vec3{278, 278, -800},
		vec3{278, 278, 0},
		vec3{0, 1, 0},
	)

	noiseTexture := noiseTexture{
		opensimplex.NewNormalized(int64(spec.RandomSeed)),
		vec3{1, 1, 1},
		0.1,
	}

	red := lambertian{solidColor{vec3{.65, .05, .05}}}
	white := lambertian{solidColor{vec3{.73, .73, .73}}}
	green := lambertian{solidColor{vec3{.12, .45, .15}}}
	light := diffuseLight{solidColor{vec3{15, 15, 15}}}
	noise := lambertian{noiseTexture}

	metal := metal{noiseTexture, 0.1}
	glass := dielectric{solidColor{vec3{1, 1, 1}}, 1.5}
	earth := lambertian{textureData["earth"]}

	world := emptyHittableList()

	world.add(createQuad(vec3{555, 0, 0}, vec3{0, 555, 0}, vec3{0, 0, 555}, green))
	world.add(createQuad(vec3{0, 0, 0}, vec3{0, 555, 0}, vec3{0, 0, 555}, red))
	world.add(createQuad(vec3{343, 554, 332}, vec3{-130, 0, 0}, vec3{0, 0, -105}, light))
	world.add(createQuad(vec3{0, 0, 0}, vec3{555, 0, 0}, vec3{0, 0, 555}, white))
	world.add(createQuad(vec3{555, 555, 555}, vec3{-555, 0, 0}, vec3{0, 0, -555}, white))
	world.add(createQuad(vec3{0, 0, 555}, vec3{555, 0, 0}, vec3{0, 555, 0}, white))

	box1 := createBox(vec3{0, 0, 0}, vec3{165, 330, 165}, noise)
	box1 = createRotationY(box1, 15)
	box1 = createTranslation(box1, vec3{265, 0, 295})
	world.add(box1)

	box2 := createBox(vec3{0, 0, 0}, vec3{165, 165, 165}, white)
	box2 = createRotationY(box2, -18)
	box2 = createTranslation(box2, vec3{130, 0, 65})
	world.add(box2)

	world.add(createSphere(vec3{400, 50, 200}, 50, metal))
	world.add(createSphere(vec3{190, 215, 147}, 50, glass))
	world.add(createSphere(vec3{270, 70, 270}, 70, earth))

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

	groundMaterial := lambertian{checkerTexture{0.32, solidColor{vec3{0.2, 0.3, 0.1}}, solidColor{vec3{0.9, 0.9, 0.9}}}}
	world.add(createSphere(vec3{0, -1000, 0}, 1000, groundMaterial))

	spheres := emptyHittableList()
	for a := -7.0; a < 7; a++ {
		for b := -5.0; b < 5; b++ {
			chooseMat := rand.Float64()
			center := vec3{a + 0.9*rand.Float64(), 0.2, b + 0.9*rand.Float64()}

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
