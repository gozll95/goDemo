package main

import (
	"fmt"
	"time"
)

func main() {

		sendTicker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("timeout")
				return
			case <-sendTicker.C:
				fmt.Println("hello world")
			}
		}


}
