package tests

import (
	"math"
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestDegreesToRadians(t *testing.T) {
	r := util.DegreesToRadians(0)
	assert.Equal(t, r, float64(0))

	r = util.DegreesToRadians(180)
	assert.Equal(t, r, math.Pi)

	r = util.DegreesToRadians(360)
	assert.Equal(t, r, math.Pi*2)

	r = util.DegreesToRadians(-180)
	assert.Equal(t, r, -math.Pi)
}
