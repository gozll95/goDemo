package main

import "fmt"

func main() {
	var (
		routineCtl chan int    = make(chan int, 20)
		feedback   chan string = make(chan string, 10000)

		msg      string
		allwork  int
		finished int
	)

	for i := 0; i < 1000; i++ {
		routineCtl <- 1
		allwork++
		go Afunction(routineCtl, feedback)
	}
}

func Afunction(routineControl chan int, feedback chan string) {
	defer func() {
		<-routineControl
		feedback <- "finish"
	}()
	fmt.Println("a")
}
