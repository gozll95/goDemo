package main

import (
	"fmt"
)

func main() {
	c := make(chan bool)
	go func() {
		fmt.Println("it is goroutine")
		c <- true
		close(c)
	}()
	for v := range c {
		fmt.Println(v)
	}
	fmt.Println("over")
}
