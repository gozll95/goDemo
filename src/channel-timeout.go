package main

import (
	"fmt"
	"time"
)

func main() {
	quit := make(chan bool)
	v := make(chan int)

	go func() {
		for {

			select {
			case <-time.After(4 * time.Second):
				t := time.Now()
				fmt.Println(t)
				quit <- true

				//break
			case <-v:
				time.Sleep(10 * time.Second)
				t := time.Now()
				fmt.Println(t)
			}

		}
	}()
	fmt.Println("xxxx")
	v <- 10
	<-quit
}
