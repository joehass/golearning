package context

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestC(t *testing.T) {
	//创建一个可取消子context,context.Background():返回一个空的Context，这个空的Context一般用于整个Context树的根节点。
	ctx, cancel := context.WithCancel(context.Background())
	ctxTwo, cancelTwo := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			//使用select调用<-ctx.Done()判断是否要结束
			case <-ctx.Done():
				fmt.Println("goroutineA exit")
				return
			default:
				fmt.Println("goroutineA running.")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)
	go func(ctx context.Context) {
		for {
			select {
			//使用select调用<-ctx.Done()判断是否要结束
			case <-ctx.Done():
				fmt.Println("goroutineB exit")
				return
			default:
				fmt.Println("goroutineB running.")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)
	go func(ctxTwo context.Context) {
		for {
			select {
			//使用select调用<-ctx.Done()判断是否要结束
			case <-ctxTwo.Done():
				fmt.Println("goroutineC exit")
				return
			default:
				fmt.Println("goroutineC running.")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctxTwo)

	time.Sleep(4 * time.Second)
	fmt.Println("main fun exit")
	//取消context
	cancel()
	cancelTwo()
	time.Sleep(5 * time.Second)

}

func TestC2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go watch(ctx, "【监控1】")
	go watch(ctx, "【监控2】")
	go watch(ctx, "【监控3】")

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "监控退出，停止了...")
			return
		default:
			time.Sleep(4 * time.Second)
			fmt.Println(name, "goroutine监控中...")
			return
		}
	}
}

func TestC3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	for i := 0; i < 10; i++ {
		go task(ctx)
	}

	time.Sleep(time.Second * 10)
}

func task(ctx context.Context) {
	ch := make(chan struct{}, 0)
	go func() {
		// 模拟4秒耗时任务
		time.Sleep(time.Second * 4)
		fmt.Println("1")
		fmt.Println("2")
		fmt.Println("3")
		fmt.Println("4")
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
		return
	}
}

var (
	mutex sync.Mutex
	id    int
)

func dosomething(ctx context.Context, val int) {
	mutex.Lock()
	defer mutex.Unlock()
	select {
	case <-ctx.Done():
		fmt.Println("op timeout", val)
		return
	default:
		time.Sleep(100 * time.Second)
		id = val
		fmt.Println("id", id)
	}

}

func TestC4(t *testing.T) {

	for i := 0; i < 3; i++ {
		done := make(chan bool)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Second)
			defer cancel()
			dosomething(ctx, i)
			done <- true
		}()
		select {
		case res := <-done:
			fmt.Println("done", res, id)
		case <-time.After(time.Duration(i) * time.Second):
			fmt.Println("timeout ", i)
		}
	}
}

func TestC5(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.TODO())
	ctx2, _ := context.WithCancel(ctx1)
	go func(ctx context.Context) {
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
	go func(ctxTwo context.Context) {
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
	time.Sleep(15 * time.Second)
	fmt.Println("exit")
}

func TestD1(t *testing.T) {
	ctx1, cancel1 := context.WithDeadline(context.TODO(), time.Now().Add(5*time.Second))
	ctx2, _ := context.WithDeadline(ctx1, time.Now().Add(5*time.Second))

	go func(ctx context.Context) {
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

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx3 exit")
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
