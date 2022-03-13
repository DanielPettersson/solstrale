package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestBvhWithEmptyList(t *testing.T) {
	assert.Panics(t, func() {
		hittable.NewBoundingVolumeHierarchy(hittable.NewHittableList())
	})
}
