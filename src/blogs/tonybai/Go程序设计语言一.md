# 参数当即求值,defer稍后执行

package main

import "fmt"

func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}
func un(s string) {
	fmt.Println("leaving:", s)
}
func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}
func b() {
	defer un(trace("b"))
	fmt.Println("in b")
	a()
}
func main() { b() }

/*
entering: b
in b
entering: a
in a
leaving: a
leaving: b
*/

# 函数字面值是闭包(closure)
函数字面值实际上是闭包。

func adder()(func(int)int){
    var x int
    return func(delta int)int{
        x+=delta
        return x
    }
}

f:=adder()
fmt.Print(f(1))
fmt.Print(f(20))
fmt.Print(f(300))

输出1 21 321 – f中的x累加。