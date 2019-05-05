package queue

import (
	"zhubeat/lib/queue"
)

var (
	ToOutputQueue *queue.Queue
)

func init() {
	ToOutputQueue = queue.NewQueue()
}
