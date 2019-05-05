package main

import (
	"fmt"
)

func main() {
	msg := make(chan bool)
	go func() {
		fmt.Println("this is a gorouting")
		<-msg
	}()
	msg <- true
}
