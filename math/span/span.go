// Package span implements an integer span of numbers.
package span

import "golang.org/x/exp/constraints"

// Span is an integer span of numbers.
type Span[T constraints.Integer] struct {
	Start, End T // Inclusive of both.

	isEmpty bool
}

// New constructs a new span. If start > end, New returns Empty().
func New[T constraints.Integer](start, end T) Span[T] {
	if start > end {
		return Empty[T]()
	}
	return Span[T]{Start: start, End: end}
}

// Empty returns the canonical empty span, containing no numbers.
func Empty[T constraints.Integer]() Span[T] {
	return Span[T]{Start: 0, End: 0, isEmpty: true}
}

// Intersection returns a new span with all numbers contained in both a and b.
// If there are no such numbers, Intersection returns Empty().
func Intersection[T constraints.Integer](a, b Span[T]) Span[T] {
	s := Span[T]{
		Start: max(a.Start, b.Start),
		End:   min(a.End, b.End),
	}
	if s.Start >= s.End {
		return Empty[T]()
	}
	return s
}

// Union returns a new span with all numbers contained in either a or b,
// with the condition that a and b share at least one number. If not, Union
// returns Empty().
func Union[T constraints.Integer](a, b Span[T]) Span[T] {
	if Intersection(a, b) == Empty[T]() {
		return Empty[T]()
	}
	return Span[T]{
		Start: min(a.Start, b.Start),
		End:   max(a.End, b.End),
	}
}

// Len returns the number of integers covered by this span.
func (s Span[T]) Len() T {
	if s == Empty[T]() {
		return 0
	}
	return s.End - s.Start + 1
}
