package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wait := sync.WaitGroup{}
	locker := new(sync.Mutex)
	cond := sync.NewCond(locker)

	for i := 0; i < 3; i++ {
		go func(i int) {
			defer wait.Done()
			wait.Add(1)
			cond.L.Lock()
			fmt.Println("Waiting start...")
			cond.Wait()
			fmt.Println("Waiting end...")
			cond.L.Unlock()

			fmt.Println("Goroutine run. Number:", i)
		}(i)
	}

	time.Sleep(2e9)
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()

	time.Sleep(2e9)
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()

	time.Sleep(2e9)
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()

	wait.Wait()
}

/*
result:
Waiting start...
Waiting start...
Waiting start...
Waiting end...
Goroutine run. Number: 1
Waiting end...
Goroutine run. Number: 0
Waiting end...
Goroutine run. Number: 2


可以看出来，每执行一次Signal就会执行一个goroutine。如果想让所有的goroutine执行，那么将所有的Signal换成一个Broadcast方法可以。
*/
