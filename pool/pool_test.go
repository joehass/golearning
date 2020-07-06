package pool

import (
	"fmt"
	"testing"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func TestWorkerPoolStart(t *testing.T) {
	wp := workpool.New(10)    // Set the maximum number of threads
	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { // 0-10 values per print
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			//time.Sleep(1 * time.Second)
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}

func TestWorkerPoolError(t *testing.T) {
	wp := workpool.New(10)    // Set the maximum number of threads
	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { // 0-10 values per print

				if j == 9 {
					return fmt.Errorf("线程%v,执行%v失败", ii, j)
				}
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			//time.Sleep(1 * time.Second)
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}
