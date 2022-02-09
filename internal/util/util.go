package util

import (
	"math"

	"github.com/valyala/fastrand"
)

var (
	// Infinity is the largest possible float64 value
	Infinity float64 = math.Inf(1)
	// AlmostZero is a value that is so small as to be almost zero
	AlmostZero float64 = 1e-8
	fastRandom fastrand.RNG
)

// SetRandomSeed set the seed to be used for following random operations
func SetRandomSeed(seed uint32) {
	fastRandom.Seed(seed)
}

// RandomNormalFloat returns a random float 0 to <1
func RandomNormalFloat() float64 {
	return float64(fastRandom.Uint32()) / float64(math.MaxUint32)
}

// RandomFloat returns a random float min to <max
func RandomFloat(min float64, max float64) float64 {
	return RandomNormalFloat()*(max-min) + min
}

// DegreesToRadians converts an angle in degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
