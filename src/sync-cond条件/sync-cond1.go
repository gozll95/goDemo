条件等待

条件等待通过Wait让goroutine等待,通过signal让一个等待的goroutine继续,通过broadcase让所有等待的goroutine继续。
在Wait之前应该手动为c.L上锁,Wait结束后手动解锁。为了避免虚假唤醒,需要将"Wait放到一个条件判断循环中"。

c.L.Lock()
for !condition(){
	c.Wait()
}
// 执行条件满足之后的动作
c.L.Unlock()

Cond 在开始使用之后，不能再被复制。

---------------------------------
type Cond struct {
    L Locker // 在“检查条件”或“更改条件”时 L 应该锁定。
} 

// 创建一个条件等待
func NewCond(l Locker) *Cond

// Broadcast 唤醒所有等待的 Wait，建议在“更改条件”时锁定 c.L，更改完毕再解锁。
func (c *Cond) Broadcast()

// Signal 唤醒一个等待的 Wait，建议在“更改条件”时锁定 c.L，更改完毕再解锁。
func (c *Cond) Signal()

// Wait 会解锁 c.L 并进入等待状态，在被唤醒时，会重新锁定 c.L
func (c *Cond) Wait()

func main(){
	condition:=false //条件不满足
	var mu sync.Mutex 
	cond:=sync.NewCond(&mu)
	// 当goroutine去创造条件
	go func(){
		mu.Lock()
		condition=true //更改条件
		cond.Signal() // 发送通知:条件已满足
		mu.Unlock()
	}()
	mu.Lock()
	// 检查条件是否满足,避免虚假通知,同时避免Signal提前于Wait执行。
	for !condition{
		// 等待条件满足的通知,如果收到虚假通知,则循环继续等待。
		cond.Wait() //等待时mu处于解锁状态,唤醒时重新锁定。
	}
	fmt.Println("条件满足,开始后续动作...")
	mu.Unlock()
}


