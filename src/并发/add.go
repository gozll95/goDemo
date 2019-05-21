package main

import "fmt"

func Add(x, y int) {
	z := x + y
	fmt.Println(z)
}

func main() {
	for i := 0; i < 10; i++ {
		go Add(i, i)
	}
}

/*
在上面的代码里,我们在一个for  中调用了10次Add()函数,它们是并发执行的。可是
当你编译执行了上面的代码,就会发现一些  的现象: “什么?!  上什么都没有,程序没有正常工作!”
是什么原因呢?明明调用了10次Add(),应该有10次   出才对。要解释这个现象,就  及Go语言的程序执行机制了。
Go程序从 始化main package并执行main()函数开始,当main()函数返回时,程序退出, 且程序并不等 其他goroutine(非主goroutine)结 。
对于上面的例子,主函数 动了10个goroutine,然后返回,这时程序就退出了,而被 动的
执行Add(i, i)的goroutine没有来得及执行,所以程序没有任何 出。
OK,问题 到了,怎么解决呢?提到这一点, 计写过多线程程序的读者就已经恍然大悟，并且摩  掌地准 使用类似WaitForSingleObject之类的调用,或者写个自己很  的 等
或者  先进一些的sleep  等 来等 所有线程执行完 。
在Go语言中有自己推荐的方式,它要比这些方法都优 得多。
要让主函数等 所有goroutine退出后再返回,如何知道goroutine都退出了呢?这就引出了多个 goroutine之间通信的问题。下一节我们将主要解决这个问题。
*/
