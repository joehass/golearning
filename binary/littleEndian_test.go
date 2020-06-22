package binary

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestLittle(t *testing.T) {
	blob := []byte("byte")

	binary.LittleEndian.PutUint32(blob, uint32(len(blob)))

	fmt.Println(string(blob))

	binary.LittleEndian.Uint32(blob)

}
