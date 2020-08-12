package insert

import (
	"fmt"
	"testing"
)

func insertSort(nums []int) {
	for i := 1; i < len(nums); i++ {
		if nums[i] < nums[i-1] {
			j := i - 1
			temp := nums[i]
			for j >= 0 && nums[j] > temp {
				nums[j+1] = nums[j]
				j--
			}
			nums[j+1] = temp
		}
	}
}

func TestSort(t *testing.T) {
	nums := []int{1, 9, 5, 3, 1, 6, 2}
	insertSort(nums)
	fmt.Println(nums)
}
