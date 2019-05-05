package main

import (
	"fmt"
	"sync"
)

func main() {
	var l *sync.Mutex
	l = new(sync.Mutex)
	l.Lock()
	defer l.Unlock()
	fmt.Println("1")
}
