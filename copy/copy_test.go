package copy

import (
	"fmt"
	"testing"
)

func TestC(t *testing.T) {
	//设置元素数量为1000
	const elementCount = 1000

	//预分配足够多的元素切片
	srcData := make([]int, elementCount)

	//将切片赋值
	for i := 0; i < elementCount; i++ {
		srcData[i] = i
	}

	//引用切片数据
	refData := srcData

	//预分配足够多的元素切片
	copyData := make([]int, elementCount)
	//将数据复制到新的切片空间中
	copy(copyData, srcData)

	//修改原始数据的第一个元素
	srcData[0] = 999

	//打印引用切片的第一个元素
	fmt.Println(refData[0])

	//打印复制切片的第一个和最后一个元素
	fmt.Println(copyData[0], copyData[elementCount-1])

	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", copyData[i])
	}

	fmt.Println(" ")

	//复制原始数据从4到6，替换到copyData的第1，2个位置
	copy(copyData, srcData[4:6])

	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", copyData[i])
	}
}

func TestC2(t *testing.T) {
	const elementCount = 10

	srcData := make([]int, elementCount)
	for i := 0; i < elementCount; i++ {
		srcData[i] = i
	}

	//左闭右开区间
	fmt.Println(srcData[1:])

	fmt.Println(srcData[:elementCount])

	ref := srcData

	copyData := make([]int, 20)

	n := copy(copyData, ref[1:3])

	fmt.Println(copyData)
	fmt.Println(n)
	fmt.Println(len(copyData))
	fmt.Println(len(ref[1:3]))

	n2 := copy(copyData[1:], ref[4:7])

	fmt.Println(copyData)
	fmt.Println(n2)
	fmt.Println(len(copyData[1:]))
	fmt.Println(len(ref[4:7]))

}

func TestC3(t *testing.T) {
	srcData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	desData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	copy(srcData[1:3], desData[2:6])

	fmt.Println(srcData)
}
