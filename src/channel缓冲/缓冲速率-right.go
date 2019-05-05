package main

import (
	"fmt"
)

func Afunction(ch chan int) {
	fmt.Println("finish")
	<-ch //goroutine执行完了就从channel取出一个数据
}

func main() {
	ch := make(chan int, 10)
	for i := 0; i < 1000; i++ {
		//每当创建goroutine的时候就向channel中放入一个数据，如果里面已经有10个数据了，就会
		//阻塞，由此我们将同时运行的goroutine的总数控制在<=10个的范围内
		ch <- 1
		go Afunction(ch)
	}
	// 这里只是示范个例子，当然，接下来应该有些更加周密的同步操作
}

/*
※从这里可以看出，对于无缓冲的channel，放入操作和取出操作不能再同一个routine中，而且应该是先确保有某个routine对它执行取出操作，然后才能在另一个routine中执行放入操作。



对于带缓冲的channel，就没那么多讲究了，因为有缓冲空间，所以只要缓冲区不满，放入操作就不会阻塞，同样，只要缓冲区不空，取出操作就不会阻塞。而且，带有缓冲的channel的放入和取出可以用在同一个routine中。

但是，并不是说有了缓冲就可以随意使用channel的放入和取出了，我们一定要注意放入和取出的速率问题。下面我们就举个例子来说明这种问题：
*/
