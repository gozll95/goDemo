//http://studygolang.com/articles/2915
package main

import (
	"fmt"
	"time"
)

func main() {
	go t()

	c := make(chan struct{})
	timer := time.AfterFunc(10*time.Second, func() {
		fmt.Println("haha")
		close(c)
	})

	go func() {
		// for {
		timer.Reset(2 * time.Second) //2s之后过期
		fmt.Println("work")
		time.Sleep(3 * time.Second)
		// }
	}()

	for {
		select {
		case <-c:
			fmt.Println("quit")
			return
		}
	}

}

func t() {
	ticker := time.NewTicker(time.Second)

	for t := range ticker.C {
		fmt.Println(t)
	}

}

//这样做就相当于给Timer一个0秒的超时时间，让Timer立刻过期。

// func main() {
// 	timer := time.NewTimer(3 * time.Second)

// 	go func() {
// 		<-timer.C
// 		fmt.Println("Timer has expired.")
// 	}()

// 	timer.Stop()
// 	time.Sleep(60 * time.Second)
// }

/*
timer.NewTimer()会启动一个新的Timer实例，并开始计时。 我们启动一个新的goroutine，来以阻塞的方式从Timer的C这个channel中，等待接收一个值，这个值是到期的时间。并打印”Timer has expired.”

到现在看起来似乎没什么问题，但是当我们执行timer.Stop()之后，3秒钟过去了，程序却没有打印那句话。说明执行timer.Stop()之后，Timer自带的channel并没有关闭，而且这个Timer已经从runtime中删除了，所以这个Timer永远不会到期。

这会导致程序逻辑错误，或者更严重的导致goroutine和内存泄露。解决的办法是，使用timer.Reset()代替timer.Stop()来停止定时器。
*/

// package main

// import (
//     "time"
//     "fmt"
// )

// func main() {
//     timer := time.NewTimer(3 * time.Second)

//     go func() {
//         <-timer.C
//         fmt.Println("Timer has expired.")
//     }()

//     //timer.Stop()
//     timer.Reset(0  * time.Second)
//     time.Sleep(60 * time.Second)
// }
