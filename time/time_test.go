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
