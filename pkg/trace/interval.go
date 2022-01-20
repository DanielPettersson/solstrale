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

func combineIntervals(a interval, b interval) interval {
	return interval{
		math.Min(a.min, b.min),
		math.Max(a.max, b.max),
	}
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

func (i interval) size() float64 {
	return i.max - i.min
}

func (i interval) expand(delta float64) interval {
	padding := delta / 2
	return interval{
		i.min - padding,
		i.max + padding,
	}
}

func (i interval) add(displacement float64) interval {
	return interval{i.min + displacement, i.max + displacement}
}
