package main

func main() {
	defer func() {
		println("after panic")
	}()
	defer func() {
		panic("panic in defer")
	}()
	defer func() {
		if e := recover(); e != nil {
			println("recover from", e.(string))
			panic("error again")
		}
	}()
	panic("unexpected error")
}
