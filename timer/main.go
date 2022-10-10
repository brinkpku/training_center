package main

import (
	"fmt"
	"time"
)

func main() {
	mstick()
}

func timer() {
	timer := time.NewTimer(time.Second * 2)
	fmt.Println("main", time.Now())
	runtimes := 0
	for {
		select {
		case <-timer.C:
			fmt.Println("start", time.Now())
			time.Sleep(time.Second * 3)
			runtimes += 1
			fmt.Println("finish", runtimes, time.Now())
			timer.Reset(time.Second * 2)
		}
	}
}

func tick() { // when run time> tick timer, need use timer
	ticker := time.NewTicker(time.Second * 2)
	fmt.Println("main", time.Now())
	runtimes := 0
	for {
		select {
		case <-ticker.C:
			fmt.Println("start", time.Now())
			time.Sleep(time.Second * 3)
			runtimes += 1
			fmt.Println("finish", runtimes, time.Now())
		}
	}
}

func mstick() {
	ticker := time.NewTicker(time.Millisecond * 20)
	for i := 0; i < 10; i++ {
		<-ticker.C
		fmt.Println(time.Now().UnixNano() / 1e6)
	}
}
