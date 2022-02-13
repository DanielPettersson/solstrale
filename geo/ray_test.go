package geo

import (
	"testing"

	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestAt(t *testing.T) {

	origin := NewVec3(1, 2, 3)
	direction := NewVec3(4, 5, 6)

	r := Ray{
		Origin:    origin,
		Direction: direction,
		Time:      util.RandomNormalFloat(),
	}

	assert.Equal(t, r.At(0), origin)
	assert.Equal(t, r.At(1), origin.Add(direction))
	assert.Equal(t, r.At(-1), origin.Sub(direction))
	assert.Equal(t, r.At(3), NewVec3(13, 17, 21))
}
