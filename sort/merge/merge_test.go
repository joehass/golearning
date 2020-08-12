package merge

import (
	"fmt"
	"testing"
)

//归并排序

func MergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	m := len(arr) / 2
	l := MergeSort(arr[:m])
	r := MergeSort(arr[m:])
	return merge(l, r)
}

func merge(a, b []int) (c []int) {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			c = append(c, a[i])
			i++
		} else {
			c = append(c, b[j])
			j++
		}
	}

	c = append(c, a[i:]...)
	c = append(c, b[j:]...)
	return c
}

func TestSort(t *testing.T) {
	nums := []int{1, 9, 5, 3, 1, 6, 2}
	num := MergeSort(nums)
	fmt.Println(num)
}
