package heap

import (
	"fmt"
	"testing"
)

//堆排序

func sort(arr []int) {
	//构建大顶堆，从最后一个非叶节点开始
	for i := len(arr)/2 - 1; i >= 0; i-- {
		//从第一个非叶子节点从上至下，从右至左调整结构
		adjustHeap(arr, i, len(arr))
	}
	//调整堆结构，交换堆顶元素和末尾元素
	for j := len(arr) - 1; j > 0; j-- {
		swap(arr, 0, j)       //将堆顶元素和末尾元素进行交换
		adjustHeap(arr, 0, j) //重新对堆进行调整
	}
}

//调整大顶堆，仅是调整过程，建立在大顶堆已构建的基础上
func adjustHeap(arr []int, i, length int) {
	temp := arr[i]
	for k := i*2 + 1; k < length; k = k*2 + 1 { //从i节点的左子节点开始，也就是2i+1处开始
		if k+1 < length && arr[k] < arr[k+1] { //如果左子节点小于右子节点，k指向右子节点
			k++
		}
		if arr[k] > temp { //如果子节点大于父节点，将子节点赋值给父节点，不用进行交换
			arr[i] = arr[k]
			i = k
		} else {
			break
		}
	}
	arr[i] = temp //将temp值放到最终位置
}

//交换元素
func swap(arr []int, a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}

func TestHeap(t *testing.T) {
	arr := []int{9, 6, 7, 5, 2, 4, 3, 1, 8}
	sort(arr)
	fmt.Println(arr)
}
