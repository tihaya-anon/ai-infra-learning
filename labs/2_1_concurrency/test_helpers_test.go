package queue

import "testing"

func assertDequeued[T comparable](t *testing.T, dequeue func() (T, bool), want T) {
	t.Helper()

	got, ok := dequeue()
	if !ok {
		t.Fatalf("TryDequeue() reported an empty queue, want %v", want)
	}
	if got != want {
		t.Fatalf("TryDequeue() = %v, want %v", got, want)
	}
}

func assertEmpty[T any](t *testing.T, dequeue func() (T, bool)) {
	t.Helper()

	if value, ok := dequeue(); ok {
		t.Fatalf("TryDequeue() = %v, want empty queue", value)
	}
}
