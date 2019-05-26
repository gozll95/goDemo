// 第一种
var add = make(chan struct{})
var total = make(chan int)

func getVistors() int {
	return <-total
}

func addVistors() {
	add <- struct{}{}
}

// 将全局变量限定在一个函数中执行,通过select case来实现串行化
func teller() {
	var visitors int
	for {
		select {
		case <-add:
			visitors += 1
		case total <- visitors:
		}
	}
}

// 第二种
var (
	c       = make(chan struct{}, 1)
	vistors int
)

func handleVistors() {
	c <- struct{}{}
	vistors++
	<-c
}