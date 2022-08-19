package main

func main() {
	defer func() {
		println("after panic") // will run even the second defer panic
	}()
	defer func() {
		panic("panic in defer") // will run even the last defer panic
	}()
	defer func() {
		if e := recover(); e != nil {
			println("recover from", e.(string))
			panic("error again") // run with LIFO order
		}
	}()
	panic("unexpected error")
}
