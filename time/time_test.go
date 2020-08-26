package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTime1(t *testing.T) {
	fmt.Println(time.Now().AddDate(0, 0, 0))
	fmt.Println(time.Now().Round(0))

	fmt.Println(time.Now().AddDate(0, 0, 0) == time.Now().Round(0))

	fmt.Println(time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
}

func TestT12(t *testing.T) {
	fmt.Println(time.Now().AddDate(0, 1, 0).Unix())
}

func TestT13(t *testing.T) {
	tm := time.Unix(1598197930, 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
}
