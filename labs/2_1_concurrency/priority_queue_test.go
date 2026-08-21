package queue

import (
	"errors"
	"testing"
	"time"
)

func TestBlockingPriorityQueueOrdersByPriorityThenFIFO(t *testing.T) {
	q := NewBlockingPriorityQueue[string]()
	for _, item := range []struct {
		value    string
		priority int
	}{
		{value: "low", priority: 1},
		{value: "high-first", priority: 10},
		{value: "medium", priority: 5},
		{value: "high-second", priority: 10},
	} {
		if err := q.Push(item.value, item.priority); err != nil {
			t.Fatal(err)
		}
	}

	for _, want := range []string{"high-first", "high-second", "medium", "low"} {
		value, err := q.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if value != want {
			t.Fatalf("Pop() = %q, want %q", value, want)
		}
	}
}

func TestBlockingPriorityQueueWakesWaitingConsumer(t *testing.T) {
	q := NewBlockingPriorityQueue[int]()
	result := make(chan int, 1)
	go func() {
		value, err := q.Pop()
		if err == nil {
			result <- value
		}
	}()

	if err := q.Push(42, 1); err != nil {
		t.Fatal(err)
	}

	select {
	case value := <-result:
		if value != 42 {
			t.Fatalf("Pop() = %d, want 42", value)
		}
	case <-time.After(time.Second):
		t.Fatal("waiting consumer was not woken")
	}
}

func TestBlockingPriorityQueueDrainsAfterClose(t *testing.T) {
	q := NewBlockingPriorityQueue[string]()
	if err := q.Push("queued", 1); err != nil {
		t.Fatal(err)
	}
	q.Close()
	q.Close()

	value, err := q.Pop()
	if err != nil || value != "queued" {
		t.Fatalf("Pop() = (%q, %v), want (queued, nil)", value, err)
	}
	if _, err := q.Pop(); !errors.Is(err, ErrClosed) {
		t.Fatalf("Pop() error = %v, want ErrClosed", err)
	}
	if err := q.Push("late", 1); !errors.Is(err, ErrClosed) {
		t.Fatalf("Push() error = %v, want ErrClosed", err)
	}
}
