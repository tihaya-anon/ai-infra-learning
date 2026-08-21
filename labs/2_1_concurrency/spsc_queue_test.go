package queue

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNewSPSCQueueRejectsInvalidCapacity(t *testing.T) {
	for _, capacity := range []int{-1, 0, 3, 6} {
		_, err := NewSPSCQueue[int](capacity)
		if !errors.Is(err, ErrInvalidCapacity) {
			t.Errorf("NewSPSCQueue(%d) error = %v, want ErrInvalidCapacity", capacity, err)
		}
	}
}

func TestSPSCQueueFullEmptyAndWrapAround(t *testing.T) {
	q, err := NewSPSCQueue[int](4)
	if err != nil {
		t.Fatal(err)
	}

	for value := 0; value < q.Capacity(); value++ {
		if !q.TryEnqueue(value) {
			t.Fatalf("TryEnqueue(%d) failed before queue was full", value)
		}
	}
	if q.TryEnqueue(4) {
		t.Fatal("TryEnqueue succeeded on a full queue")
	}

	assertDequeued(t, q.TryDequeue, 0)
	assertDequeued(t, q.TryDequeue, 1)
	if !q.TryEnqueue(4) || !q.TryEnqueue(5) {
		t.Fatal("TryEnqueue failed after space became available")
	}

	for _, want := range []int{2, 3, 4, 5} {
		assertDequeued(t, q.TryDequeue, want)
	}
	assertEmpty(t, q.TryDequeue)
}

func TestSPSCQueueConcurrentProducerAndConsumer(t *testing.T) {
	q, err := NewSPSCQueue[int](256)
	if err != nil {
		t.Fatal(err)
	}

	const total = 100_000
	consumerResult := make(chan error, 1)
	go func() {
		for want := 0; want < total; {
			value, ok := q.TryDequeue()
			if !ok {
				runtime.Gosched()
				continue
			}
			if value != want {
				consumerResult <- fmt.Errorf("dequeued %d, want %d", value, want)
				return
			}
			want++
		}
		consumerResult <- nil
	}()

	for value := 0; value < total; {
		if q.TryEnqueue(value) {
			value++
			continue
		}
		runtime.Gosched()
	}

	select {
	case err := <-consumerResult:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("consumer did not finish")
	}
}
