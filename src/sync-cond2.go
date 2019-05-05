//cond的主要作用就是获取锁之后,wait()方法会等待一个通知，来进行下一步锁释放等操作，以此控制锁何时释放，释放频率。

package main

import (
	"fmt"
	"sync"
	"time"
)

var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)

func test(x int) {
	cond.L.Lock() //获取锁
	cond.Wait()   //等待通知  暂时阻塞
	fmt.Println(x)
	// time.Sleep(time.Second * 1)
	cond.L.Unlock() //释放锁
}
func main() {
	for i := 0; i < 40; i++ {
		go test(i)
	}
	fmt.Println("start all")
	time.Sleep(time.Second * 3)
	fmt.Println("broadcast")
	cond.Signal() // 下发一个通知给已经获取锁的goroutine
	time.Sleep(time.Second * 3)
	cond.Signal() // 3秒之后 下发一个通知给已经获取锁的goroutine
	time.Sleep(time.Second * 3)
	cond.Broadcast() //3秒之后 下发广播给所有等待的goroutine
	time.Sleep(time.Second * 60)
}

/*
package main

import (
	"fmt"
	"sync"
	"time"
)

var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)
var wait = sync.WaitGroup{}

func test(x int) {
	defer wait.Done()
	cond.L.Lock()
	cond.Wait()
	fmt.Println(x)
	cond.L.Unlock()
}

func main() {
	for i := 0; i < 40; i++ {
		wait.Add(1)
		go test(i)
	}
	fmt.Println("start all")
	time.Sleep(1 * time.Second)
	fmt.Println("signal--1")
	cond.Signal()
	time.Sleep(1 * time.Second)
	fmt.Println("signal--2")
	cond.Signal()
	time.Sleep(1 * time.Second)
	fmt.Println("broadcast")
	cond.Broadcast()
	wait.Wait()
}

*/