package trace

import "math"

var (
	empty_interval    interval = interval{math.Inf(1), math.Inf(-1)}
	universe_interval interval = interval{math.Inf(-1), math.Inf(1)}
)

type interval struct {
	min float64
	max float64
}

func (i interval) contains(x float64) bool {
	return i.min <= x && x <= i.max
}

func (i interval) clamp(x float64) float64 {
	if x < i.min {
		return i.min
	}
	if x > i.max {
		return i.max
	}
	return x
}
