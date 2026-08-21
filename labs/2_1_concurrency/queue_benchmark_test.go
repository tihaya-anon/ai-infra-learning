package queue

import "testing"

func BenchmarkMutexQueueRoundTrip(b *testing.B) {
	var q MutexQueue[int]
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		q.TryDequeue()
	}
}

func BenchmarkSPSCQueueRoundTrip(b *testing.B) {
	q, err := NewSPSCQueue[int](1024)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.TryEnqueue(i)
		q.TryDequeue()
	}
}

func BenchmarkLockFreeQueueRoundTrip(b *testing.B) {
	q := NewLockFreeQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		q.TryDequeue()
	}
}
