// 另外一个chanchan的例子参考gpool-weixin-blog.go

package main

import ()

func (queue *BufferedQueue) FlushBuffer() {
	flushing := make(chan bool, 1)
	queue.flushStatus <- flushing
	queue.flushCommand <- true
	<-flushing
	return
}

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
				for i := 0; i < len(queue.flushStatus); i++ {
					c := <-queue.flushStatus
					c <- true
				}
				queue.nextWrite = time.Now().Unix() + 1
			}
			<-queue.flushCommand
		}
	}()
}
