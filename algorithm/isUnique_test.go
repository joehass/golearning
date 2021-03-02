package algorithm

import (
	"fmt"
	"testing"
)

func isUnique(astr string) bool {
	mark := 0
	for i := range astr {
		off := int(astr[i]) - int('a')
		markOff := 1 << off
		if markOff&mark == 0 {
			mark = markOff | mark
		} else {
			return false
		}
	}

	return true
}

func TestM1(t *testing.T) {
	fmt.Println(isUnique("abc"))
}
