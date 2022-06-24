package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {
	a := decimal.Decimal{}
	fmt.Println(a)               // 0
	a.Add(decimal.NewFromInt(1)) // ok
	p := decimal.NewFromInt(40)
	used := decimal.NewFromInt(0)
	for i := 0; i < 20; i++ {
		used = used.Add(decimal.NewFromInt(1).Div(p))
	}
	fmt.Println("precision", used.String()) //precision 0.5
	if used.LessThanOrEqual(decimal.NewFromFloat32(1 / float32(4) * float32(3-1))) {
		fmt.Println("congrats") // congrats
	} else {
		fmt.Println("oh no")
	}
	if res, err := decimal.NewFromString(""); err != nil {
		fmt.Println(err) // can't convert  to decimal
	} else {
		fmt.Println(res)
	}
	// 1 is exactly, 0.1 is not
	f, exact := decimal.NewFromFloat32(1).Float64()
	fmt.Println(f, exact) // 1 true
}
