package main

import "fmt"

func main() {
	s := "z"
	defer func() {
		fmt.Println(s)
	}()
	s = "ttt"
}
