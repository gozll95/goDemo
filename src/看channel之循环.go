package main

import (
	"gopl.io/ch8/thumbnail"
	"log"
	"os"
	"sync"
)

//无goroutine
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

//添加了goroutine
func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f)
			ch <- struct{}{}
		}(f)
	}
	for range filenames {
		<-ch
	}
}

//添加error channle,但是该会发生goroutine泄漏，造成oom
// 这个程序有一个微秒的bug。当它遇到第一个非nil的error时会直接将error返回到调用方，
// 使得没有一个goroutine去排空errors channel。这样剩下的worker goroutine在向这
// 个channel中发送值时，都会永远地阻塞下去，并且永远都不会退出。这种情况叫做goroutine
// 泄露(§8.4.4)，可能会导致整个程序卡住或者跑出out of memory的错误。

// 最简单的解决办法就是用一个具有合适大小的buffered channel，这样这些worker goroutine向
// channel中发送测向时就不会被阻塞。(一个可选的解决办法是创建一个另外的goroutine，当main goroutine
// 	返回第一个错误的同时去排空channel)
func makeThumbnails4(filenames []string) {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			error <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err
		}
	}
	return nil
}

//使用了带缓冲的
func makeThumbnails5(filenames []string) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))
	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}
	return thumbfiles, nil
}

//知道最后一个goroutine什么时候结束的
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup
	for f := range filenames {
		wg.Add(1)
		//worker
		go func(f string) {
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb)
			sizes <- info.Size()
		}(f)
	}

	//closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
