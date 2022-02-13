package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRandomSeed(t *testing.T) {
	SetRandomSeed(123)
	r := RandomNormalFloat()
	assert.Equal(t, 0.007376669907797284, r)
}

func TestRandomNormalFloat(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomNormalFloat()
		assert.True(t, r >= 0 && r < 1)
	}
}

func TestRandomFloat(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomFloat(-2, 2)
		assert.True(t, r >= -2 && r < 2)
	}
}

func TestDegreesToRadians(t *testing.T) {
	r := DegreesToRadians(0)
	assert.Equal(t, r, float64(0))

	r = DegreesToRadians(180)
	assert.Equal(t, r, math.Pi)

	r = DegreesToRadians(360)
	assert.Equal(t, r, math.Pi*2)

	r = DegreesToRadians(-180)
	assert.Equal(t, r, -math.Pi)
}
