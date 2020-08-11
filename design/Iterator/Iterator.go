package Iterator

/**
迭代器模式，为一个容器设置一个迭代函数，可以使用这个迭代函数来顺序访问其中的每一个元素，而外部无需知道底层实现
*/

type IAggregate interface {
	Iterator() IIterator
}

type IIterator interface {
	HasNext() bool
	Current() int
	Next() bool
}

type Aggregate struct {
	container []int //容器中装载int型数据
}

//迭代器
type Iterator struct {
	cursor    int        //当前游标
	aggregate *Aggregate //对应的容器指针
}

//判断是否迭代到最后，如果没有，则返回true
func (i *Iterator) HasNext() bool {
	if i.cursor+1 < len(i.aggregate.container) {
		return true
	}

	return false
}

//获取当前迭代元素（从容器中去除当前游标对应的元素）
func (i *Iterator) Current() int {
	return i.aggregate.container[i.cursor]
}

//将游标指向下一个元素
func (i *Iterator) Next() bool {
	if i.cursor < len(i.aggregate.container) {
		i.cursor++
		return true
	}
	return false
}

//创建一个迭代器，并让迭代器中的指针指向当前对象
func (a *Aggregate) Iterator() IIterator {
	i := new(Iterator)
	i.aggregate = a
	return i
}
