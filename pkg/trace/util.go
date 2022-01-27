package trace

import (
	"math"

	"github.com/valyala/fastrand"
)

// colors

var (
	black     vec3 = vec3{}
	white     vec3 = vec3{1, 1, 1}
	lightBlue vec3 = vec3{0.5, 0.7, 1}
)

// math

var (
	infinity   float64 = math.Inf(1)
	fastRandom fastrand.RNG
)

func setRandomSeed(seed uint32) {
	fastRandom.Seed(seed)
}

func randomNormalFloat() float64 {
	return float64(fastRandom.Uint32()) / float64(math.MaxUint32)
}

func randomFloat(min float64, max float64) float64 {
	return randomNormalFloat()*(max-min) + min
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
