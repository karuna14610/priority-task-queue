package queue

import (
	"log"
	"sync"
)

// Pool runs tasks from a PriorityQueue using a fixed number of workers.
type Pool struct {
	q       *PriorityQueue
	workers int
	wg      sync.WaitGroup
}

// NewPool creates a worker pool bound to a queue.
func NewPool(q *PriorityQueue, workers int) *Pool {
	if workers < 1 {
		workers = 1
	}
	return &Pool{q: q, workers: workers}
}

// Start launches the workers. They exit once the queue is closed and drained.
func (p *Pool) Start() {
	for i := 1; i <= p.workers; i++ {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			for {
				t, err := p.q.Pop()
				if err != nil {
					return
				}
				if err := t.Run(); err != nil {
					log.Printf("worker %d: task %s failed: %v", id, t.ID, err)
				}
			}
		}(i)
	}
}

// Wait blocks until all workers have finished.
func (p *Pool) Wait() { p.wg.Wait() }
