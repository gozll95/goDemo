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