// Package set provides a set type based on Go's built-in maps.
package set

import (
	"iter"
	"maps"
)

// Set represents an unsorted set of values. It is not thread-safe; use together
// with e.g. a sync.Mutex to ensure thread-safety.
type Set[T comparable] map[T]struct{}

// Of initializes a set that holds the given elements.
func Of[T comparable](vs ...T) Set[T] {
	s := make(Set[T], (len(vs)))
	for _, v := range vs {
		s.Add(v)
	}
	return s
}

// Union returns a new set holding a ⋃ b.
func Union[T comparable](a, b Set[T]) Set[T] {
	c := make(Set[T], a.Len()+b.Len())
	for v := range a.All() {
		c.Add(v)
	}
	for v := range b.All() {
		c.Add(v)
	}
	return c
}

// Intersection returns a new set holding a ⋂ b.
func Intersection[T comparable](a, b Set[T]) Set[T] {
	c := make(Set[T], max(a.Len(), b.Len()))
	for v := range a.All() {
		if b.Contains(v) {
			c.Add(v)
		}
	}
	return c
}

// Minus returns a new set holding a \ b.
func Minus[T comparable](a, b Set[T]) Set[T] {
	c := a.Clone()
	for v := range b.All() {
		c.Delete(v)
	}
	return c
}

// Clone returns a copy of s.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// Add ensures v is contained in s, and returns true if s was modified.
func (s Set[T]) Add(v T) bool {
	added := !s.Contains(v)
	s[v] = struct{}{}
	return added
}

// Delete ensures v is not contained in s, and returns true if s was modified.
func (s Set[T]) Delete(v T) bool {
	deleted := s.Contains(v)
	delete(s, v)
	return deleted
}

// Contains returns true if s contains v.
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// All iterates over the elements in the set in an undefined, non-deterministic
// order.
func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s)
}
