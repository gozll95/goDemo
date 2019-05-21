package main

import (
	"fmt"
	"time"
)

func Count(ch chan int, i int) {
	ch <- i
	fmt.Println("Counting", i)
}
func main() {
	chs := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}
	for i := 0; i < 10; i++ {
		<-chs[i]
		time.Sleep(1 * time.Millisecond)
	}

}
