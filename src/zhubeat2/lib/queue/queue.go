package queue

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type Queue struct {
	data *list.List
	lock sync.RWMutex
}

func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

func (q *Queue) Push(v interface{}) {
	defer q.lock.Unlock()
	q.lock.Lock()
	q.data.PushFront(v)
}

func (q *Queue) Pop() (interface{}, error) {
	defer q.lock.Unlock()
	q.lock.Lock()
	iter := q.data.Back()
	if iter == nil {
		return nil, errors.New("no value in queue")
	}
	v := iter.Value
	q.data.Remove(iter)
	return v, nil
}

func (q *Queue) Len() int {
	defer q.lock.RUnlock()
	q.lock.RLock()
	return q.data.Len()
}

func (q *Queue) Dump() []interface{} {
	defer q.lock.Unlock()
	q.lock.Lock()

	var res []interface{}
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item:", iter.Value)
		res = append(res, iter.Value)
	}
	return res
}
