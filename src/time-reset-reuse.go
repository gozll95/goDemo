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

	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			if !timer.Stop() { //timer过期才会drain channel
				fmt.Println("stop")
				<-timer.C
			}
			fmt.Println("reset")
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
