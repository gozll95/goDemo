
chan T          // 可以接收和发送类型为 T 的数据
chan<- float64  // 只可以用来发送 float64 类型的数据
<-chan int      // 只可以用来接收 int 类型的数据




for …… range语句可以处理Channel。

func main() {
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c <- i
		}
		close(c)
	}()
	for i := range c {
		fmt.Println(i)
	}
	fmt.Println("Finished")
}
range c产生的迭代值为Channel中发送的值，它会一直迭代直到channel被关闭。上面的例子中如果把close(c)注释掉，程序会一直阻塞在for …… range那一行。