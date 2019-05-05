package main

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}

//time.Tick函数表现得好像它创建了一个在循环中调用time.Sleep的goroutine，
// 每次被唤醒时发送一个事件。当countdown函数返回时，它会停止从tick中接收事件，
// 但是ticker这个goroutine还依然存活，继续徒劳地尝试从channel中发送值，
// 然而这时候已经没有其它的goroutine会从该channel中接收值了
// --这被称为goroutine泄露(§8.4.4)。
