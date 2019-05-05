//blog:http://tonybai.com/articles/

package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 7)
			c <- false
		}

		time.Sleep(time.Second * 7)
		c <- true
	}()

	//consumer
	go func() {
		// try to read from channel, block at most 5s.
		// if timeout, print time event and go on loop.
		// if read a message which is not the type we want(we want true, not false),
		// retry to read.
		timer := time.NewTimer(time.Second * 5)
		for {
			// timer may be not active, and fired
			if !timer.Stop() {
				select {
				case <-timer.C: //try to drain from the channel
				default:
				}
			}
			timer.Reset(time.Second * 5)
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

如果我们不想重复创建这么多Timer实例，而是reuse现有的Timer实例，
那么我们就要用到Timer的Reset方法，见下面example2.go，考虑篇幅，
这里仅列出consumer routine代码，其他保持不变：


按照go 1.7 doc中关于Reset使用的建议：

To reuse an active timer, always call its Stop method first and—if it had expired—drain the value from its channel. For example:

if !t.Stop() {
        <-t.C
}
t.Reset(d)
*/

/*
producer的发送行为发生了变化，Comsumer routine在收到第一个数据前有了一次time expire的事件，
for loop回到loop的开始端。这时timer.Stop函数返回的不再是true，而是false，因为timer已经expire，
最小堆中已经不包含该timer了，Stop在最小堆中找不到该timer，返回false。于是example3代码尝试抽干(drain)timer.C中的数据。
但timer.C中此时并没有数据，于是routine block在channel recv上了。
*/

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run time-reset-reuse-最终版.go
2017-08-04 13:40:54.961439814 +0800 CST :timer expired
2017-08-04 13:40:56.961541396 +0800 CST :recv false. continue
2017-08-04 13:41:01.961898788 +0800 CST :timer expired
2017-08-04 13:41:03.962523929 +0800 CST :recv false. continue
2017-08-04 13:41:08.963837941 +0800 CST :timer expired
2017-08-04 13:41:10.963789902 +0800 CST :recv false. continue
2017-08-04 13:41:15.965099688 +0800 CST :timer expired
2017-08-04 13:41:17.965060352 +0800 CST :recv false. continue
2017-08-04 13:41:22.966411879 +0800 CST :timer expired
*/
