同步

channel可以用在goroutine之间的同步。
下面的例子中main goroutine通过done channel等待worker完成任务。 worker做完任务后只需往channel发送一个数据就可以通知main goroutine任务完成。


import (
	"fmt"
	"time"
)
func worker(done chan bool) {
	time.Sleep(time.Second)
	// 通知任务已完成
	done <- true
}
func main() {
	done := make(chan bool, 1)
	go worker(done)
	// 等待任务完成
	<-done
}