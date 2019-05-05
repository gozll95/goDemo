package send

import (
	"sync"

	"github.com/qbusapi/app/models"
)

type MessageQueue struct {
	First *MessageNode
	Tail  *MessageNode
	mutex sync.Mutex
}

type MessageNode struct {
	Message *models.MessageModel
	Next    *MessageNode
}

func (queue *MessageQueue) IsEmpty() bool {
	if queue.First == nil || queue.Tail == nil {
		return true
	}

	return false
}

func (queue *MessageQueue) AddMessage(message *models.MessageModel) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	node := &MessageNode{
		Message: message,
		Next:    nil,
	}

	// empty queue
	if queue.IsEmpty() {
		queue.First = node
		queue.Tail = node
		return
	}

	queue.Tail.Next = node
	queue.Tail = node
}

func (queue *MessageQueue) FetchMessage() (message *models.MessageModel) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	// empty queue
	if queue.IsEmpty() {
		return nil
	}

	message = queue.First.Message
	queue.First = queue.First.Next
	return
}
