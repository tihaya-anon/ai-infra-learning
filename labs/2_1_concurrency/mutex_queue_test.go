package queue

import "testing"

func TestMutexQueueFIFO(t *testing.T) {
	var q MutexQueue[int]
	q.Enqueue(10)
	q.Enqueue(20)

	assertDequeued(t, q.TryDequeue, 10)
	assertDequeued(t, q.TryDequeue, 20)
	assertEmpty(t, q.TryDequeue)

	if got := q.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}
