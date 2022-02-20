package tests

import (
	"math"
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestSetRandomSeed(t *testing.T) {
	rand := util.NewRandom(123)
	r := rand.RandomNormalFloat()
	assert.Equal(t, 0.007376669907797284, r)
}

func TestRandomNormalFloat(t *testing.T) {
	rand := util.NewRandom(0)
	for i := 0; i < 100; i++ {
		r := rand.RandomNormalFloat()
		assert.True(t, r >= 0 && r < 1)
	}
}

func TestRandomFloat(t *testing.T) {
	rand := util.NewRandom(0)
	for i := 0; i < 100; i++ {
		r := rand.RandomFloat(-2, 2)
		assert.True(t, r >= -2 && r < 2)
	}
}

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
