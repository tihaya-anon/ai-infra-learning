package queue

import (
	"container/heap"
	"errors"
	"sync"
)

var ErrClosed = errors.New("queue is closed")

// BlockingPriorityQueue returns higher priorities first and preserves FIFO
// order between values with equal priority. It uses a mutex and is not lock-free.
// A queue must be created with NewBlockingPriorityQueue and must not be copied.
type BlockingPriorityQueue[T any] struct {
	mu       sync.Mutex
	ready    *sync.Cond
	items    priorityHeap[T]
	nextSeq  uint64
	isClosed bool
}

type priorityItem[T any] struct {
	value    T
	priority int
	sequence uint64
}

type priorityHeap[T any] []priorityItem[T]

func NewBlockingPriorityQueue[T any]() *BlockingPriorityQueue[T] {
	q := &BlockingPriorityQueue[T]{}
	q.ready = sync.NewCond(&q.mu)
	return q
}

func (q *BlockingPriorityQueue[T]) Push(value T, priority int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isClosed {
		return ErrClosed
	}

	heap.Push(&q.items, priorityItem[T]{
		value:    value,
		priority: priority,
		sequence: q.nextSeq,
	})
	q.nextSeq++
	q.ready.Signal()
	return nil
}

// Pop blocks until an item is available or the closed queue has been drained.
func (q *BlockingPriorityQueue[T]) Pop() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.isClosed {
		q.ready.Wait()
	}
	if len(q.items) == 0 {
		var zero T
		return zero, ErrClosed
	}

	item := heap.Pop(&q.items).(priorityItem[T])
	return item.value, nil
}

func (q *BlockingPriorityQueue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isClosed {
		return
	}
	q.isClosed = true
	q.ready.Broadcast()
}

func (q *BlockingPriorityQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.items)
}

func (h priorityHeap[T]) Len() int { return len(h) }

func (h priorityHeap[T]) Less(i, j int) bool {
	if h[i].priority == h[j].priority {
		return h[i].sequence < h[j].sequence
	}
	return h[i].priority > h[j].priority
}

func (h priorityHeap[T]) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *priorityHeap[T]) Push(value any) {
	*h = append(*h, value.(priorityItem[T]))
}

func (h *priorityHeap[T]) Pop() any {
	old := *h
	last := len(old) - 1
	item := old[last]
	var zero priorityItem[T]
	old[last] = zero
	*h = old[:last]
	return item
}
