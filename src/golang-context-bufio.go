package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// ctx 一般用作第一个参数
func Inc(ctx context.Context, a int) {

	// 此处用 ctx 的 value 来获取 一个值。 不过 ctx 的 value 通常情况下不用来传递参数。
	// 此处只是用来说明 context.Value 的用法
	intival := ctx.Value("interval").(time.Duration)

	for {
		select {
		// 当子 ctx 的 cancelfunc 被调用的时候或者 context 到期， Done 关闭
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				fmt.Println("function canceled")
				return
			} else if ctx.Err() == context.DeadlineExceeded {
				// WithDeadline 或 WithTimeout 设定的超时被触发
				fmt.Println("function time out ")
				return
			}
		default:
			time.Sleep(intival)
			a++
		}
	}
}

func main() {
	dura := time.Second * 1
	ctx := context.WithValue(context.Background(), "interval", dura)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	// 先从 Background context 继承了一个带 value 的 context，再继承了一个带超时的 ctx。

	go Inc(ctx, 0) // 新开一个 goutine, 传入当前 goroutine 的上下文
	// context 是 goroutine 安全的， 可以在多个 goroutine 中同时访问该 context

	fmt.Print("Inc function is runnning")
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf(" %s second left to run the Inc function .\n", deadline.Sub(time.Now()).String())
	}
	// ctx.Deadline 主要获取ctx 的过期时间，如果没有设置超时的话 ok 会返回 false
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Printf("  Press A to abort: ")
		if line, err := reader.ReadString('\n'); err == nil {
			command := strings.TrimSuffix(line, "\n")
			switch command {
			case "A":
				cancel() // 手动调用 ctx 的cancelfunc ctx Done 返回的 channel 会关闭
				time.Sleep(500 * time.Millisecond)
				return
			default:
				if deadline, ok := ctx.Deadline(); ok {
					fmt.Printf("%s second left, Press A to stop `\n", deadline.Sub(time.Now()).String())
				}
			}
		}
	}
}
