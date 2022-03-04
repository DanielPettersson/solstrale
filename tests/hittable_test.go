package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestNonPdfUsingHittable(t *testing.T) {

	h := hittable.NonPdfUsingHittable{}

	assert.Panics(t, func() {
		h.PdfValue(geo.RandomVec3(-1, 1), geo.RandomVec3(-1, 1))
	})

	assert.Panics(t, func() {
		h.RandomDirection(geo.RandomVec3(-1, 1))
	})
}
