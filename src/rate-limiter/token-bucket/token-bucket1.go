package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //桶上次放token的时间戳 s

	lock sync.Mutex
}

func (l *TokenBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().Unix()
	l.tokens = l.tokens + (now-l.lastTokenSec)*l.rate // 先添加令牌
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	if l.tokens > 0 {
		// 还有令牌，领取令牌
		l.tokens--
		return true
	} else {
		// 没有令牌,则拒绝
		return false
	}
}

func (l *TokenBucket) Set(r, c int64) {
	l.rate = r
	l.capacity = c
	l.tokens = 0
	l.lastTokenSec = time.Now().Unix()
}

func main() {
	var wg sync.WaitGroup
	var lr TokenBucket
	lr.Set(3, 3) //每秒访问速率限制为3个请求，桶容量为3

	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)

		fmt.Println("Create req", i, time.Now())
		go func(i int) {
			if lr.Allow() {
				fmt.Println("Respon req", i, time.Now())
			}
			wg.Done()
		}(i)

		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}

/*
类似的，这里初始化桶容量为 3 个单位，桶内无令牌，每 1s 产生 3 个令牌，主程序阻塞 1s 以便让桶中储备好 3 个令牌，而后每 1s 创建 5 个请求，获取访问令牌：


Create req 0 2018-09-12 20:56:30.926267 +0800 CST m=+1.002077637
Respon req 0 2018-09-12 20:56:30.926512 +0800 CST m=+1.002322722
Create req 1 2018-09-12 20:56:31.127552 +0800 CST m=+1.203356343
Respon req 1 2018-09-12 20:56:31.127643 +0800 CST m=+1.203445860
Create req 2 2018-09-12 20:56:31.328242 +0800 CST m=+1.404040714
Respon req 2 2018-09-12 20:56:31.328311 +0800 CST m=+1.404110620
Create req 3 2018-09-12 20:56:31.529957 +0800 CST m=+1.605746538
Respon req 3 2018-09-12 20:56:31.530069 +0800 CST m=+1.605862175
Create req 4 2018-09-12 20:56:31.734506 +0800 CST m=+1.810291673 x
Create req 5 2018-09-12 20:56:31.938578 +0800 CST m=+2.014347368 x
Create req 6 2018-09-12 20:56:32.141289 +0800 CST m=+2.217061017
Respon req 6 2018-09-12 20:56:32.141379 +0800 CST m=+2.217153670
Create req 7 2018-09-12 20:56:32.341492 +0800 CST m=+2.417260734
Respon req 7 2018-09-12 20:56:32.341571 +0800 CST m=+2.417339558
Create req 8 2018-09-12 20:56:32.543497 +0800 CST m=+2.619259765
Respon req 8 2018-09-12 20:56:32.543554 +0800 CST m=+2.619317332
Create req 9 2018-09-12 20:56:32.746344 +0800 CST m=+2.822096642 x
*/
