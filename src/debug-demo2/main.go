//go-debug-profile-optimization/step0/demo.go
package main

import (
	"fmt"
	"time"
)

var (
	c       = make(chan struct{}, 1)
	vistors int
)

func handleVistors() {
	c <- struct{}{}
	vistors++
	<-c

}
func main() {
	for i := 0; i < 1000; i++ {
		go handleVistors()
	}
	time.Sleep(4 * time.Second)
	fmt.Println(vistors)
}
