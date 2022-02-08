package util

import "math"

var (
	EmptyInterval    Interval = Interval{math.Inf(1), math.Inf(-1)}
	UniverseInterval Interval = Interval{math.Inf(-1), math.Inf(1)}
)

type Interval struct {
	Min float64
	Max float64
}

func CombineIntervals(a Interval, b Interval) Interval {
	return Interval{
		math.Min(a.Min, b.Min),
		math.Max(a.Max, b.Max),
	}
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Expand(delta float64) Interval {
	padding := delta / 2
	return Interval{
		i.Min - padding,
		i.Max + padding,
	}
}

func (i Interval) Add(displacement float64) Interval {
	return Interval{i.Min + displacement, i.Max + displacement}
}
