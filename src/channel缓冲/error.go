package main

import "fmt"

func Afuntion(ch chan int) {
	fmt.Println("finish")
	<-ch
}

func main() {
	ch := make(chan int) //无缓冲的channel
	//只是把这两行的代码顺序对调一下
	ch <- 1
	go Afuntion(ch)

	// 输出结果：
	// 死锁，无结果
}

/*
第二段代码：

        1.  创建了一个无缓冲的channel

        2.  主routine要向channel中放入一个数据，但是因为channel没有缓冲，相当于channel一直都是满的，所以这里会发生阻塞。可是下面的那个goroutine还没有创建呢，主routine在这里一阻塞，整个程序就只能这么一直阻塞下去了，然后。。。然后就没有然后了。。死锁！

※从这里可以看出，对于无缓冲的channel，放入操作和取出操作不能再同一个routine中，而且应该是先确保有某个routine对它执行取出操作，然后才能在另一个routine中执行放入操作。
*/
