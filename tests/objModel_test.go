package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {

	o, err := hittable.NewObjModel("obj/", "missing.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to read obj file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestMissingMaterialFile(t *testing.T) {

	o, err := hittable.NewObjModel("obj/", "missingMaterialLib.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "Failed to read material file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestMissingImageFile(t *testing.T) {

	o, err := hittable.NewObjModel("obj/", "missingImage.obj")

	assert.Equal(t, nil, o)
	assert.Equal(t, "Failed to load image texture obj/missing.jpg. Got error: open obj/missing.jpg: no such file or directory", err.Error())
}

func TestInvalidImageFile(t *testing.T) {

	o, err := hittable.NewObjModel("obj/", "invalidImage.obj")

	assert.Equal(t, nil, o)
	assert.Equal(t, "Failed to decode image texture obj/invalidImage.mtl. Got error: image: unknown format", err.Error())
}
