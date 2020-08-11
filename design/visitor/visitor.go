package visitor

import "fmt"

/**
访问者模式：就是在不更改这个结构体的前提下，更改这个结构体中方法所能执行的逻辑
*/
type IVisitor interface {
	Visit()
}

type ProductionVisitor struct {
}

func (v ProductionVisitor) Visit() {
	fmt.Println("这是生产环境")
}

type TestingVisitor struct {
}

func (t TestingVisitor) Visit() {
	fmt.Println("这是测试环境")
}

type Element struct {
}

func (el Element) Accept(visitor IVisitor) {
	visitor.Visit()
}

type EnvExample struct {
	Element
}

func (e EnvExample) Print(visitor IVisitor) {
	e.Element.Accept(visitor)
}
