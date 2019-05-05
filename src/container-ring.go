package main

import (
	"container/ring" //闭环包引入，详见/usr/local/go/src/pkg/container/ring
	"fmt"
)

func main() {
	//创建10个元素的闭环
	r := ring.New(10)

	//给闭环中的元素赋值
	for i := 1; i <= r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}

	printRing(r)

	//获得当前元素之后的第5个元素
	r5 := r.Move(5)

	printRing(r5)
	//fmt.Println(r5)
	//fmt.Println(r)
	printRing(r)

	//链接当前元素r与r5，相当于删除了r与r5之间的元素
	r1 := r.Link(r5)
	printRing(r1)
	//fmt.Println(r1)
	//fmt.Println(r)
	printRing(r)
}


func printRing(r *ring.Ring){
	fmt.Println("start to print")
		//循环打印闭环中的元素值
	r.Do(
		func(p interface{}) {
			println(p.(int))
		})
	fmt.Println("end print")
}