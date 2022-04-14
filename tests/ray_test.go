package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/random"
	"github.com/stretchr/testify/assert"
)

func TestRayAt(t *testing.T) {
	origin := geo.NewVec3(1, 2, 3)
	direction := geo.NewVec3(4, 5, 6)
	l := direction.Length()

	r := geo.NewRay(
		origin,
		direction,
		random.RandomNormalFloat(),
	)

	assert.Equal(t, r.At(0), origin)
	assert.True(t, r.At(l).Sub(origin.Add(direction)).NearZero())
	assert.True(t, r.At(-l).Sub(origin.Sub(direction)).NearZero())
	assert.True(t, r.At(l*3).Sub(geo.NewVec3(13, 17, 21)).NearZero())
}
