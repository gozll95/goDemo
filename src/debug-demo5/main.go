//go-debug-profile-optimization/step0/demo.go
package main

import (
	"fmt"
	"time"
)

var add = make(chan struct{})
var total = make(chan int)

func getVistors() int {
	return <-total
}

func addVistors() {
	add <- struct{}{}
}

func teller() {
	var visitors int
	for {
		select {
		case <-add:
			visitors += 1
			fmt.Println(visitors)
		case total <- visitors:
		}
	}
}

func handleVistor() {
	addVistors()

}
func main() {
	go teller()
	for i := 0; i < 10; i++ {
		go handleVistor()
	}
	time.Sleep(5 * time.Second)
	vistor := getVistors()
	fmt.Println(vistor)

}
