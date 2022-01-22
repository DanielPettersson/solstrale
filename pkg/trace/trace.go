package trace

import (
	"math/rand"
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

func RayTrace(spec TraceSpecification, output chan TraceProgress) {
	rand.Seed(int64(spec.RandomSeed))

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

	groundMaterial := lambertian{vec3{0.5, 0.5, 0.5}}
	world.add(createSphere(vec3{0, -1000, 0}, 1000, groundMaterial))

	spheres := emptyHittableList()
	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			chooseMat := rand.Float64()
			center := vec3{a + 0.9*rand.Float64(), 0.2, b + 0.9*rand.Float64()}

			if center.sub(vec3{4, 0.2, 0}).length() > 0.9 {

				if chooseMat < 0.8 {
					material := lambertian{randomVec3(0, 1).mul(randomVec3(0, 1))}
					sphere := createSphere(center, 0.2, material)
					blur := createMotionBlur(sphere, vec3{0, randomFloat(0, 0.5), 0})
					spheres.add(blur)
				} else if chooseMat < 0.95 {
					material := metal{randomVec3(0.5, 1), randomFloat(0, 0.5)}
					spheres.add(createSphere(center, 0.2, material))
				} else {
					material := dielectric{vec3{1, 1, 1}, 1.5}
					spheres.add(createSphere(center, 0.2, material))
				}

			}
		}
	}

	spheres.add(createSphere(vec3{0, 1, 0}, 1.0, dielectric{white, 1.5}))
	spheres.add(createSphere(vec3{-4, 1, 0}, 1, lambertian{vec3{0.4, 0.2, 0.1}}))
	spheres.add(createSphere(vec3{4, 1, 0}, 1.0, metal{vec3{0.7, 0.6, 0.5}, 0}))
	world.add(createBvh(spheres))

	scene{
		world:  world,
		cam:    camera,
		spec:   spec,
		output: output,
	}.render()

}
