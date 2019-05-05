package main

import (
	"fmt"
	"runtime"
)

func Go(c chan bool, i int) {
	fmt.Println("i is", i)
	c <- true
}

func main() {
	//c := make(chan bool, 10)
	c := make(chan bool)
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < 10; i++ {
		go Go(c, i)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}
