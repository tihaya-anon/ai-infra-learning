package queue

import "sync/atomic"

// LockFreeQueue is an unbounded Michael-Scott FIFO queue. Multiple producers
// and consumers may call it concurrently. Go's garbage collector keeps nodes
// alive while another goroutine can still observe them. A queue must be created
// with NewLockFreeQueue and must not be copied.
type LockFreeQueue[T any] struct {
	head atomic.Pointer[linkedNode[T]]
	tail atomic.Pointer[linkedNode[T]]
}

type linkedNode[T any] struct {
	value T
	next  atomic.Pointer[linkedNode[T]]
}

func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	dummy := &linkedNode[T]{}
	q := &LockFreeQueue[T]{}
	q.head.Store(dummy)
	q.tail.Store(dummy)
	return q
}

func (q *LockFreeQueue[T]) Enqueue(value T) {
	newNode := &linkedNode[T]{value: value}

	for {
		tail := q.tail.Load()
		next := tail.next.Load()
		if tail != q.tail.Load() {
			continue
		}

		if next != nil {
			q.tail.CompareAndSwap(tail, next)
			continue
		}

		if tail.next.CompareAndSwap(nil, newNode) {
			q.tail.CompareAndSwap(tail, newNode)
			return
		}
	}
}

func (q *LockFreeQueue[T]) TryDequeue() (T, bool) {
	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.next.Load()
		if head != q.head.Load() {
			continue
		}

		if next == nil {
			var zero T
			return zero, false
		}

		if head == tail {
			q.tail.CompareAndSwap(tail, next)
			continue
		}

		value := next.value
		if q.head.CompareAndSwap(head, next) {
			return value, true
		}
	}
}

// IsEmpty is only an observation. Another goroutine may enqueue immediately
// after it returns.
func (q *LockFreeQueue[T]) IsEmpty() bool {
	head := q.head.Load()
	return head.next.Load() == nil
}
