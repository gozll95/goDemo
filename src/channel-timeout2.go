package main

import (
	"fmt"
	"time"
)

func main() {
	timeout := make(chan bool, 1)
	//timeout := make(chan bool)
	go func() {
		time.Sleep(1e9) // sleep one second
		timeout <- true
	}()
	ch := make(chan int)
	select {
	case <-ch:
	case <-timeout:
		fmt.Println("timeout!")
	}
}
