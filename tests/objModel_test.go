package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {

	o, err := hittable.NewObjModel("missing.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to read obj file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestMissingMaterialFile(t *testing.T) {

	o, err := hittable.NewObjModel("missingMaterialLib.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to read material file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestMissingImageFile(t *testing.T) {

	o, err := hittable.NewObjModel("missingImage.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to open image file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestInvalidImageFile(t *testing.T) {

	o, err := hittable.NewObjModel("invalidImage.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to read image file: image: unknown format")
}
