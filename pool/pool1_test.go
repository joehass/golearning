package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestWPool(t *testing.T) {
	wp := NewWPool(5, 20)
	for i := 0; i < 20; i++ {
		k := i
		wp.Do(func() error {
			if k == 2 || k == 7 {
				time.Sleep(4 * time.Second)
			} else {
				time.Sleep(2 * time.Second)
			}
			fmt.Println(fmt.Sprintf("线程%v", k))
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}

func TestWPoolError(t *testing.T) {
	wp := NewWPool(5, 20)
	for i := 0; i < 20; i++ {
		k := i
		wp.Do(func() error {
			if k == 2 || k == 7 {
				return fmt.Errorf("线程%v执行失败", k)
			} else {
				time.Sleep(2 * time.Second)
			}
			fmt.Println(fmt.Sprintf("线程%v", k))
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}
