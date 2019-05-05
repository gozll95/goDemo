# redis设计list/set
```
const (
	connectionsKey                   = "rmq::connections"                                           // Set of connection names
	connectionHeartbeatTemplate      = "rmq::connection::{connection}::heartbeat"                   // expires after {connection} died
	connectionQueuesTemplate         = "rmq::connection::{connection}::queues"                      // Set of queues consumers of {connection} are consuming
	connectionQueueConsumersTemplate = "rmq::connection::{connection}::queue::[{queue}]::consumers" // Set of all consumers from {connection} consuming from {queue}
	connectionQueueUnackedTemplate   = "rmq::connection::{connection}::queue::[{queue}]::unacked"   // List of deliveries consumers of {connection} are currently consuming

	queuesKey             = "rmq::queues"                     // Set of all open queues
	queueReadyTemplate    = "rmq::queue::[{queue}]::ready"    // List of deliveries in that {queue} (right is first and oldest, left is last and youngest)
	queueRejectedTemplate = "rmq::queue::[{queue}]::rejected" // List of rejected deliveries from that {queue}

	phConnection = "{connection}" // connection name
	phQueue      = "{queue}"      // queue name
	phConsumer   = "{consumer}"   // consumer name (consisting of tag and token)

	defaultBatchTimeout = time.Second
	purgeBatchSize      = 100
)
```

# 实例:
- Connection:维护多个Queue
- Queue:相当于一个中间层
    - 当作为`Producer`的时候,push Delivery `ready`队列
    - 当作为`Consumer`的时候,
        - 1)从`ready`队列里BRPOPLUSH delivery 到`unacked`队列,
        - 2)并且将delivery放入chan中
        - 3)`AddConsumer`来range delivery消费chan中的delivery[在这里,queue有点类似调度者]
- Stat:监控各种队列/set的状态
- Cleaner:清空`无心跳`的queue,并将queue里面`unacked`/`rejected`的**return**回`read队列`
- Delivery:
    - Ack: 会将unacked list里对应的delivery删除
    - Reject: 会将delivery push进 reject list


# 流程:
```
Producer: Open a Connection[
            1.将connection name add into `SET conn`
            2.go connection.heartbeat[每隔一段时间 set hearbeatkey 1 XXs]
            ]----->connection Open a queue-------> queue.Publish(xxx string)[LPush into `ready list`]
Consumer: Open a Connection[
            1.将connection name add into `SET conn`
            2.go connection.heartbeat[每隔一段时间 set hearbeatkey 1 XXs]
            ]----->
         connection Open a queue-------> 
         start Consuming[
              1.从`ready`队列里BRPOPLUSH delivery 到`unacked`队列,
              2)并且将delivery放入chan中
        ] ----->
        queue.AddConsumer(xx consumer[实现了Consume(delivery rmq.Delivery)接口])[
                1.add consumer to consumer-set
                2.`go func` range delivery chan then consumer delivery
        ]
        或者
        queue.AddBatchConsumer(xx consumer[实现了Consume(delivery rmq.Delivery)接口])[ // 批量消费delivery
                1.add consumer to consumer-set
                2.`go func` 每隔xx时间或者达到了batchSize 就 从 chan delivery 出来 一堆 deliverys
        ]
Consumer的Consume delivery方法中 调用 delivery.Ack()/Reject()
```


# 精致代码
```
//每隔xx时间或者达到了batchSize 就 从 chan delivery 出来 一堆 deliverys

func (queue *redisQueue) consumerConsume(consumer Consumer) {
	for delivery := range queue.deliveryChan {
		// debug(fmt.Sprintf("consumer consume %s %s", delivery, consumer)) // COMMENTOUT
		consumer.Consume(delivery)
	}
}

func (queue *redisQueue) consumerBatchConsume(batchSize int, timeout time.Duration, consumer BatchConsumer) {
	batch := []Delivery{}
	for {
		// Wait for first delivery
		delivery, ok := <-queue.deliveryChan
		if !ok {
			// debug("batch channel closed") // COMMENTOUT
			return
		}
		batch = append(batch, delivery)
		// debug(fmt.Sprintf("batch consume added delivery %d", len(batch))) // COMMENTOUT
		batch, ok = queue.batchTimeout(batchSize, batch, timeout)
		consumer.Consume(batch)
		if !ok {
			// debug("batch channel closed") // COMMENTOUT
			return
		}
		batch = batch[:0] // reset batch
	}
}

func (queue *redisQueue) batchTimeout(batchSize int, batch []Delivery, timeout time.Duration) (fullBatch []Delivery, ok bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// debug("batch timer fired") // COMMENTOUT
			// debug(fmt.Sprintf("batch consume consume %d", len(batch))) // COMMENTOUT
			return batch, true
		case delivery, ok := <-queue.deliveryChan:
			if !ok {
				// debug("batch channel closed") // COMMENTOUT
				return batch, false
			}
			batch = append(batch, delivery)
			// debug(fmt.Sprintf("batch consume added delivery %d", len(batch))) // COMMENTOUT
			if len(batch) >= batchSize {
				// debug(fmt.Sprintf("batch consume wait %d < %d", len(batch), batchSize)) // COMMENTOUT
				return batch, true
			}
		}
	}
}
```

```
// 心跳

go connection.heartbeat()

// heartbeat keeps the heartbeat key alive
func (connection *redisConnection) heartbeat() {
	for {
		if !connection.updateHeartbeat() {
			// log.Printf("rmq connection failed to update heartbeat %s", connection)
		}

		time.Sleep(time.Second)

		if connection.heartbeatStopped {
			// log.Printf("rmq connection stopped heartbeat %s", connection)
			return
		}
	}
}


func (connection *redisConnection) updateHeartbeat() bool {
	ok := connection.redisClient.Set(connection.heartbeatKey, "1", heartbeatDuration)
	return ok
}
```

```
// 封装
type Consumer interface {
	Consume(delivery Delivery)
}

type ConsumerFunc func(Delivery)

func (consumerFunc ConsumerFunc) Consume(delivery Delivery) {
	consumerFunc(delivery)
}


type Handler struct {
	connection rmq.Connection
}

func NewHandler(connection rmq.Connection) *Handler {
	return &Handler{connection: connection}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	layout := request.FormValue("layout")
	refresh := request.FormValue("refresh")

	queues := handler.connection.GetOpenQueues()
	stats := handler.connection.CollectStats(queues)
	log.Printf("queue stats\n%s", stats)
	fmt.Fprint(writer, stats.GetHtml(layout, refresh))
}
```