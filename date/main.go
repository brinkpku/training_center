package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Unix(1613356602, 0)
	for ; start.Before(time.Now()); start = start.Add(time.Hour * 24) {
		fmt.Println(start.Format("sensego/plate/2006/01/02"))
	}
}
