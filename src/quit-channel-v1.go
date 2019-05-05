package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	quit := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- true
	}()
	fibonnacci(c, quit)
}

func fibonnacci(c chan int, quit chan bool) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			fmt.Println("i am xxxxxxxxxx")
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return

			// case c <- x:
			// 	fmt.Println("i am xxxxxxxxxx")
			// 	x, y = y, x+y
		}
	}

}
