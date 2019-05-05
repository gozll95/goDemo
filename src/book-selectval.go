package main

import (
	"fmt"
)

var intChan1 chan int
var intChan2 chan int
var channels = []chan int{}

func init() {
	chancap := 5
	intChan1 = make(chan int, chancap)
	intChan2 = make(chan int, chancap)

	channels = append(channels, intChan1)
	channels = append(channels, intChan2)
}

func main() {

	select {
	case getChant(0) <- 1:
		fmt.Println("11")
	case getChant(1) <- 2:
		fmt.Println("22")
	default:
		fmt.Println("default")
	}
}

func getChant(i int) chan int {
	fmt.Println(i)
	fmt.Println(cap(channels[i]))
	return channels[i]
}

// select会执行所有case右边的语句
