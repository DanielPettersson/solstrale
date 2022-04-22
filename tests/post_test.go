package tests

import (
	"testing"

	"github.com/DanielPettersson/solstrale/post"
	"github.com/stretchr/testify/assert"
)

func TestNewOidnExecutableNotFound(t *testing.T) {
	oidnPost, err := post.NewOidn("notFound")

	assert.Nil(t, oidnPost)
	assert.Equal(t, "Oidn path does not exist: stat notFound: no such file or directory", err.Error())
}

func TestNewOidnExecutableNotExecutable(t *testing.T) {
	oidnPost, err := post.NewOidn("post_test.go")

	assert.Nil(t, oidnPost)
	assert.Equal(t, "Oidn path is not executable: post_test.go", err.Error())
}

func TestNewOidnExecutable(t *testing.T) {
	oidnPost, err := post.NewOidn("mock_oidn.sh")

	assert.NotNil(t, oidnPost)
	assert.Nil(t, err)
}
