package util

import (
	"math"

	"github.com/valyala/fastrand"
)

// math

var (
	Infinity   float64 = math.Inf(1)
	fastRandom fastrand.RNG
)

func SetRandomSeed(seed uint32) {
	fastRandom.Seed(seed)
}

func RandomNormalFloat() float64 {
	return float64(fastRandom.Uint32()) / float64(math.MaxUint32)
}

func RandomFloat(min float64, max float64) float64 {
	return RandomNormalFloat()*(max-min) + min
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
