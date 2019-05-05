package bank

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

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance() int {
	mu.Lock
	defer mu.Unlock()
	return balance
}

//This function requires that the lock be held.
func deposit(amount int) { balance += amount }

/*
一个通用的解决方案是将一个函数分离为多个函数，比如我们把Deposit分离成两个：
一个不导出的函数deposit，这个函数假设锁总是会被保持并去做实际的操作，
另一个是导出的函数Deposit，这个函数会调用deposit，但在调用前会先去获取锁。同理我们可以将Withdraw也表示成这种形式：
*/
