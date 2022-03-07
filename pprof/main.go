package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/felixge/fgprof"
)

const cpuTime = 1000 * time.Millisecond

func main() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	go func() {
		http.Handle("/debug/fgprof", fgprof.Handler())
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	for {
		cpuIntensiveTask()
		slowNetworkRequest()
	}
}

func cpuIntensiveTask() {
	start := time.Now()

	for time.Since(start) <= cpuTime {
		for i := 0; i < 1000; i++ {
			_ = i
		}
	}
}

func slowNetworkRequest() {
	resp, err := http.Get("http://httpbin.org/delay/1")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
}
