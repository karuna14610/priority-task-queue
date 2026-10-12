// Package queue implements a thread-safe priority queue for tasks.
// Higher Priority values are served first; tasks with equal priority
// are served in the order they were submitted (FIFO).
package queue

import (
	"container/heap"
	"errors"
	"sync"
)

// Task is a unit of work.
type Task struct {
	ID       string
	Priority int // higher number = more urgent
	Run      func() error

	seq   uint64 // submission order, used to break priority ties
	index int    // position in the heap (managed by container/heap)
}

// ErrClosed is returned when pushing to or popping from a closed, drained queue.
var ErrClosed = errors.New("queue: closed")

// taskHeap implements heap.Interface (max-heap on Priority, then FIFO).
type taskHeap []*Task

func (h taskHeap) Len() int { return len(h) }
func (h taskHeap) Less(i, j int) bool {
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	return h[i].seq < h[j].seq
}
func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *taskHeap) Push(x any) {
	t := x.(*Task)
	t.index = len(*h)
	*h = append(*h, t)
}
func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	t := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return t
}

// PriorityQueue is safe for concurrent use.
type PriorityQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  taskHeap
	nextID uint64
	closed bool
}

// New creates an empty queue.
func New() *PriorityQueue {
	q := &PriorityQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push adds a task to the queue.
func (q *PriorityQueue) Push(t *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrClosed
	}
	t.seq = q.nextID
	q.nextID++
	heap.Push(&q.items, t)
	q.cond.Signal()
	return nil
}

// Pop blocks until a task is available and returns the highest-priority one.
// After Close, Pop keeps returning remaining tasks, then ErrClosed.
func (q *PriorityQueue) Pop() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 && !q.closed {
		q.cond.Wait()
	}
	if len(q.items) == 0 {
		return nil, ErrClosed
	}
	return heap.Pop(&q.items).(*Task), nil
}

// Len returns the number of waiting tasks.
func (q *PriorityQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Close stops accepting new tasks and wakes all waiting consumers.
func (q *PriorityQueue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	q.cond.Broadcast()
}
