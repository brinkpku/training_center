package main

import (
	"fmt"
	"math"
)

// see decimal
func main() {
	p := float32(40)
	used := float32(0)
	for i := 0; i < 20; i++ {
		used += 1 / p
	}
	fmt.Println(used, 1/float32(4)*float32(3-1))
	if GetValidFloat32(used) <= 1/float32(4)*float32(3-1) {
		fmt.Println("congrats")
	} else {
		fmt.Println("oh no")
	}
}

func GetValidFloat32(val float32) float32 {
	// [IEEE 754]
	return float32(RoundFloat64(float64(val), 7))
}

func RoundFloat64(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
