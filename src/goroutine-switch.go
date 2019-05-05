package main

import (
	"fmt"
	"runtime"
	"time"
)

/* step 1
func main() { //这个main函数也是一个goroutine
	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				// fmt.Printf("Hello from "+
				// 	"goroutine %d\n", i) //这个会有I/O操作,会触发协程的切换
				a[i]++            //会始终在这个协程里,没有机会交出控制权
				runtime.Gosched() //手动交出控制权,让别人也有机会运行,这行代码很少写,因为有其他机会进行切换
			}
		}(i)
	}
	time.Sleep(time.Millisecond) //因为没人交出控制权,所以main goroutine一直执行不到这里
	fmt.Println(a)
}
*/

/*
卡死
*/

/* step 2
func main() { //这个main函数也是一个goroutine
	var a [10]int
	for i := 0; i < 10; i++ {
		go func() {
			for {
				a[i]++            //会始终在这个协程里,没有机会交出控制权
				runtime.Gosched() //手动交出控制权,让别人也有机会运行,这行代码很少写,因为有其他机会进行切换
			}
		}()
	}
	time.Sleep(time.Millisecond) //因为没人交出控制权,所以main goroutine一直执行不到这里
	fmt.Println(a)
}
*/

/*
这是闭包,到最后一次的时候i是10
panic: runtime error: index out of range

goroutine 10 [running]:
main.main.func1(0xc42001a140, 0xc420016098)
        /Users/flow/workspace/goDemo/src/test1/main.go:36 +0x45
created by main.main
        /Users/flow/workspace/goDemo/src/test1/main.go:34 +0x95
exit status 2
*/

func main() { //这个main函数也是一个goroutine
	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++            //会始终在这个协程里,没有机会交出控制权
				runtime.Gosched() //手动交出控制权,让别人也有机会运行,这行代码很少写,因为有其他机会进行切换
			}
		}(i)
	}
	time.Sleep(time.Millisecond) //因为没人交出控制权,所以main goroutine一直执行不到这里
	fmt.Println(a)
}

/*
一遍读 一遍并发写  这里其实需要channel
==================
WARNING: DATA RACE
Read at 0x00c420090000 by main goroutine:
  main.main()
      /Users/flow/workspace/goDemo/src/test1/main.go:70 +0xfe

Previous write at 0x00c420090000 by goroutine 6:
  main.main.func1()
      /Users/flow/workspace/goDemo/src/test1/main.go:64 +0x6b

Goroutine 6 (running) created at:
  main.main()
      /Users/flow/workspace/goDemo/src/test1/main.go:62 +0xc6
==================
[2935 2164 2068 426 429 437 455 384 397 431]
Found 1 data race(s)
exit status 66
*/

/*
goroutine可能的切换点
I/O,select
channel
等待锁
函数调用(有时)
runtime.Gosched()
*/
