package main

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "", log.LstdFlags))
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc), cron.WithChain(cron.SkipIfStillRunning(logger)))
	if _, err := c.AddFunc("*/1 * * * *", job); err != nil {
		log.Fatalf("add cron job error:%v", err)
	}

	ch := make(chan int)
	c.Start()
	<-ch
}

func job() {
	log.Println("job started")
	time.Sleep(time.Second * 65)
	log.Println("job finished")
}
