package main

import (
	"fmt"
	"sync"
)

func main() {
	c := make(chan int, 10)
	wg := sync.WaitGroup{}
	wg.Add(10)
	go fibonnacci(&wg, 10, c)

	go func() {
		wg.Wait()
		close(c)
	}()

	for v := range c {
		fmt.Println(v)
	}
	fmt.Println("quit")

}

func fibonnacci(wg *sync.WaitGroup, n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		wg.Done()
		x, y = y, x+y
	}

}
