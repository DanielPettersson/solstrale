package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/material"
	"github.com/stretchr/testify/assert"
)

func TestNonPdfGeneratingMaterial(t *testing.T) {
	m := material.NonPdfGeneratingMaterial{}
	r := geo.NewRay(
		geo.RandomVec3(-1, 1),
		geo.RandomVec3(-1, 1),
		0,
	)

	assert.Equal(t, 0., m.ScatteringPdf(nil, r))
}

func TestNonLightEmittingMaterial(t *testing.T) {
	m := material.NonLightEmittingMaterial{}
	assert.Equal(t, geo.ZeroVector, m.Emitted(nil))
}
