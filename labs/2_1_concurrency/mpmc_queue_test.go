package queue

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLockFreeQueueFIFO(t *testing.T) {
	q := NewLockFreeQueue[int]()
	if !q.IsEmpty() {
		t.Fatal("new queue is not empty")
	}

	q.Enqueue(10)
	q.Enqueue(20)
	assertDequeued(t, q.TryDequeue, 10)
	assertDequeued(t, q.TryDequeue, 20)
	assertEmpty(t, q.TryDequeue)

	if !q.IsEmpty() {
		t.Fatal("drained queue is not empty")
	}
}

func TestLockFreeQueueMultipleProducersAndConsumers(t *testing.T) {
	const (
		producerCount  = 4
		consumerCount  = 4
		itemsPerWriter = 25_000
		totalItems     = producerCount * itemsPerWriter
	)

	q := NewLockFreeQueue[int]()
	seen := make([]atomic.Uint32, totalItems)
	var consumed atomic.Int64
	var failed atomic.Bool
	failure := make(chan error, 1)
	recordFailure := func(err error) {
		if failed.CompareAndSwap(false, true) {
			failure <- err
		}
	}

	var producers sync.WaitGroup
	producers.Add(producerCount)
	for producer := 0; producer < producerCount; producer++ {
		go func(producer int) {
			defer producers.Done()
			start := producer * itemsPerWriter
			for offset := 0; offset < itemsPerWriter; offset++ {
				q.Enqueue(start + offset)
			}
		}(producer)
	}

	var consumers sync.WaitGroup
	consumers.Add(consumerCount)
	for range consumerCount {
		go func() {
			defer consumers.Done()
			for consumed.Load() < totalItems && !failed.Load() {
				value, ok := q.TryDequeue()
				if !ok {
					runtime.Gosched()
					continue
				}
				if value < 0 || value >= totalItems {
					recordFailure(fmt.Errorf("dequeued out-of-range value %d", value))
					return
				}
				if count := seen[value].Add(1); count != 1 {
					recordFailure(fmt.Errorf("dequeued value %d %d times", value, count))
					return
				}
				consumed.Add(1)
			}
		}()
	}

	finished := make(chan struct{})
	go func() {
		producers.Wait()
		consumers.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-failure:
		t.Fatal(err)
	case <-time.After(10 * time.Second):
		t.Fatalf("timed out after consuming %d/%d values", consumed.Load(), totalItems)
	}

	if got := consumed.Load(); got != totalItems {
		t.Fatalf("consumed %d values, want %d", got, totalItems)
	}
	for value := range totalItems {
		if got := seen[value].Load(); got != 1 {
			t.Fatalf("value %d was seen %d times", value, got)
		}
	}
}
