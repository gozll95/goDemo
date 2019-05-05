close

内建的close方法可以用来关闭channel。

总结一下channel关闭后sender的receiver操作。
如果channel c已经被关闭,继续往它发送数据会导致panic: send on closed channel:

import "time"
func main() {
	go func() {
		time.Sleep(time.Hour)
	}()
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	c <- 3
}
但是从这个关闭的channel中不但可以读取出已发送的数据，还可以不断的读取零值:


c := make(chan int, 10)
c <- 1
c <- 2
close(c)
fmt.Println(<-c) //1
fmt.Println(<-c) //2
fmt.Println(<-c) //0
fmt.Println(<-c) //0
但是如果通过range读取，channel关闭后for循环会跳出：


c := make(chan int, 10)
c <- 1
c <- 2
close(c)
for i := range c {
	fmt.Println(i)
}
通过i, ok := <-c可以查看Channel的状态，判断值是零值还是正常读取的值。


c := make(chan int, 10)
close(c)
i, ok := <-c
fmt.Printf("%d, %t", i, ok) //0, false
