package runtime

import (
	"fmt"
	"runtime"
	"testing"
)

//得到函数调用名称
func TestR1(t *testing.T) {
	foo()
}

func foo() {
	fmt.Printf("我是 %s, %s 在调用我!\n", printMyName(), printCallerName())
	bar()
}

func bar() {
	fmt.Printf("我是 %s, %s 又在调用我!\n", printMyName(), printCallerName())
}

func printMyName() string {
	//caller可以返回函数调用栈的某一层的程序计数器、文件信息、行号
	//0表示当前函数，也是调用runtime.Caller的函数。1代表上一层调用者，以此类推
	pc, _, _, _ := runtime.Caller(1)

	//FuncForPC 可以吧程序计数器地址对应的函数的信息获取出来。如果因为内联程序计数器对应多个函数，它返回最外面的函数
	return runtime.FuncForPC(pc).Name()
}

func printCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}
