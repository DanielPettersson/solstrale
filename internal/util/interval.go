package util

import "math"

var (
	// EmptyInterval contains nothing
	EmptyInterval Interval = Interval{math.Inf(1), math.Inf(-1)}
	// UniverseInterval contains everything
	UniverseInterval Interval = Interval{math.Inf(-1), math.Inf(1)}
)

// Interval defines a range between Min and Max inclusive
type Interval struct {
	Min float64
	Max float64
}

// CombineIntervals creates a new interval that is the union of the two given.
// If there is a gap between the intervals, that is included in the returned interval.
func CombineIntervals(a Interval, b Interval) Interval {
	return Interval{
		math.Min(a.Min, b.Min),
		math.Max(a.Max, b.Max),
	}
}

// Contains checks if the interval contains a given value
func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

// Clamp returns the given value clamped to the interval
func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

// Size return the size of the interval
func (i Interval) Size() float64 {
	return i.Max - i.Min
}

// Expand returns a new interval that is larger by given value delta.
// Delta is added equally to both sides of the interval
func (i Interval) Expand(delta float64) Interval {
	padding := delta / 2
	return Interval{
		i.Min - padding,
		i.Max + padding,
	}
}

// Add returns a new interval that is increased with given value.
// The returned inteval has same size as original
func (i Interval) Add(displacement float64) Interval {
	return Interval{i.Min + displacement, i.Max + displacement}
}
