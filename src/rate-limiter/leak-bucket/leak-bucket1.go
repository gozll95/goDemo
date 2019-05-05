package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	rate       float64 //固定每秒出水速率
	capacity   float64 //桶的容量
	water      float64 //桶中当前水量
	lastLeakMs int64   //桶上次漏水时间戳 ms

	lock sync.Mutex
}

func (l *LeakyBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1e6
	eclipse := float64((now - l.lastLeakMs)) * l.rate / 1000 //先执行漏水
	l.water = l.water - eclipse                              //计算剩余水量
	l.water = math.Max(0, l.water)                           //桶干了
	l.lastLeakMs = now
	if (l.water + 1) < l.capacity {
		// 尝试加水,并且水还未满
		l.water++
		return true
	} else {
		// 水满，拒绝加水
		return false
	}
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}

func main() {
	var wg sync.WaitGroup
	var lr LeakyBucket
	lr.Set(3, 3) //每秒访问速率限制为3个请求，桶容量为3

	for i := 0; i < 10; i++ {
		wg.Add(1)

		fmt.Println("Create req", i, time.Now())
		go func(i int) {
			if lr.Allow() {
				fmt.Println("Respon req", i, time.Now())
			}
			wg.Done()
		}(i)

		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}
