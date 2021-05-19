package tunny

import (
	"context"
	"sync"
	"time"
)

type params struct {
	wait int64
	name string
}

func main() {
	numCPUs := 2
	wg := sync.WaitGroup{}
	wg.Add(3)

	pool := NewFunc(numCPUs, func(payload interface{}) interface{} {
		param := payload.(*params)
		println(param.name, "start")
		// TODO: Something CPU heavy with payload
		time.Sleep(time.Second * time.Duration(param.wait))
		println(param.name, "finished")
		wg.Done()
		return nil
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
	cancel()
	time.Sleep(time.Second * 2)
	go pool.Process(param3)
	wg.Wait()
}
