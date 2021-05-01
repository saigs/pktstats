package queue

import (
	"math/rand"
	"testing"
	"time"
)

func TestQueueAddRemove(t *testing.T) {
	var (
		q = NewQueue(10)
	)
	t.Logf("testing queue add remove\n")
	t.Logf("is queue empty\n")
	if !q.IsEmpty() {
		t.Fatalf("queue not empty!")
	}
	t.Logf("add to queue\n")
	rand.Seed(time.Now().UTC().UnixNano())
	for !q.IsFull() {
		q.Queue(rand.Intn(100))
	}
	t.Logf("remove from queue\n")
	for !q.IsEmpty() {
		q.Dequeue()
	}
	t.Logf("is queue empty\n")
	if !q.IsEmpty() {
		t.Fatalf("queue not empty!")
	}
}

func TestQueueVarious(t *testing.T) {
	var (
		s = "string"
		i = 100
		b = rune('x')
		p = &s
	)
	var (
		q = NewQueue(10)
	)
	add := func(e interface{}) {
		if err := q.Queue(e); err != nil {
			t.Fatalf("queue error: %v", err)
		}
	}
	t.Logf("testing queue add remove special struct\n")
	add(s)
	add(i)
	add(b)
	add(p)
	for !q.IsEmpty() {
		q.Dequeue()
	}
	if !q.IsEmpty() {
		t.Fatalf("queue not empty!")
	}
}

func TestQueueStruct(t *testing.T) {
	type MyData struct {
		i int
		r rune
		s string
		a []int
	}
	var (
		q = NewQueue(10)
		m = MyData{
			i: 10,
			r: 'x',
			s: "string",
		}
	)
	q.Queue(m)
	for !q.IsEmpty() {
		q.Dequeue()
	}
	if !q.IsEmpty() {
		t.Fatalf("queue not empty!")
	}
}
