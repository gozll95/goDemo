//虽然chanel也能实现，但是觉得如果涉及不到子线程与主线程数据同步，这个感觉不错。
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	num = 10
)

func main() {
	TestFunc("testchan", TestChan)
}

func TestFunc(name string, f func()) {
	st := time.Now().UnixNano()
	f()
	fmt.Printf("task %s cost %d \r\n", name, (time.Now().UnixNano()-st)/int64(time.Millisecond))
}

func TestChan() {
	var wg sync.WaitGroup
	c := make(chan string)
	wg.Add(1)

	go func() {
		for _ = range c {
		}
		wg.Done()
	}()

	for i := 0; i < num; i++ {
		c <- "123"
	}

	close(c)
	wg.Wait()
}
