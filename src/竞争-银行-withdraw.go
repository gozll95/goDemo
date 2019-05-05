package main

import (
	"fmt"
)

var deposits = make(chan int)
var balances = make(chan int)
var withdraws = make(chan int)
var result = make(chan bool)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) string {
	withdraws <- amount
	select {
	case a := <-result:
		if a {
			//成功
			return fmt.Sprintf("withdraw:%d,remain balance is:%d,result:sucessful", amount, Balance())
		} else {
			//失败
			return fmt.Sprintf("withdraw:%d,remain balance is:%d,result:failed", amount, Balance())
		}
	}

}

//balance这个变量限制在了一个goroutine里
func teller() {
	var balnace int
	for {
		select {
		case amount := <-deposits:
			balnace += amount
		case balances <- balnace:
		case withdraw := <-withdraws:
			balnace -= withdraw
			if balnace < 0 {
				balnace -= withdraw
				result <- false
			} else {
				result <- true
			}
		}
	}
}

func main() {
	done := make(chan struct{})
	go teller()

	//Alice
	go func() {
		Deposit(200)
		fmt.Println("balance=", Balance())
		done <- struct{}{}
	}()

	//Bob
	go func() {
		Deposit(100)
		fmt.Println("balance=", Balance())
		done <- struct{}{}
	}()

	//Cindy
	go func() {
		fmt.Println("balance=", Balance())
		fmt.Println(Withdraw(150))
		done <- struct{}{}
	}()

	//Wait for both transacions
	<-done
	<-done
	<-done

	want := 100
	got := Balance()

	fmt.Printf("Balance = %d, want %d\n", got, want)

}

/*
由于其他goroutine不能够直接访问变量,他们只能通过一个channel来发送给
指定的goroutine来查询更新变量。
也就是GO的口头禅"不要使用共享数据来通信,使用通信来共享数据"。
一个提供对一个指定的变量通过channel来请求的gotoutine叫做这个变量的监控(monitor)goroutine。
*/
