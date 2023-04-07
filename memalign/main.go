package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Foo struct {
	w byte   //1 byte
	x byte   //1 byte
	y uint64 //8 bytes
}
type Bar struct {
	x byte   //1 byte
	y uint64 //8 bytes
	w byte   // 1 byte
}

func main() {
	fmt.Println(runtime.GOARCH)
	newFoo := new(Foo)
	fmt.Println(unsafe.Sizeof(*newFoo))
	newBar := new(Bar)
	fmt.Println(unsafe.Sizeof(*newBar))
}
