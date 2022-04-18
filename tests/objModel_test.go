package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/hittable"
	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {

	o, err := hittable.NewObjModel("missing.obj")

	assert.Equal(t, nil, o)
	assert.Contains(t, err.Error(), "no such file or directory")
}
