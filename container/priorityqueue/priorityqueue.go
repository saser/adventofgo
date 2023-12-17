// Package priorityqueue implements a priority queue based on a heap implemented
// on a slice.
package priorityqueue

import (
	"cmp"
	"container/heap"
)

type queue[T any] struct {
	q    []T
	less func(x, y T) bool
}

var _ heap.Interface = (*queue[string])(nil)

func (q *queue[T]) Len() int           { return len(q.q) }
func (q *queue[T]) Less(i, j int) bool { return q.less(q.q[i], q.q[j]) }
func (q *queue[T]) Swap(i, j int)      { q.q[i], q.q[j] = q.q[j], q.q[i] }

func (q *queue[T]) Push(x any) {
	q.q = append(q.q, x.(T))
}

func (q *queue[T]) Pop() any {
	n := len(q.q)
	x := q.q[n-1]
	q.q = q.q[:n-1]
	return x
}

// Queue is a priority queue containing elements of type T.
type Queue[T any] struct {
	q *queue[T]
}

// New creates a new min-queue for ordered types.
func New[T cmp.Ordered]() *Queue[T] {
	return NewFunc(cmp.Less[T])
}

// NewFunc creates a new min-queue for any type using the provided comparison
// function.
func NewFunc[T any](less func(x, y T) bool) *Queue[T] {
	return &Queue[T]{
		q: &queue[T]{less: less},
	}
}

// Len is the number of elements currently in the queue.
func (q *Queue[T]) Len() int { return q.q.Len() }

// Push adds a new element to the queue.
func (q *Queue[T]) Push(x T) { heap.Push(q.q, x) }

// Pop removes and returns the minimum element according to the ordering this
// queue was created with.
func (q *Queue[T]) Pop() T { return heap.Pop(q.q).(T) }
