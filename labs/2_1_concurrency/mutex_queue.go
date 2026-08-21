package queue

import "sync"

// MutexQueue is an unbounded FIFO queue protected by a mutex.
// It is the baseline for comparing the lock-free implementations in this lab.
// A queue must not be copied after first use.
type MutexQueue[T any] struct {
	mu     sync.Mutex
	values []T
}

func (q *MutexQueue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.values = append(q.values, value)
}

func (q *MutexQueue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.values) == 0 {
		var zero T
		return zero, false
	}

	value := q.values[0]
	var zero T
	q.values[0] = zero
	if len(q.values) == 1 {
		q.values = q.values[:0]
	} else {
		q.values = q.values[1:]
	}

	return value, true
}

func (q *MutexQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.values)
}
