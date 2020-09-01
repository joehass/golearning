package minKoa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//实现一个最小的洋葱模型
func TestMid1(t *testing.T) {
	str := ""
	e := New()
	e.Use(func() {
		str += "a"
		e.Next()
		str += "c"
	})

	e.Use(func() {
		str += "b"
		e.Next()
		str += "d"
	})

	e.Run()
	assert.Equal(t, "abdc", str)
}
