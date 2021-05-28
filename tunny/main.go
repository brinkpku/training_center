package main

import (
	"context"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
)

type params struct {
	wait int64
	name string
}

func main() {
	numCPUs := 2
	wg := sync.WaitGroup{}
	wg.Add(3)

	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		param := payload.(*params)
		println(param.name, "start")
		timer := time.NewTimer(time.Second * time.Duration(param.wait))
		ticker := time.NewTicker(time.Second)
		runned := 0
		for {
			select {
			case <-timer.C:
				println(param.name, "finished")
				wg.Done()
				return nil
			case <-ticker.C:
				runned++
				println(param.name, "is running", runned)
			}
		}
	})
	defer pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	param1 := &params{
		wait: 10,
		name: "worker1",
	}
	param2 := &params{
		wait: 10,
		name: "worker2",
	}
	param3 := &params{
		wait: 10,
		name: "worker3",
	}
	go pool.ProcessCtx(ctx, param1)
	go pool.ProcessCtx(ctx, param2)
	time.Sleep(time.Second * 3)
	cancel() // can't cancel, because it has run
	go pool.Process(param3)
	wg.Wait()
}
