package quick

import (
	"fmt"
	"testing"
)

//快速排序

func quickSort(values []int, left, right int) {
	temp := values[left]
	p := left
	i, j := left, right

	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}

		for i <= p && values[i] <= temp {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}
	}

	values[p] = temp

	if p-left > 1 {
		quickSort(values, left, p-1)
	}
	if right-p > 1 {
		quickSort(values, p+1, right)
	}
}

func TestSort(t *testing.T) {
	nums := []int{3, 9, 5, 1, 1, 6, 2}
	quickSort(nums, 0, len(nums)-1)
	fmt.Println(nums)
}
