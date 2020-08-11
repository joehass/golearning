package heap

import (
	"container/heap"
	"fmt"
	"sort"
	"testing"
)

/**
go 中container/heap的使用
*/

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

/**
自定义的类，实现相关接口后，交由heap.Init()去构建堆，
从堆中pop()后，数据就被从heap中移除类，
升降序由less()来决定
自定义类也可以直接用sort来排序，因为实现了相关接口
*/
func TestH1(t *testing.T) {
	h := &IntHeap{100, 16, 4, 8, 70, 2, 36, 22, 5, 12}

	fmt.Println("\nHeap:")
	heap.Init(h)

	fmt.Printf("最小值：%d\n", (*h)[0])

	//依次输出最小值
	fmt.Println("\nHeap sort:")
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	//增加一个新值
	fmt.Println("\nPush(h, 3),然后输出堆看看:")
	heap.Push(h, 3)
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	fmt.Println("\n 使用sort.Sort排序：")

	h2 := IntHeap{100, 16, 4, 8, 70, 2, 36, 22, 5, 12}
	sort.Sort(h2)
	for _, v := range h2 {
		fmt.Printf("%d ", v)
	}
}
