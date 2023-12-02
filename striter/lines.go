package striter

// Lines provides iteration over each line in the input. The newline character
// '\n' is used to separate lines. Note that this will likely not work on
// platforms that don't use '\n' as the newline separator.
type Lines struct {
	// Delegate to a Split iterator using "\n" as the separator.
	sp *Split
}

// OverLines constructs a Lines iterator over the input string.
func OverLines(s string) *Lines {
	return &Lines{
		sp: OverSplit(s, "\n"),
	}
}

// Next returns the next line in the input, and reports whether the value is
// valid.
func (i *Lines) Next() (string, bool) {
	return i.sp.Next()
}
