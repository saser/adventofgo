// Package striter contains an iterator-like interface for strings intended as
// efficient alternatives to splitting functions in package strings.
//
// Many Advent of Code puzzles requires you to iterate over lines (or some other
// segment) of the input and perform some calculation on it. In many cases the
// calculation of each line can be performed completely independent over the
// other lines. This means that functions such as strings.Split, that allocate a
// new (potentially large) slice for each call, may perform worse than iteration
// over substrings. This package is intended to provide convenient ways to
// implement the latter.
package striter

// Iter represents an iterator over a string. It is intended to be very similar
// to the propose Iter[E any] interface from
// https://github.com/golang/go/discussions/54245.
type Iter interface {
	// Next returns the next string in the iteration if there is one, and
	// reports whether the returned value is valid. Once Next returns ok==false,
	// the iteration is over, and all subsequent calls will return ok==false.
	Next() (s string, ok bool)
}

// Collect consumes all remaining values in the iterator and collects them in a
// slice. Any subsequent calls to iter.Next() will return ok==false.
func Collect(iter Iter) []string {
	var ss []string
	for s, ok := iter.Next(); ok; s, ok = iter.Next() {
		ss = append(ss, s)
	}
	return ss
}
