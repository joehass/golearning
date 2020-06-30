package time

import (
	"fmt"
	"testing"
	"time"
)

func TestT1(t *testing.T) {
	d := time.Duration(time.Second * 2)

	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		<-timer.C

		fmt.Println("timeout...")
		// need reset
		timer.Reset(time.Second * 2)
	}
}

func TestT2(t *testing.T) {
	timer := time.NewTicker(3 * time.Second)
	defer timer.Stop()

	fmt.Println(time.Now())
	time.Sleep(4 * time.Second)

	for {

		select {

		case <-timer.C:

			fmt.Println(time.Now())

		}

	}
}
