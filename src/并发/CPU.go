package main

import (
	"fmt"
)

type Vector []float64

//分配给每个CPU的计算任务

func (v Vector) DoSome(i, n int, u Vector, c chan int) {
	for ; i < n; i++ {
		v[i] += i
	}
	c <- 1
}

const NCPU = 16

func (v Vector) DoAll(u Vector) {
	c := make(chan int, NCPU) //用于接收每个CPU的任务完成信号

	for i := 0; i < NCPU; i++ {
		go v.DoSome(i*len(v)/NCPU, (i+1)*len(v)/NCPU, u, c)
	}

	//等待所有CPU的任务完成

	for i := 0; i < NCPU; i++ {
		<-c //获取到一个数据，表示一个CPU计算完成了
	}

	//到这里表示所有计算已经结束
}

func main() {
	var v Vector
	v = []float64{1e9, 134}
	v.DoAll(v)
	fmt.Println("hehe")
}

/*
在Go语言升级到 认支持多CPU的某个 本之前,我们可以先通过设置  变量GOMAXPROCS的值来 制使用多少个CPU核心。
具体操作方法是通过直接设置  变量 GOMAXPROCS的值,或者在代码中 动goroutine之前先调用以下这个语句以设置使用16个CPU 核心:
	runtime.GOMAXPROCS(16)
到底应该设置多少个CPU核心呢,其实runtime包中还提供了另外一个函数NumCPU()来获取核心数。


*/
