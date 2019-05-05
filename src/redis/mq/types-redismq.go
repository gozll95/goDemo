
import "github.com/go-redis/redis"

// Queue
type Queue struct {
	redisClient    *redis.Client
	Name           string
	rateStatsCache map[int64]map[string]int64
	rateStatsChan  chan (*dataPoint)
	lastStatsWrite int64
}

type dataPoint struct {
	name  string
	value int64
	incr  bool
}


// 每秒在内存里,过了1s就刷到redis里
func (queue *Queue) startStatsWriter() {
	queue.rateStatsCache = make(map[int64]map[string]int64)
	queue.rateStatsChan = make(chan *dataPoint, 2E6)
	writing := false
	go func() {
		for dp := range queue.rateStatsChan {
			now := time.Now().UTC().Unix()
			if queue.rateStatsCache[now] == nil {
				queue.rateStatsCache[now] = make(map[string]int64)
			}
			queue.rateStatsCache[now][dp.name] += dp.value
			if now > queue.lastStatsWrite && !writing {
				writing = true
				queue.writeStatsCacheToRedis(now)
				writing = false
			}
		}
	}()
	return
}

// Consumer
type Consumer struct {
	Name  string
	Queue *Queue

	cancel         context.CancelFunc
	contextCleared <-chan struct{}
}


func (consumer *Consumer) startHeartbeat() {
	firstWrite := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	consumer.cancel = cancel

	waitForClear := make(chan struct{}, 1)
	consumer.contextCleared = waitForClear

	go func() {
		firstRun := true
		for {
			consumer.Queue.redisClient.Set(
				consumerHeartbeatKey(consumer.Queue.Name, consumer.Name),
				"ping",
				time.Second,
			)
			if firstRun {
				// use close instead
				close(firstWrite)
				firstRun = false
			}
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				// remove heart beat immediately
				consumer.Queue.redisClient.Del(consumerHeartbeatKey(consumer.Queue.Name, consumer.Name))
				close(waitForClear)
				return
			}
		}
	}()
	<-firstWrite // 确保已经set heartbeat key了
	return
}



type Package struct {
	Payload    string
	CreatedAt  time.Time
	Queue      interface{} `json:"-"`
	Consumer   *Consumer   `json:"-"`
	Collection *[]*Package `json:"-"`
	Acked      bool        `json:"-"`
	//TODO add Headers or smth. when needed
	//wellle suggested error headers for failed packages
}



func (consumer *Consumer) Quit() {
	if consumer.cancel == nil {
		return
	}

	consumer.cancel()
	// wait until heart beat mark is removed
	<-consumer.contextCleared

	consumer.cancel = nil
}


// 十分牛掰
// FlushBuffer tells the background writer to flush the buffer to redis
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