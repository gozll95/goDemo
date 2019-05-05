package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhu/qvm/server/lib/aliyun/mns/manager"
)

func main() {
	// abort := make(chan struct{})
	// go func() {
	// 	os.Stdin.Read(make([]byte, 1)) // read a single byte
	// 	abort <- struct{}{}
	// }()

	// m, err := send.NewSendManager(2, 3*time.Second, 60*time.Second, DoWithRandomTicket)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// go func() {
	// 	time.Sleep(20 * time.Second)
	// 	m.Close()
	// }()

	// go func() {
	// 	sigs := make(chan os.Signal, 1)
	// 	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	// 	for {
	// 		select {
	// 		case <-sigs:
	// 			m.Close()
	// 			os.Exit(-1)
	// 			abort <- struct{}{}
	// 		}
	// 	}

	// }()

	// m.Start()
	// <-abort

	abort := make(chan struct{})

	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	m, err := manager.NewMnsManager(100, 20*time.Second, 5*time.Second, 1*time.Second, DoWithEmptyTicket)
	if err != nil {
		return
	}

	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	m.Close()
	// }()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

		for {
			select {
			case <-sigs:
				m.Close()
				os.Exit(-1)
				abort <- struct{}{}
			}
		}

	}()

	m.Start()
	<-abort

}
