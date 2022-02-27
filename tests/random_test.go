package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/random"
	"github.com/stretchr/testify/assert"
)

func TestRandomNormalFloat(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := random.RandomNormalFloat()
		assert.True(t, r >= 0 && r < 1)
	}
}

func TestRandomFloat(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := random.RandomFloat(-2, 2)
		assert.True(t, r >= -2 && r < 2)
	}
}
