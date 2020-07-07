package math

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestMath(t *testing.T) {
	//var j float64
	wordCount := 2789
	bill := 5
	var price int

	price = wordCount / 1000 * bill
	if wordCount <= 100 {
		price += getPrice(0)
	} else if wordCount%1000 > 500 || (wordCount > 100 && wordCount < 1000) {
		price += getPrice(1000)
	} else if wordCount%1000 <= 500 && wordCount%1000 != 0 {
		price += getPrice(500)
	}
	fmt.Println(price)
}

func getPrice(wordCount int) int {
	r, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(wordCount)/float64(1000)*float64(5)), 64)
	price := int(math.Ceil(r))
	return price
}

func TestT2(t *testing.T) {
	fmt.Println(math.Ceil(3.1))
}
