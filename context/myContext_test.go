package context

import (
	"fmt"
	"testing"
	"time"
)

func TestMyC1(t *testing.T) {
	ctx1, cancel1 := WithCancel(TODO())
	ctx2, _ := WithCancel(ctx1)
	go func(ctx Context) {
		for {
			select {
			//使用select调用<-ctx.Done()判断是否要结束
			case <-ctx.Done():
				fmt.Println("ctx1 exit")
				return
			default:
				fmt.Println("goCtx1 running.")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx1)
	go func(ctxTwo Context) {
		for {
			select {
			//使用select调用<-ctx.Done()判断是否要结束
			case <-ctxTwo.Done():
				fmt.Println("ctx2 exit")
				return
			default:
				fmt.Println("goCtx2 running.")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx2)
	time.Sleep(4 * time.Second)
	fmt.Println("main fun exit")
	//取消context
	cancel1()
	//cancelTwo()
	time.Sleep(5 * time.Second)
	fmt.Println("exit")
}

func TestMyD1(t *testing.T) {
	ctx1, cancel1 := WithDeadline(TODO(), time.Now().Add(5*time.Second))
	ctx2, _ := WithDeadline(ctx1, time.Now().Add(5*time.Second))

	go func(ctx Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx1 exit")
				return
			default:
				fmt.Println("ctx1 running ")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx1)

	go func(ctx Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx2 exit")
				return
			default:
				fmt.Println("ctx2 running ")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx2)

	time.Sleep(4 * time.Second)
	fmt.Println("main fun exit")
	//取消context
	cancel1()
	//cancelTwo()
	time.Sleep(5 * time.Second)
	fmt.Println("exit")
}
