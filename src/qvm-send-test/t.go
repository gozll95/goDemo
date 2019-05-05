package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/zhu/qvm/server/lib/aliyun/mns/utils"

	"github.com/astaxie/beego"
	"github.com/zhu/qvm/server/lib/aliyun/mns/manager"
)

//m[0]->gid
//m[1]->timeout
func DoWithRandomTicket(t manager.Manager, args ...interface{}) {
	var (
		timeout     <-chan time.Time
		gId         int
		isSuccess   bool //true->update ticket success;false -> update ticket false
		ok          bool
		timeoutChan chan struct{}
	)

	timeoutChan = make(chan struct{})

	//get gId
	if gId, ok = args[0].(int); !ok {
		return
	}
	fmt.Println("this is goroutine ", gId)

	//get timeout
	if managerTimeout, ok := args[1].(time.Duration); ok {
		timeout = time.After(managerTimeout)
	} else {
		return
	}

	items, err := getRandomTicket() //ge Ticket and set it free--busy
	beego.Info(gId, "items: ", items)
	if err != nil {
		//包含not found
		return
	}

	defer func() {
		//需要考虑意外停止的时候
		if err := recover(); err != nil {
			time.Sleep(1 * time.Second)
			beego.Info("panic")
			fmt.Println(gId, "update ticket idle ")
			return
		}
		if isSuccess {
			fmt.Println(gId, "update ticket success ")
		} else {
			fmt.Println(gId, "update ticket idle ")
		}
		return
	}()

	if len(items) == 0 {
		select {
		case <-t.Pause():
		default:
			fmt.Println("start send pause sig")
			t.Pause() <- struct{}{}
			fmt.Println("finish send pause sig")
		}
		return
	}

	select {
	case result := <-generateResult(t, gId, items, timeoutChan): //这里do getTicket 并且返回ticket number
		if result {
			isSuccess = true
			return
		}
		isSuccess = false
		return
	case <-timeout:
		beego.Info(gId, " timeout")
		isSuccess = false
		close(timeoutChan)
		return
	case <-t.Quit():
		beego.Info(gId, " receive close signal")
		isSuccess = false
		// panic(closed)
		return
	}

}

func DoWithEmptyTicket(t manager.Manager, args ...interface{}) {
	var (
		timeout     <-chan time.Time
		gId         int
		isSuccess   bool //true->update ticket success;false -> update ticket false
		ok          bool
		timeoutChan chan struct{}
	)

	timeoutChan = make(chan struct{})

	//get gId
	if gId, ok = args[0].(int); !ok {
		return
	}
	fmt.Println("this is goroutine ", gId)

	//get timeout
	if managerTimeout, ok := args[1].(time.Duration); ok {
		timeout = time.After(managerTimeout)
	} else {
		return
	}

	ticketId, items, err := getEmptyTicket() //ge Ticket and set it free--busy
	beego.Info(gId, "items: ", items)
	if err != nil {
		//包含not found
		return
	}

	defer func() {
		//需要考虑意外停止的时候
		if err := recover(); err != nil {
			time.Sleep(1 * time.Second)
			beego.Info("panic")
			fmt.Println(gId, "update ticket idle ", ticketId)
			return
		}
		if isSuccess {
			fmt.Println(gId, "update ticket success ", ticketId)
		} else {
			fmt.Println(gId, "update ticket idle ", ticketId)
		}
		return
	}()

	if len(items) == 0 {
		select {
		case <-t.Pause():
		default:
			fmt.Println("start send pause sig")
			t.Pause() <- struct{}{}
			fmt.Println("finish send pause sig")
		}
		return
	}

	select {
	case result := <-generateResult(t, gId, items, timeoutChan): //这里do getTicket 并且返回ticket number
		if result {
			isSuccess = true
			return
		}
		isSuccess = false
		return
	case <-timeout:
		beego.Info(gId, " timeout")
		isSuccess = false
		close(timeoutChan)
		return
	case <-t.Quit():
		beego.Info(gId, " receive close signal")
		isSuccess = false
		// panic(closed)
		return
	}

}

