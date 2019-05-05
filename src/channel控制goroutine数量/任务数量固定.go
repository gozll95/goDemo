package main

func Afunction(ch chan int) {
	ch <- 1
}

func main() {
	var (
		ch        chan int = make(chan int, 20) //可以同时运行的routine数量为20
		dutycount int      = 500
	)
	for i := 0; i < dutycount; i++ {
		go Afunction(ch)
	}

	//知道了任务总量，可以像这样利用固定循环次数的循环检测所有的routine是否工作完毕
	for i := 0; i < dutycount; i++ {
		<-ch
	}
}
