// Package span implements an integer span of numbers.
package span

// Span is an integer span of numbers.
type Span struct {
	Start, End int // Inclusive of both.
}

// Empty returns the canonical empty span, containing no numbers.
func Empty() Span {
	return Span{Start: 0, End: -1}
}

// Intersection returns a new span with all numbers contained in both a and b.
// If there are no such numbers, Intersection returns Empty().
func Intersection(a, b Span) Span {
	s := Span{
		Start: max(a.Start, b.Start),
		End:   min(a.End, b.End),
	}
	if s.Start >= s.End {
		return Empty()
	}
	return s
}

// Union returns a new span with all numbers contained in either a or b,
// with the condition that a and b share at least one number. If not, Union
// returns Empty().
func Union(a, b Span) Span {
	if Intersection(a, b) == Empty() {
		return Empty()
	}
	return Span{
		Start: min(a.Start, b.Start),
		End:   max(a.End, b.End),
	}
}

// Len returns the number of integers covered by this span.
func (s Span) Len() int {
	if s == Empty() {
		return 0
	}
	return s.End - s.Start + 1
}
