package decorator

import (
	"fmt"
	"testing"
)

type Beverage interface {
	getDescription() string
	cost() int
}

// 实现咖啡的过程
type Coffee struct {
	description string
}

func (this Coffee) getDescription() string {
	return this.description
}

func (this Coffee) cost() int {
	return 1
}

// Mocha 实现
type Mocha struct {
	beverage    Beverage
	description string
}

func (this Mocha) getDescription() string {
	return fmt.Sprintf("%s, %s", this.beverage.getDescription(), this.description)
}

func (this Mocha) cost() int {
	return this.beverage.cost() + 1
}

// Whip 实现
type Whip struct {
	beverage    Beverage
	description string
}

func (this Whip) getDescription() string {
	return fmt.Sprintf("%s, %s", this.beverage.getDescription(), this.description)
}

func (this Whip) cost() int {
	return this.beverage.cost() + 1
}

func TestD(t *testing.T) {
	var beverage Beverage
	// 买了一杯咖啡
	beverage = Coffee{description: "houseBlend"}
	// 给咖啡加上 Mocha
	beverage = Mocha{beverage: beverage, description: "Mocha"}
	// 给咖啡加上 Whip
	beverage = Whip{beverage: beverage, description: "whip"}
	// 最后计算 Coffee 的价格
	fmt.Println(beverage.getDescription(), ", cost is ", beverage.cost())
}
