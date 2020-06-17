package composite

import (
	"fmt"
	"strings"
	"testing"
)

//组合模式
type ICompany interface {
	add(ic ICompany)
	remove(ic ICompany)
	display(depth int)
	lineOfDuty()
}

type concreteCompany struct {
	name string
	list map[ICompany]ICompany
}

func NewConcreteCompany(name string) *concreteCompany {
	list := make(map[ICompany]ICompany)

	return &concreteCompany{
		name: name,
		list: list,
	}
}

func (cc *concreteCompany) add(ic ICompany) {
	cc.list[ic] = ic
}

func (cc *concreteCompany) remove(ic ICompany) {
	delete(cc.list, ic)
}

func (cc *concreteCompany) display(depth int) {
	fmt.Println(strings.Repeat("-", depth), " ", cc.name)
	for _, ccc := range cc.list {
		ccc.display(depth + 2)
	}
}

func (cc *concreteCompany) lineOfDuty() {
	for _, ccc := range cc.list {
		ccc.lineOfDuty()
	}
}

type HRDepartment struct {
	name string
}

func newHRDepartment(name string) *HRDepartment {
	return &HRDepartment{
		name: name,
	}
}

func (hrd *HRDepartment) add(ic ICompany) {
}

func (hrd *HRDepartment) remove(ic ICompany) {
}

func (hrd *HRDepartment) display(depth int) {
	fmt.Println(strings.Repeat("-", depth), " ", hrd.name)
}

func (hrd *HRDepartment) lineOfDuty() {
	fmt.Println(hrd.name, "员工招聘培训管理")
}

type FinanceDepartment struct {
	name string
}

func newFinanceDepartment(name string) *FinanceDepartment {
	return &FinanceDepartment{
		name: name,
	}
}

func (fd *FinanceDepartment) add(ic ICompany) {
}

func (fd *FinanceDepartment) remove(ic ICompany) {
}

func (fd *FinanceDepartment) display(depth int) {
	fmt.Println(strings.Repeat("-", depth), " ", fd.name)
}

func (fd *FinanceDepartment) lineOfDuty() {
	fmt.Println(fd.name, "公司财务收支管理")
}

func TestC(t *testing.T) {
	root := NewConcreteCompany("北京总公司")
	root.add(newHRDepartment("总公司人力资源部"))
	root.add(newFinanceDepartment("总公司财务"))

	com := NewConcreteCompany("上海华东分公司")
	com.add(newHRDepartment("上海华东分公司人力资源部"))
	com.add(newFinanceDepartment("上海华东分公司财务"))
	root.add(com)

	com1 := NewConcreteCompany("南京办事处")
	com1.add(newHRDepartment("南京办事处人力资源部"))
	com1.add(newFinanceDepartment("南京办事处财务"))
	com.add(com1)

	com2 := NewConcreteCompany("杭州办事处")
	com2.add(newHRDepartment("杭州办事处人力资源部"))
	com2.add(newFinanceDepartment("杭州办事处财务"))
	com.add(com2)

	fmt.Println("结构图:")
	root.display(1)
	fmt.Println("职责:")
	root.lineOfDuty()
}
