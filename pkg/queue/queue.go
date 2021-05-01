package queue

import (
	"fmt"
)

type Queue struct {
	current int
	max     int
	arr     []interface{}
}

type QueueInterface interface {
	Queue(e interface{}) error
	Dequeue() (interface{}, error)
	IsEmpty() bool
	IsFull() bool
	Len() int
}

func NewQueue(max int) *Queue {
	q := &Queue{}
	q.max = max
	q.current = -1
	q.arr = make([]interface{}, max)
	return q
}

func (q *Queue) Queue(elem interface{}) error {
	if q.IsFull() {
		return fmt.Errorf("queue max exceeded")
	}
	q.current++
	q.arr[q.current] = elem
	return nil
}

func (q *Queue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue is empty")
	}
	out := q.arr[q.current]
	q.current--
	return out, nil
}

func (q *Queue) IsEmpty() bool {
	if q.current < 0 {
		return true
	}
	return false
}

func (q *Queue) IsFull() bool {
	if q.current >= q.max-1 {
		return true
	}
	return false
}

func (q *Queue) Len() int {
	return len(q.arr)
}
