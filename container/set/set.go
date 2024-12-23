// Package set provides a set type based on Go's built-in maps.
package set

import (
	"iter"
	"maps"
)

// Set represents an unsorted set of values. It is not thread-safe; use together
// with e.g. a sync.Mutex to ensure thread-safety.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New initializes a new set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// NewSize initializes a new set that has allocted enough memory to hold n
// elements.
func NewSize[T comparable](n int) *Set[T] {
	return &Set[T]{m: make(map[T]struct{}, n)}
}

// Of initializes a set that holds the given elements.
func Of[T comparable](vs ...T) *Set[T] {
	s := NewSize[T](len(vs))
	for _, v := range vs {
		s.Add(v)
	}
	return s
}

// Union returns a new set holding a ⋃ b.
func Union[T comparable](a, b *Set[T]) *Set[T] {
	c := NewSize[T](a.Len() + b.Len())
	for v := range a.All() {
		c.Add(v)
	}
	for v := range b.All() {
		c.Add(v)
	}
	return c
}

// Intersection returns a new set holding a ⋂ b.
func Intersection[T comparable](a, b *Set[T]) *Set[T] {
	c := NewSize[T](max(a.Len(), b.Len()))
	for v := range a.All() {
		if b.Contains(v) {
			c.Add(v)
		}
	}
	return c
}

// Minus returns a new set holding a \ b.
func Minus[T comparable](a, b *Set[T]) *Set[T] {
	c := a.Clone()
	for v := range b.All() {
		c.Delete(v)
	}
	return c
}

// Clone returns a copy of s.
func (s *Set[T]) Clone() *Set[T] {
	return &Set[T]{m: maps.Clone(s.m)}
}

// Add ensures v is contained in s, and returns true if s was modified.
func (s *Set[T]) Add(v T) bool {
	added := !s.Contains(v)
	s.m[v] = struct{}{}
	return added
}

// Delete ensures v is not contained in s, and returns true if s was modified.
func (s *Set[T]) Delete(v T) bool {
	deleted := s.Contains(v)
	delete(s.m, v)
	return deleted
}

// Contains returns true if s contains v.
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

// All iterates over the elements in the set in an undefined, non-deterministic
// order.
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}
