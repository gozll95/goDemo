package main

import (
	"fmt"
)

func fibonnacci(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	//close(c)
}

func main() {
	c := make(chan int, 10)
	go fibonnacci(cap(c), c)
	/*
		for v := range c {
			fmt.Println(v)
		}
	*/
	for i := 0; i < 10; i++ {
		v := <-c
		fmt.Println(v)
	}
}
