package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	for ; start.Before(time.Now()); start = start.Add(time.Hour * 24) {
		fmt.Println(start.Format("sensego/plate/2006/01/02"))
	}
	today := getTodayZero()
	fmt.Println(today.Add(-time.Hour * 24 * 90).Format("2006/1/2"))
}

func getTodayZero() time.Time {
	now := time.Now()
	zero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return zero
}
