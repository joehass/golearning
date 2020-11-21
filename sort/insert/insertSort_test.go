package insert

import (
	"fmt"
	"testing"
)

//直接插入排序
func insertSort(arr []int) {
	len := len(arr)
	for i := 0; i < len; i++ {
		selected := arr[i]
		for j := i - 1; j >= 0; j-- {
			if arr[j] > selected {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			} else {
				arr[j+1] = selected
				break
			}
		}
	}
}

func TestSort(t *testing.T) {
	nums := []int{1, 9, 5, 3, 1, 6, 2}
	insertSort(nums)
	fmt.Println(nums)
}
