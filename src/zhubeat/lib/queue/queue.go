package queue

import (
	"container/list"
	"errors"
)

type Queue struct {
	data *list.List
}

func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

func (q *Queue) Push(v interface{}) {
	q.data.PushFront(v)
}

func (q *Queue) Pop() (interface{}, error) {
	iter := q.data.Back()
	if iter == nil {
		return nil, errors.New("no value in queue")
	}
	v := iter.Value
	q.data.Remove(iter)
	return v, nil
}

func (q *Queue) Len() int {
	return q.data.Len()
}

func (q *Queue) Dump() []interface{} {
	var res []interface{}
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		res = append(res, iter.Value)
	}
	return res
}

func (q *Queue) Reset() {
	q.data.Init()
}
