package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCombineIntervals(t *testing.T) {
	interval := util.CombineIntervals(util.Interval{Min: 0, Max: 2}, util.Interval{Min: 1, Max: 3})
	assert.Equal(t, util.Interval{Min: 0, Max: 3}, interval)

	interval = util.CombineIntervals(util.Interval{Min: 0, Max: 1}, util.Interval{Min: 2, Max: 3})
	assert.Equal(t, util.Interval{Min: 0, Max: 3}, interval)

	interval = util.CombineIntervals(util.Interval{Min: 3, Max: 3}, util.Interval{Min: -1, Max: -1})
	assert.Equal(t, util.Interval{Min: -1, Max: 3}, interval)
}

func TestIntervalContains(t *testing.T) {
	assert.False(t, util.Interval{Min: -2, Max: 2}.Contains(-3))
	assert.True(t, util.Interval{Min: -2, Max: 2}.Contains(-2))
	assert.True(t, util.Interval{Min: -2, Max: 2}.Contains(2))
	assert.False(t, util.Interval{Min: -2, Max: 2}.Contains(3))
}

func TestIntervalClamp(t *testing.T) {
	assert.Equal(t, float64(-2), util.Interval{Min: -2, Max: 2}.Clamp(-3))
	assert.Equal(t, float64(-2), util.Interval{Min: -2, Max: 2}.Clamp(-2))
	assert.Equal(t, float64(0), util.Interval{Min: -2, Max: 2}.Clamp(0))
	assert.Equal(t, float64(2), util.Interval{Min: -2, Max: 2}.Clamp(2))
	assert.Equal(t, float64(2), util.Interval{Min: -2, Max: 2}.Clamp(3))
}

func TestIntervalSize(t *testing.T) {
	assert.Equal(t, float64(0), util.Interval{}.Size())
	assert.Equal(t, float64(2), util.Interval{Min: -1, Max: 1}.Size())
	assert.Equal(t, float64(-2), util.Interval{Min: 1, Max: -1}.Size())
}

func TestIntervalExpand(t *testing.T) {
	interval := util.Interval{Min: -2, Max: 2}

	assert.Equal(t, util.Interval{Min: -3, Max: 3}, interval.Expand(2))
	assert.Equal(t, util.Interval{Min: -3.5, Max: 3.5}, interval.Expand(3))
	assert.Equal(t, util.Interval{Min: -1, Max: 1}, interval.Expand(-2))
}

func TestIntervalAdd(t *testing.T) {
	interval := util.Interval{Min: -2, Max: 2}

	assert.Equal(t, util.Interval{Min: 0, Max: 4}, interval.Add(2))
	assert.Equal(t, util.Interval{Min: 1, Max: 5}, interval.Add(3))
	assert.Equal(t, util.Interval{Min: -4, Max: 0}, interval.Add(-2))
}
