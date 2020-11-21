package _interface

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInter(t *testing.T) {
	var a interface{}
	fmt.Println(a)
}

func TestI2(t *testing.T) {
	a := new(interface{})
	fmt.Println(a)
}

func TestT3(t *testing.T) {
	a := new(interface{})

	v := reflect.ValueOf(a)
	fmt.Println(v.IsNil())
}
