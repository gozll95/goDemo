package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	o := &sync.Once{}
	go do(o)
	go do(o)
	time.Sleep(time.Second * 2)
}

func do(o *sync.Once) {
	fmt.Println("Start do")
	o.Do(func() {
		fmt.Println("Doing something...")
	})
	fmt.Println("Do end")
}

/*
输出结果：

Start do
Doing something...
Do end
Start do
Do end
这里 Doing something 只被调用了一次。

go的源码实现:

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
就是用了原子计数记录被执行的次数。使用Mutex Lock Unlock锁定被执行函数，防止被重复执行。


*/