func DoWithErrTicket(t manager.Manager, args ...interface{}) {
	var (
		timeout     <-chan time.Time
		gId         int
		isSuccess   bool //true->update ticket success;false -> update ticket false
		ok          bool
		timeoutChan chan struct{}
	)
	timeoutChan = make(chan struct{})

	//get gId
	if gId, ok = args[0].(int); !ok {
		return
	}
	fmt.Println("this is goroutine ", gId)

	//get timeout
	if managerTimeout, ok := args[1].(time.Duration); ok {
		timeout = time.After(managerTimeout)
	} else {
		return
	}

	ticketId, items, err := getErrTicket() //ge Ticket and set it free--busy

	if err != nil {
		//包含not found
		beego.Error(gId, "err: ", err)
		return
	}
	beego.Info(gId, "items: ", items)

	defer func() {
		//需要考虑意外停止的时候
		if err := recover(); err != nil {
			time.Sleep(1 * time.Second)
			beego.Info("panic")
			fmt.Println(gId, "update ticket idle ", ticketId)
			return
		}
		if isSuccess {
			fmt.Println(gId, "update ticket success ", ticketId)
		} else {
			fmt.Println(gId, "update ticket idle ", ticketId)
		}
		return
	}()

	if len(items) == 0 {
		select {
		case <-t.Pause():
		default:
			fmt.Println("start send pause sig")
			t.Pause() <- struct{}{}
			fmt.Println("finish send pause sig")
		}
		return
	}
	select {
	case result := <-generateResult(t, gId, items, timeoutChan): //这里do getTicket 并且返回ticket number
		if result {
			isSuccess = true
			return
		}
		isSuccess = false
		return
	case <-timeout:
		beego.Info(gId, " timeout")
		isSuccess = false
		close(timeoutChan)
		return
	case <-t.Quit():
		beego.Info(gId, " receive close signal")
		isSuccess = false
		// panic(closed)
		return
	}

}

func DoWithFullTicket(t manager.Manager, args ...interface{}) {
	var (
		timeout     <-chan time.Time
		gId         int
		isSuccess   bool //true->update ticket success;false -> update ticket false
		ok          bool
		timeoutChan chan struct{}
	)
	timeoutChan = make(chan struct{})

	//get gId
	if gId, ok = args[0].(int); !ok {
		return
	}
	fmt.Println("this is goroutine ", gId)

	//get timeout
	if managerTimeout, ok := args[1].(time.Duration); ok {
		timeout = time.After(managerTimeout)
	} else {
		return
	}

	ticketId, items, err := getFullTicket() //ge Ticket and set it free--busy
	beego.Info(gId, "items: ", items)
	if err != nil {
		//包含not found
		return
	}

	defer func() {
		//需要考虑意外停止的时候
		if err := recover(); err != nil {
			time.Sleep(1 * time.Second)
			beego.Info("panic")
			fmt.Println(gId, "update ticket idle ", ticketId)
			return
		}
		if isSuccess {
			fmt.Println(gId, "update ticket success ", ticketId)
		} else {
			fmt.Println(gId, "update ticket idle ", ticketId)
		}
		return
	}()

	if len(items) == 0 {
		select {
		case <-t.Pause():
		default:
			fmt.Println("start send pause sig")
			t.Pause() <- struct{}{}
			fmt.Println("finish send pause sig")
		}
		return
	}

	select {
	case result := <-generateResult(t, gId, items, timeoutChan): //这里do getTicket 并且返回ticket number
		if result {
			isSuccess = true
			return
		}
		isSuccess = false
		return
	case <-timeout:
		beego.Info(gId, " timeout")
		isSuccess = false
		close(timeoutChan)
		return
	case <-t.Quit():
		beego.Info(gId, " receive close signal")
		isSuccess = false
		// panic(closed)
		return
	}
}

func doSomething(t manager.Manager, gId, id int, status chan bool, timeoutChan chan struct{}) {
	select {
	case <-timeoutChan:
		return
	default:
		timeSleep := utils.GetRand()
		fmt.Printf("gId:%d-taskId:%d-sleep %d s\n", gId, id, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
		fmt.Printf("gId:%d-taskId:%d-do\n", gId, id)
		status <- true
	}
}

func generateResult(t manager.Manager, gId int, items []int, timeoutChan chan struct{}) <-chan bool {
	cc := make(chan bool, 1)

	go func() {
		var result bool
		result = true
		if len(items) > 0 {
			status := make(chan bool, len(items))
			for _, item := range items {
				select {
				case <-timeoutChan:
					beego.Error(gId, " gid timeout")
					return
				default:
					go doSomething(t, gId, item, status, timeoutChan)
				}

			}
			for i := 0; i < len(items); i++ {
				c := <-status
				if !c {
					result = false
				}
			}
			beego.Info(gId, " gid result", result)
			cc <- result
		} else {
			select {
			case <-t.Pause():
			default:
				beego.Error("start send pause sig")
				t.Pause() <- struct{}{}
				beego.Error("finish send pause sig")
			}
			result = false
			cc <- false
		}
	}()
	return cc
}

func getRandomTicket() (items []int, err error) {
	x := utils.GetRand()
	if x%2 == 0 {
		for i := 0; i < 100; i++ {
			items = append(items, i)
		}
		return items, nil
	}
	return []int{}, nil
}

func getEmptyTicket() (ticketId int, items []int, err error) {
	return 0, []int{}, nil
}

func getErrTicket() (ticketId int, items []int, err error) {
	return 0, []int{}, errors.New("xxx")
}

func getFullTicket() (ticketId int, items []int, err error) {
	for i := 0; i < 100; i++ {
		items = append(items, i)
	}
	return 0, items, nil
}
