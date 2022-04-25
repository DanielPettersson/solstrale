package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/geo"
	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestBvhWithEmptyList(t *testing.T) {
	assert.Panics(t, func() {
		hittable.NewBoundingVolumeHierarchy([]hittable.Hittable{})
	})
}

func TestSortHittablesByCenter(t *testing.T) {

	hittables := []hittable.Hittable{
		hittable.NewSphere(geo.NewVec3(7, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(5, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(3, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(1, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(-1, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(-3, 0, 0), 1, nil),
	}

	mid := hittable.SortHittablesByCenter(hittables, 2, 0)
	assert.Equal(t, 3, mid)
	assert.Equal(t, -3., hittables[0].Center().X)
	assert.Equal(t, -1., hittables[1].Center().X)
	assert.Equal(t, 1., hittables[2].Center().X)
	assert.Equal(t, 3., hittables[3].Center().X)
	assert.Equal(t, 5., hittables[4].Center().X)
	assert.Equal(t, 7., hittables[5].Center().X)
}

func TestSortHittablesByCenterOnlyOneOff(t *testing.T) {

	hittables := []hittable.Hittable{
		hittable.NewSphere(geo.NewVec3(-3, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(5, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(1, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(3, 0, 0), 1, nil),
		hittable.NewSphere(geo.NewVec3(-1, 0, 0), 1, nil),
	}

	mid := hittable.SortHittablesByCenter(hittables, 2, 0)
	assert.Equal(t, 3, mid)
	assert.Equal(t, -3., hittables[0].Center().X)
	assert.Equal(t, -1., hittables[1].Center().X)
	assert.Equal(t, 1., hittables[2].Center().X)
	assert.Equal(t, 3., hittables[3].Center().X)
	assert.Equal(t, 5., hittables[4].Center().X)
}
