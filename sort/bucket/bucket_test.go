package bucket

/**
桶排序
*/

func tong(arr []int) []int {
	t := make([]int, 10)
	for _, val := range arr {
		t[val]++
	}
	res := make([]int, 0, len(arr))
	for index, val := range t {
		//循环把排序元素添加到新的数组中
		for ; val > 0; val-- {
			res = append(res, index)
		}
	}
}
