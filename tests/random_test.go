package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/random"
	"github.com/stretchr/testify/assert"
)

func TestSetRandomSeed(t *testing.T) {
	rand := random.NewRandom(123)
	r := rand.RandomNormalFloat()
	assert.Equal(t, 0.007376669907797284, r)
}

func TestRandomNormalFloat(t *testing.T) {
	rand := random.NewRandom(0)
	for i := 0; i < 100; i++ {
		r := rand.RandomNormalFloat()
		assert.True(t, r >= 0 && r < 1)
	}
}

func TestRandomFloat(t *testing.T) {
	rand := random.NewRandom(0)
	for i := 0; i < 100; i++ {
		r := rand.RandomFloat(-2, 2)
		assert.True(t, r >= -2 && r < 2)
	}
}
