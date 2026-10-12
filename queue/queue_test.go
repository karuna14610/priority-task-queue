package queue

import "testing"

func TestPriorityOrder(t *testing.T) {
	q := New()
	for _, tc := range []struct {
		id string
		p  int
	}{{"low", 1}, {"high", 10}, {"mid", 5}} {
		_ = q.Push(&Task{ID: tc.id, Priority: tc.p})
	}
	q.Close()

	want := []string{"high", "mid", "low"}
	for _, w := range want {
		got, err := q.Pop()
		if err != nil || got.ID != w {
			t.Fatalf("want %s, got %v (err %v)", w, got, err)
		}
	}
	if _, err := q.Pop(); err != ErrClosed {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}

func TestFIFOWithinSamePriority(t *testing.T) {
	q := New()
	for _, id := range []string{"a", "b", "c"} {
		_ = q.Push(&Task{ID: id, Priority: 3})
	}
	q.Close()
	for _, w := range []string{"a", "b", "c"} {
		got, _ := q.Pop()
		if got.ID != w {
			t.Fatalf("want %s, got %s", w, got.ID)
		}
	}
}

func TestPushAfterClose(t *testing.T) {
	q := New()
	q.Close()
	if err := q.Push(&Task{ID: "x"}); err != ErrClosed {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}

func TestPoolRunsAllTasks(t *testing.T) {
	q := New()
	done := make(chan string, 10)
	for i := 0; i < 10; i++ {
		id := string(rune('a' + i))
		_ = q.Push(&Task{ID: id, Priority: i, Run: func() error { done <- id; return nil }})
	}
	q.Close()
	p := NewPool(q, 3)
	p.Start()
	p.Wait()
	if len(done) != 10 {
		t.Fatalf("expected 10 tasks run, got %d", len(done))
	}
}
