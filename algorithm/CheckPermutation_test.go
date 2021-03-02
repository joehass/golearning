package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func CheckPermutation(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	m := make(map[rune]int)

	for _, v := range s1 {
		m[v]++
	}
	for _, v := range s2 {
		m[v]--
	}

	for _, v := range m {
		if v > 0 {
			return false
		}
	}

	return true
}

func TestCheckPermutation(t *testing.T) {
	assert.Equal(t, true, CheckPermutation("abc", "bca"))
	assert.Equal(t, false, CheckPermutation("aab", "abb"))
}
