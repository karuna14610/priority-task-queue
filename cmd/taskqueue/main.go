// Demo: submit tasks with different priorities and watch the worker pool
// process the most urgent ones first.
package main

import (
	"fmt"
	"time"

	"github.com/YOUR_USERNAME/priority-task-queue/queue"
)

func main() {
	q := queue.New()

	jobs := []struct {
		name     string
		priority int
	}{
		{"send newsletter", 1},
		{"process payment", 10},
		{"resize image", 3},
		{"security alert", 20},
		{"generate report", 2},
		{"password reset email", 8},
	}

	for _, j := range jobs {
		job := j
		_ = q.Push(&queue.Task{
			ID:       job.name,
			Priority: job.priority,
			Run: func() error {
				fmt.Printf("running %-22s (priority %d)\n", job.name, job.priority)
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		})
	}
	q.Close() // no more tasks; workers exit when the queue is drained

	pool := queue.NewPool(q, 1) // 1 worker makes the ordering easy to see
	pool.Start()
	pool.Wait()
	fmt.Println("all tasks done")
}
