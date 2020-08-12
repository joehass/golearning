package bubble

import (
	"fmt"
	"testing"
)

func bubbleSort(nums *[]int) {
	//控制比较次数
	for i := 0; i < len(*nums); i++ {
		////j<len(*nums)-i:是为了不去比较已经到最后的最大数据
		//j =1:是为了每次都要从开始比较
		for j := 1; j < len(*nums)-i; j++ {
			if (*nums)[j] < (*nums)[j-1] {
				(*nums)[j], (*nums)[j-1] = (*nums)[j-1], (*nums)[j]
			}
		}
	}
}

func TestSort(t *testing.T) {
	nums := []int{1, 9, 5, 3, 1, 6, 2}
	bubbleSort(&nums)
	fmt.Println(nums)
}
