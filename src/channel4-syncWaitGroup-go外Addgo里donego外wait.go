package main

import (
	"fmt"
	"runtime"
	"sync"
)

func Go(wg *sync.WaitGroup, index int) {
	fmt.Println("index is", index)
	wg.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go Go(&wg, i)
	}
	wg.Wait()
}
