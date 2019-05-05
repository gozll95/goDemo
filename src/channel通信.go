在Go编程中，Channel是被推荐的执行体间通信的方法，Go的编译器和运行态都会尽力对其优化。

对一个Channel的发送操作(send) happens-before 相应Channel的接收操作完成
关闭一个Channel happens-before 从该Channel接收到最后的返回值0
不带缓冲的Channel的接收操作（receive） happens-before 相应Channel的发送操作完成


var c = make(chan int, 10)
var a string
func f() {
    a = "hello, world"
    c <- 0
}
func main() {
    go f()
    <-c
    print(a)
}

上述代码可以确保输出hello, world，因为a = "hello, world" happens-before c <- 0，print(a) happens-after <-c， 

根据上面的规则1）以及happens-before的可传递性，a = "hello, world" happens-beforeprint(a)。

根据规则2）把c<-0替换成close(c)也能保证输出hello,world，因为关闭操作在<-c接收到0之前发送。

var c = make(chan int)
var a string
func f() {
    a = "hello, world"
    <-c
}
func main() {
    go f()
    c <- 0
    print(a)
}
根据规则3），因为c是不带缓冲的Channel，a = "hello, world" happens-before <-c happens-before c <- 0 happens-before print(a)， 但如果c是缓冲队列，如定义c = make(chan int, 1), 那结果就不确定了。


/*
有缓冲，先赋才能取
关闭一个channel才能不断取零值
无缓冲,先取才能赋->持保留意见
*/