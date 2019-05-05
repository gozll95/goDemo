package main

import (
	"context"
	"log"
	"os"
	"time"
)

var logg *log.Logger

func someHandler() {
	//1.context.WithCancel
	//ctx, cancel := context.WithCancel(context.Background())
	//2.context.WithDeadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// go doStuff(ctx)
	go doTimeOutStuff(ctx)

	//10秒后取消doStuff
	time.Sleep(10 * time.Second)
	cancel()

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func doTimeOutStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)

		if deadline, ok := ctx.Deadline(); ok { //设置了deadl
			logg.Printf("deadline set")
			if time.Now().After(deadline) {
				logg.Printf(ctx.Err().Error())
				return
			}

		}

		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func main() {
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("end")
}

/*
context.WithCancel
result:
11:21:27 work
11:21:28 work
11:21:29 work
11:21:30 work
11:21:31 work
11:21:32 work
11:21:33 work
11:21:34 work
11:21:35 work
11:21:36 down
*/

/*
context.WithDeadline
result:
15:59:22 work
15:59:24 work
15:59:25 work
15:59:26 work
15:59:27 done
15:59:31 end
*/

/*
go doTimeOutStuff(ctx)
result:
11:32:55 deadline set
11:32:55 work
11:32:56 deadline set
11:32:56 work
11:32:57 deadline set
11:32:57 work
11:32:58 deadline set
11:32:58 work
11:32:59 deadline set
11:32:59 context deadline exceeded
11:33:04 end
*/
