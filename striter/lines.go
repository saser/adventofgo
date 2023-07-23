package striter

import "strings"

// Lines provides iteration over each line in the input. The newline character
// '\n' is used to separate lines. Note that this will likely not work on
// platforms that don't use '\n' as the newline separator.
type Lines struct {
	// s is the whole input string.
	s string
	// pos points to the current position within s.
	pos int
}

// OverLines constructs a Lines iterator over the input string.
func OverLines(s string) *Lines {
	return &Lines{
		s:   s,
		pos: 0,
	}
}

// Next returns the next line in the input, and reports whether the value is
// valid.
func (i *Lines) Next() (string, bool) {
	// Note that this is > rather than >=. This is because i.pos is set to
	// len(i.s)+1 upon reaching a point where there are no remaining newlines.
	// This also covers the cases of empty strings, trailing newlines, etc.
	if i.pos > len(i.s) {
		return "", false
	}
	start := i.pos
	delta := strings.IndexRune(i.s[start:], '\n')
	end := start + delta
	if delta == -1 {
		end = len(i.s)
	}
	// If there is a next newline, then i.pos will skip over it.
	// If not, end will point to len(i.s)+1.
	i.pos = end + 1 // Skip over the next newline.
	return i.s[start:end], true
}
