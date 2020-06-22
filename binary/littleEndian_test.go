package binary

import (
	"encoding/binary"
	"testing"
)

func TestLittle(t *testing.T) {
	blob := []byte("byte")

	blob1 := []byte("byte1")

	binary.LittleEndian.PutUint64(blob, blob1)

}
