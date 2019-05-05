# 设计
同rmq

# 精华
- 设计了local buffer,在内存中缓存多个package,超过一定size/超过一定时间才flush到redis里
- 同样也采用了heartbeat
- stat rates

# 分析local buffer
`这里讲channel用的挺好的`
CreateBufferedQueue 
->Start[
	queue.startHeartbeat()
	queue.startWritingBufferToRedis()
	queue.startPacemaker()
]
 

#### startHeartbeat
`这里通过bool+channel保证事件一定操作了`
 ```
 func (queue *BufferedQueue) startHeartbeat() {
	firstWrite := make(chan bool, 1) // 这个不是多余的,是确保已经set heartbeatkey了
	go func() {
		firstRun := true
		for {
			queue.redisClient.Set(queueHeartbeatKey(queue.Name), "ping", time.Second)
			if firstRun {
				firstWrite <- true
				firstRun = false
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	<-firstWrite 
	return
}
 ```


#### startWritingBufferToRedis()
`这里太巧妙了,使用了channel的len特性`
```
func (queue *BufferedQueue) startWritingBufferToRedis() {
	go func() {
		queue.nextWrite = time.Now().Unix()
		for {
			if len(queue.Buffer) >= queue.BufferSize || time.Now().Unix() >= queue.nextWrite {
				size := len(queue.Buffer)
				a := []string{}
				for i := 0; i < size; i++ {
					p := <-queue.Buffer
					a = append(a, p.getString())
				}
				queue.redisClient.LPush(queueInputKey(queue.Name), a...)
				queue.incrRate(queueInputRateKey(queue.Name), int64(size))
				for i := 0; i < len(queue.flushStatus); i++ { // 这里太巧妙了,使用了channel的len特性
					c := <-queue.flushStatus
					c <- true
				}
				queue.nextWrite = time.Now().Unix() + 1
			}
			<-queue.flushCommand
		}
	}()
}


func (queue *BufferedQueue) startPacemaker() {
	go func() {
		for {
			queue.flushCommand <- true
			time.Sleep(10 * time.Millisecond)
		}
	}()
}
```


####  FlushBuffer()
`巧妙使用了channel channel`
```
// FlushBuffer tells the background writer to flush the buffer to redis
func (queue *BufferedQueue) FlushBuffer() {
	flushing := make(chan bool, 1)
	queue.flushStatus <- flushing
	queue.flushCommand <- true
	<-flushing
	return
}
```


# 总结
`channel channel + 有缓冲的channel可以很灵活`