package main

var (
	sema    = make(chan struct{}, 1) //a binary semaphore guarding balance
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} //acquire token
	balance += amount
	<-sema //release token
}

func Balance() int {
	sema <- struct{}{}
	b := balance
	<-sema
	return b
}


####################################
var (
    mu      sync.Mutex // guards balance
    balance int
)

func Deposit(amount int) {
    mu.Lock()
    balance = balance + amount
    mu.Unlock()
}

func Balance() int {
    mu.Lock()
    b := balance
    mu.Unlock()
    return b
}

=>

func Balance() int {
    mu.Lock()
    defer mu.Unlock()
    return balance
}