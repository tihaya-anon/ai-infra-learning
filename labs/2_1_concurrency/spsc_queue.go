package queue

import (
	"errors"
	"sync/atomic"
)

const cacheLineSize = 64

var ErrInvalidCapacity = errors.New("capacity must be a positive power of two")

// SPSCQueue is a bounded lock-free ring buffer for exactly one producer and
// one consumer. Calling either side from multiple goroutines is unsupported.
// A queue must not be copied after construction.
type SPSCQueue[T any] struct {
	slots []T
	mask  uint64

	_        [cacheLineSize]byte
	readPos  atomic.Uint64
	_        [cacheLineSize]byte
	writePos atomic.Uint64
	_        [cacheLineSize]byte
}

func NewSPSCQueue[T any](capacity int) (*SPSCQueue[T], error) {
	if capacity <= 0 || capacity&(capacity-1) != 0 {
		return nil, ErrInvalidCapacity
	}

	return &SPSCQueue[T]{
		slots: make([]T, capacity),
		mask:  uint64(capacity - 1),
	}, nil
}

func (q *SPSCQueue[T]) TryEnqueue(value T) bool {
	writePos := q.writePos.Load()
	if writePos-q.readPos.Load() == uint64(len(q.slots)) {
		return false
	}

	q.slots[writePos&q.mask] = value
	q.writePos.Store(writePos + 1)
	return true
}

func (q *SPSCQueue[T]) TryDequeue() (T, bool) {
	readPos := q.readPos.Load()
	if readPos == q.writePos.Load() {
		var zero T
		return zero, false
	}

	index := readPos & q.mask
	value := q.slots[index]
	var zero T
	q.slots[index] = zero
	q.readPos.Store(readPos + 1)
	return value, true
}

func (q *SPSCQueue[T]) Capacity() int {
	return len(q.slots)
}

// Len is a moment-in-time observation and may be stale immediately after it
// returns when the producer and consumer are active.
func (q *SPSCQueue[T]) Len() int {
	readPos := q.readPos.Load()
	return int(q.writePos.Load() - readPos)
}
