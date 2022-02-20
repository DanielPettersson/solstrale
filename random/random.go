// Package random provides a random number generator to be used by ray tracer.
// Currently depends on fastrand implementation, for a good enough random the is really fast.
package random

import (
	"math"

	"github.com/valyala/fastrand"
)

// Random number generator
type Random struct {
	fastRandom fastrand.RNG
}

// Creates a new random number generator with given seed
// A seed value of 0 is random seed
func NewRandom(seed uint32) Random {
	var fr fastrand.RNG
	fr.Seed(seed)

	return Random{
		fastRandom: fr,
	}
}

// RandomNormalFloat returns a random uint32
func (r *Random) RandomUint32() uint32 {
	return r.fastRandom.Uint32()
}

// RandomNormalFloat returns a random float 0 to <1
func (r *Random) RandomNormalFloat() float64 {
	return float64(r.fastRandom.Uint32()) / float64(math.MaxUint32)
}

// RandomFloat returns a random float min to <max
func (r *Random) RandomFloat(min float64, max float64) float64 {
	return r.RandomNormalFloat()*(max-min) + min
}
