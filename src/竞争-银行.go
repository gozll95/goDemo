package main

import (
	"fmt"
)

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

//balance这个变量限制在了一个goroutine里
func teller() {
	var balnace int
	for {
		select {
		case amount := <-deposits:
			balnace += amount
		case balances <- balnace:
		}
	}
}

func main() {
	done := make(chan struct{})
	go teller()

	//Alice
	go func() {
		Deposit(200)
		fmt.Println("=", Balance())
		done <- struct{}{}
	}()

	//Bob
	go func() {
		Deposit(100)
		done <- struct{}{}
	}()

	//Wait for both transacions
	<-done
	<-done

	if got, want := Balance(), 300; got != want {
		fmt.Errorf("Balance = %d, want %d", got, want)
	}

}

/*
由于其他goroutine不能够直接访问变量,他们只能通过一个channel来发送给
指定的goroutine来查询更新变量。
也就是GO的口头禅"不要使用共享数据来通信,使用通信来共享数据"。
一个提供对一个指定的变量通过channel来请求的gotoutine叫做这个变量的监控(monitor)goroutine。
*/
