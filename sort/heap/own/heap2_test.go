package own

import (
	"fmt"
	"testing"
)

/**
自己实现的堆排序
*/

var (
	heap = []int{100, 16, 4, 8, 70, 2, 36, 22, 5, 12}
)

func TestH21(t *testing.T) {
	MakeHeap()
	fmt.Println("\n构建树后:")
	Print(heap)
	fmt.Println("\n增加 90,30,1 :")
	Push(90)
	Push(30)
	Push(1)
	Print(heap)

	n := Pop()
	fmt.Println("\nPop出最小值(", n, ")后:")
	Print(heap)

	fmt.Println("\nRemove()掉idx为3即值", heap[3-1], "后:")
	Remove(3)
	Print(heap)

	fmt.Println("\nHeapSort()后:")
	HeapSort()
	Print(heap)
}

//构建堆
func MakeHeap() {
	n := len(heap)
	//从非叶子节点开始（叶节点不用调整，第一个非叶子节点：arr.length/2-1），从左至右，从下至上进行调整
	for i := n/2 - 1; i >= 0; i-- {
		down(i, n)
	}
}

//由父节点至子节点依次建堆
//parent:i
//left child:2*i+1
//right child:2*i+2
func down(i, n int) {
	//构建最小堆，父小于两子节点值
	for {
		j1 := 2*i + 1 //左节点
		//j1 <0 溢出
		if j1 >= n || j1 < 0 {
			break
		}
		//找到两个节点中最小的
		j := j1
		//j2 <n 是为了pop时，不比较最后一个元素
		if j2 := j1 + 1; j2 < n && !Less(j1, j2) { //如果右节点大于左节点
			j = j2 // = 2*i + 2 //右节点
		}

		//然后和父节点比较，如果父节点小于这个子节点最小值，则break，否则swap
		if !Less(j, i) {
			break
		}

		Swap(i, j)
		i = j
	}
}

func Push(x interface{}) {
	heap = append(heap, x.(int))
	up(len(heap) - 1)
	return
}

func up(j int) {
	for {
		i := (j - 1) / 2 //parent，第一个非叶节点
		//i==j:表示没有数据，插入的是第一个元素
		//父节点小于子节点，符合最小堆条件
		if i == j || !Less(j, i) {
			break
		}
		//子节点比父节点小，互换
		Swap(i, j)
		j = i
	}
}

func Pop() interface{} {
	n := len(heap) - 1
	//取出最小元素后，需要将最后一个元素和第一个元素进行调换
	Swap(0, n)
	down(0, n)
	old := heap

	n = len(old)
	x := old[n-1]
	heap = old[0 : n-1]
	return x
}

func Swap(a, b int) {
	heap[a], heap[b] = heap[b], heap[a]
}

func Less(a, b int) bool {
	return heap[a] < heap[b]
}

func Remove(i int) interface{} {
	n := len(heap) - 1
	if n != i {
		Swap(i, n)
		down(i, n)
		up(i)
	}
	return Pop()
}

func HeapSort() {
	//升序 Less(heap[a] > heap[b])	//最大堆
	//降序 Less(heap[a] < heap[b])	//最小堆
	for i := len(heap) - 1; i > 0; i-- {
		//移除顶部元素到数组末尾，然后剩下的重建堆，依次循环
		Swap(0, i)
		down(0, i)
	}
}

func Print(arr []int) {
	for _, v := range arr {
		fmt.Printf("%d ", v)
	}
}
