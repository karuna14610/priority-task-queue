# Priority Task Queue (Go)

A small, thread-safe priority queue with a worker pool, written in Go.
Urgent tasks are processed first; tasks with equal priority run in the order they were submitted.

## Features
- Max-heap priority queue built on `container/heap`
- FIFO ordering within the same priority
- Safe for concurrent use (`sync.Mutex` + `sync.Cond`)
- Worker pool that drains the queue and shuts down cleanly
- Unit tests, including a race-detector run

## Run
```bash
go run ./cmd/taskqueue
go test ./... -race
```

Example output (one worker, so the order is visible):
```
running security alert         (priority 20)
running process payment        (priority 10)
running password reset email   (priority 8)
...
```

## Project layout
```
cmd/taskqueue/   demo program
queue/           queue.go (heap + locking), worker.go (pool), tests
```

## Ideas for next steps
- Retries with exponential backoff
- Task timeouts using `context.Context`
- Priority aging so low-priority tasks never starve
- HTTP API to submit tasks
- Integrate with a workflow engine such as Cadence or Temporal
