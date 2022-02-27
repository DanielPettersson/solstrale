// Package random provides a random number generator to be used by ray tracer.
// Currently is using https://github.com/valyala/fastrand, for a good enough random that is really fast.
package random

import (
	"math"

	"github.com/valyala/fastrand"
)

// RandomNormalFloat returns a random float 0 to <1
func RandomNormalFloat() float64 {
	return float64(fastrand.Uint32()) / float64(math.MaxUint32)
}

// RandomFloat returns a random float min to <max
func RandomFloat(min float64, max float64) float64 {
	return RandomNormalFloat()*(max-min) + min
}
