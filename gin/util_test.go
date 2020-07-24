package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, "", joinPaths("", ""))
	assert.Equal(t, "/", joinPaths("", "/"))
	assert.Equal(t, "/a", joinPaths("/a", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", "/"))
	assert.Equal(t, "/a/", joinPaths("/a", "/"))
	assert.Equal(t, "/a/hola", joinPaths("/a", "/hola"))
	assert.Equal(t, "/a/hola", joinPaths("/a/", "/hola"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola/"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola//"))
}
