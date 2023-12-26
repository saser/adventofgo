// Package span implements an integer span of numbers.
package span

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Span is an integer half-open span of numbers. Span encodes a span like [a, b)
// meaning a is contained in the span but b is not.
type Span[T constraints.Integer] struct {
	Start T // Inclusive.
	End   T // Exclusive.
}

// New constructs a new span. If start > end, New returns Empty().
func New[T constraints.Integer](start, end T) Span[T] {
	if start >= end {
		return Empty[T]()
	}
	return Span[T]{Start: start, End: end}
}

// Empty returns the canonical empty span, containing no numbers.
func Empty[T constraints.Integer]() Span[T] {
	return Span[T]{Start: 0, End: 0}
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

// Union returns a new span with all numbers contained in either a or b. If both
// a and b are non-empty and share no numbers, Union returns Empty().
func Union[T constraints.Integer](a, b Span[T]) Span[T] {
	if a.Len() == 0 {
		return b
	}
	if b.Len() == 0 {
		return a
	}
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
	return s.End - s.Start
}

// Contains looks for v in the span [start, end). The return value is:
// * < 0 if v < start (i.e. v comes before the span).
// * == 0 if start <= v < end (i.e. v is covered by the span).
// * > 0 if end <= v (i.e. v comes after the span).
func (s Span[T]) Contains(v T) int {
	if v < s.Start {
		return -1
	}
	if v >= s.End {
		return +1
	}
	return 0
}

func (s Span[T]) String() string {
	return fmt.Sprintf("[%d, %d)", s.Start, s.End)
}

// Split splits s at v into two Spans.
// If s contains v and two non-empty splits can be created, Split returns:
//
//	first  = [s.Start, v)
//	second = [v, s.End)
//
// Otherwise, Split returns:
//
//	first  = s
//	second = Empty()
func (s Span[T]) Split(v T) (first, second Span[T]) {
	if s.Contains(v) != 0 {
		return s, Empty[T]()
	}
	if v == s.Start {
		return s, Empty[T]()
	}
	first = New(s.Start, v)
	second = New(v, s.End)
	return first, second
}

// SplitOver splits s into one or more non-overlapping spans, based on t. The
// split happens as follows:
//
// t contains s entirely:
//
//	s =     aaaaaaaaaa
//	t = bbbbbbbbbbbbbbbbb
//	=>    [ iiiiiiiiii ] (= s)
//
// s contains t entirely:
//
//	s =   aaaaaaaaaaaaaaaaa
//	t =       bbbbbbbbbb
//	=>  [ xxxx,
//	          iiiiiiiiii,
//	                    yyy ]
//
// t contains parts of s:
//
//	s =   aaaaaaaaaa
//	t =       bbbbbbbbbbbbbbbbb
//	=>  [ xxxx,
//	          iiiiii ]
//
//	s =     aaaaaaaaaaaaaaaaa
//	t = bbbbbbbbbb
//	=>    [ iiiiii,
//	              yyyyyyyyyyy ]
//
// t contains no part of s:
//
//	s =       aaaaaaaaaaaaaaaaa
//	t = bbbb
//	=>      [ yyyyyyyyyyyyyyyyy ]
//
//	s =   aaaa
//	t =         bbbbbbbbbbbbbbbbb
//	=>  [ xxxx ]
func (s Span[T]) SplitOver(t Span[T]) []Span[T] {
	i := Intersection(s, t)
	if i == Empty[T]() || i == s {
		return []Span[T]{s}
	}
	spans := make([]Span[T], 0, 3)
	if x := New(s.Start, i.Start); x != Empty[T]() {
		spans = append(spans, x)
	}
	spans = append(spans, i)
	if y := New(i.End, s.End); y != Empty[T]() {
		spans = append(spans, y)
	}
	return spans
}
