package trace

import (
	"math"
	"math/rand"
)

// colors

var (
	black     vec3 = vec3{}
	white     vec3 = vec3{1, 1, 1}
	lightBlue vec3 = vec3{0.5, 0.7, 1}
)

// math

var (
	infinity float64 = math.Inf(1)
	pi       float64 = math.Pi
)

func randomFloat(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
