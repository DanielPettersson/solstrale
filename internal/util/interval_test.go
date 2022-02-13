package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineIntervals(t *testing.T) {
	interval := CombineIntervals(Interval{0, 2}, Interval{1, 3})
	assert.Equal(t, Interval{0, 3}, interval)

	interval = CombineIntervals(Interval{0, 1}, Interval{2, 3})
	assert.Equal(t, Interval{0, 3}, interval)

	interval = CombineIntervals(Interval{3, 3}, Interval{-1, -1})
	assert.Equal(t, Interval{-1, 3}, interval)
}

func TestContains(t *testing.T) {
	assert.False(t, Interval{-2, 2}.Contains(-3))
	assert.True(t, Interval{-2, 2}.Contains(-2))
	assert.True(t, Interval{-2, 2}.Contains(2))
	assert.False(t, Interval{-2, 2}.Contains(3))
}

func TestClamp(t *testing.T) {
	assert.Equal(t, float64(-2), Interval{-2, 2}.Clamp(-3))
	assert.Equal(t, float64(-2), Interval{-2, 2}.Clamp(-2))
	assert.Equal(t, float64(0), Interval{-2, 2}.Clamp(0))
	assert.Equal(t, float64(2), Interval{-2, 2}.Clamp(2))
	assert.Equal(t, float64(2), Interval{-2, 2}.Clamp(3))
}

func TestSize(t *testing.T) {
	assert.Equal(t, float64(0), Interval{}.Size())
	assert.Equal(t, float64(2), Interval{-1, 1}.Size())
	assert.Equal(t, float64(-2), Interval{1, -1}.Size())
}

func TestExpand(t *testing.T) {
	interval := Interval{-2, 2}

	assert.Equal(t, Interval{-3, 3}, interval.Expand(2))
	assert.Equal(t, Interval{-3.5, 3.5}, interval.Expand(3))
	assert.Equal(t, Interval{-1, 1}, interval.Expand(-2))
}

func TestAdd(t *testing.T) {
	interval := Interval{-2, 2}

	assert.Equal(t, Interval{0, 4}, interval.Add(2))
	assert.Equal(t, Interval{1, 5}, interval.Add(3))
	assert.Equal(t, Interval{-4, 0}, interval.Add(-2))
}
