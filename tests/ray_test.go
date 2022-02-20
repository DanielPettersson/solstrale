package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestRayAt(t *testing.T) {

	rand := util.NewRandom(0)
	origin := geo.NewVec3(1, 2, 3)
	direction := geo.NewVec3(4, 5, 6)

	r := geo.Ray{
		Origin:    origin,
		Direction: direction,
		Time:      rand.RandomNormalFloat(),
	}

	assert.Equal(t, r.At(0), origin)
	assert.Equal(t, r.At(1), origin.Add(direction))
	assert.Equal(t, r.At(-1), origin.Sub(direction))
	assert.Equal(t, r.At(3), geo.NewVec3(13, 17, 21))
}
