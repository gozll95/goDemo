package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("the 1")
	tc := time.Tick(time.Second * 2)
	for i := 1; i <= 2; i++ {
		<-tc
		fmt.Println("hello", time.Now())
	}
}
