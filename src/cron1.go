package main

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	now := time.Now()
	i := 0
	c := cron.New()
	spec := "@hourly"
	c.AddFunc(spec, func() {
		i++
		log.Println("start", i)
	})
	c.AddFunc("@hourly", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour") })
	c.Start()
	c.AddFunc("@hourly", func() { fmt.Println("Every day") })
	for _, e := range c.Entries() {
		log.Println(e.Schedule.Next(now))
	}
	select {} //阻塞主线程不退出

}
