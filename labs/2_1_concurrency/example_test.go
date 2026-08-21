package queue_test

import (
	"fmt"

	queue "github.com/tihaya-anon/ai-infra-learning/labs/2_1_concurrency"
)

func ExampleLockFreeQueue() {
	q := queue.NewLockFreeQueue[string]()
	q.Enqueue("first")
	q.Enqueue("second")

	for value, ok := q.TryDequeue(); ok; value, ok = q.TryDequeue() {
		fmt.Println(value)
	}

	// Output:
	// first
	// second
}

func ExampleBlockingPriorityQueue() {
	q := queue.NewBlockingPriorityQueue[string]()
	_ = q.Push("background", 1)
	_ = q.Push("urgent", 10)
	q.Close()

	first, _ := q.Pop()
	second, _ := q.Pop()
	fmt.Println(first)
	fmt.Println(second)

	// Output:
	// urgent
	// background
}
