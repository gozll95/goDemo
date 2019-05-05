package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 1)
			c <- false
		}

		// time.Sleep(time.Second * 1)
		// c <- true
	}()

	go func() {
		for {
			timer := time.NewTimer(time.Second * 6)
			defer timer.Stop()
			select {
			case b := <-c:
				if b == false {
					fmt.Println(time.Now(), ":recv false. continue")
					continue
				}
				//we want true, not false
				fmt.Println(time.Now(), ":recv true. return")
				return
			case <-timer.C:
				fmt.Println(time.Now(), ":timer expired")
				continue
			}
		}
	}()

	//to avoid that all goroutine blocks.
	var s string
	fmt.Scanln(&s)
}

/*
综上，作为Timer的使用者，我们要做的就是尽量减少在使用Timer时对最小堆管理goroutine和GC的压力即可，
即：及时调用timer的Stop方法从最小堆删除timer element(如果timer 没有expire)以及reuse active timer。
*/
