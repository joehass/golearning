package time

import "time"

type Timer struct {
	C <-chan time.Time
	r runtimeTimer
}

type runtimeTimer struct {
	tb uintptr //存储当前定时器的数组地址
	i  int     //存储当前定时器的数组下标

	when   int64                      //当前定时器触发时间
	period int64                      //当前定时器周期触发间隔
	f      func(interface{}, uintptr) //定时器触发时执行的函数
	arg    interface{}                //定时器触发时执行函数传递的参数一
	seq    uintptr                    //定时器触发时执行函数传递的参数二（该参数只在网络收发场景下使用）
}

//func NewTimer(d time.Duration) *Timer {
//	c := make(time.Time,1) //创建一个管道
//	t := &Timer{
//		C: c,
//		r: runtimeTimer{
//			when:
//		},
//	}
//}

//func when(d time.Duration) int64 {
//
//}
