package Iterator

import (
	"fmt"
	"testing"
)

func TestI(t *testing.T) {
	//创建容器，并放入初始化数据
	c := &Aggregate{container: []int{1, 2, 3, 4}}

	//获取迭代器
	iterator := c.Iterator()

	for {
		//打印当前数据
		fmt.Println(iterator.Current())
		if iterator.HasNext() {
			iterator.Next()
		} else {
			break
		}
	}
}
