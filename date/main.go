package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Unix(1625766611, 0)
	time.Sleep(time.Second)
	for ; start.Before(time.Now()); start = start.Add(time.Hour * 24) {
		fmt.Println(start.Format("sensego/plate/2006/01/02"))
	}
}
