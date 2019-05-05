package main

func Afunction(ch chan int) {
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1

	<-ch
}

func main() {
	//主routine的操作同上面那段代码
	ch := make(chan int, 10)
	for i := 0; i < 100; i++ {
		ch <- 1
		go Afunction(ch)
	}

	// 这段代码运行的结果为死锁
}

/*
上面这段运行和之前那一段基本上原理是一样的，但是运行后却会发生死锁。为什么呢？其实总结起来就一句话，"放得太快，取得太慢了"。



按理说，我们应该在我们主routine中创建子goroutine并每次向channel中放入数据，而子goroutine负责从channel中取出数据。但是我们的这段代码在创建了子goroutine后，每个routine会向channel中放入5个数据。这样，每向channel中放入6个数据才会执行一次取出操作，这样一来就可能会有某一时刻，channel已经满了，但是所有的routine都在执行放入操作(因为它们当前执行放入操作的概率是执行取出操作的6倍)，这样一来，所有的routine都阻塞了，从而导致死锁。



在使用带缓冲的channel时一定要注意放入与取出的速率问题。
*/
