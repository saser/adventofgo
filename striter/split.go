package striter

// TODO: add tests.

import "strings"

// Split is an iterator analogous to strings.Split from the standard library.
// However, it might not handle all the edge cases around UTF-8 encoded strings
// that strings.Split does, as Split is intended to be used in an Advent of Code
// context, where all strings are assumed to be ASCII only.
type Split struct {
	// s is the whole input string
	s string
	// sep is the separator.
	sep string
	// pos points to the current position within s.
	pos int
}

// OverSplit creates an iterator over all chunks of s separated by sep.
func OverSplit(s string, sep string) *Split {
	return &Split{
		s:   s,
		sep: sep,
		pos: 0,
	}
}

// Next returns the next chunk in the input string, and reports whether the
// value is valid.
func (i *Split) Next() (string, bool) {
	// Note that this is > rather than >=. This is because i.pos is set to
	// len(i.s)+1 upon reaching a point where there are no remaining separators.
	// This also covers the cases of empty strings, trailing separators, etc.
	if i.pos > len(i.s) {
		return "", false
	}
	start := i.pos
	delta := strings.Index(i.s[start:], i.sep)
	end := start + delta
	if delta == -1 {
		end = len(i.s)
	}
	// If there is a next separator, then i.pos will skip over it.
	// If not, end will point to len(i.s)+1.
	i.pos = end + len(i.sep) // Skip over the next separator.
	return i.s[start:end], true
}
