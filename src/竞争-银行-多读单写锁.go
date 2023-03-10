package bank

import (
	"sync"
)

var mu sync.RWMutex
var balance int

func Balance() int {
	mu.RLock() //readers lock
	defer mu.RUnlock()
	return balance
}

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance < 0 {
		deposit(amount)
		return false //insufficient funds
	}
	return true
}

//This function requires that the lock be held.
func deposit(amount int) { balance += amount }

/*
由于Balance函数只需要读取变量的状态，所以我们同时让多个Balance调用并发运行事实上是安全的，
只要在运行的时候没有存款或者取款操作就行。在这种场景下我们需要一种特殊类型的锁，其允许多个
只读操作并行执行，但写操作会完全互斥。这种锁叫作“多读单写”锁(multiple readers, single
writer lock)，Go语言提供的这样的锁是sync.RWMutex：
*/

/*
Balance函数现在调用了RLock和RUnlock方法来获取和释放一个读取或者共享锁。Deposit函数没有变化，会调用mu.Lock和mu.Unlock方法来获取和释放一个写或互斥锁。

在这次修改后，Bob的余额查询请求就可以彼此并行地执行并且会很快地完成了。锁在更多的时间范围可用，并且存款请求也能够及时地被响应了。

RLock只能在临界区共享变量没有任何写入操作时可用。一般来说，我们不应该假设逻辑上的只读函数/方法也不会去更新某一些变量。比如一个方法功能是访问一个变量，但它也有可能会同时去给一个内部的计数器+1(译注：可能是记录这个方法的访问次数啥的)，或者去更新缓存--使即时的调用能够更快。如果有疑惑的话，请使用互斥锁。

RWMutex只有当获得锁的大部分goroutine都是读操作，而锁在竞争条件下，也就是说，goroutine们必须等待才能获取到锁的时候，RWMutex才是最能带来好处的。RWMutex需要更复杂的内部记录，所以会让它比一般的无竞争锁的mutex慢一些。
*/
