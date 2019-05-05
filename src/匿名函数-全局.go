package main

import (
	"fmt"
)

func square() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

func main() {
	f := square()
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
}

// 1
// 4
// 9
// 16

/*
对于匿名函数来说，是全局变量
squares返回一个匿名函数
该匿名函数每次被调用时都会返回下一个数的平方
*/
