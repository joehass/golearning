package binary

import (
	"encoding/binary"
	"fmt"
	"testing"
)

/**
转换有两种方式，也就是大端和小端。
大端就是内存中低地址对应着整数的高位，所以0123的顺序平成int32，整数最高8位是0，接着是1，依次类推，所以是66051
小端就是反过来，最高8位是3，也就是00000101，就是50462976

在计算机内部，小端序被广泛应用于cpu内部存储数据，而在其他场景譬如网络传输和文件存储使用大端序
*/

func TestB(t *testing.T) {
	var a []byte = []byte{0, 1, 2, 3}
	fmt.Println(a)
	fmt.Println(binary.BigEndian.Uint32(a))
	fmt.Println(binary.LittleEndian.Uint32(a))
}
