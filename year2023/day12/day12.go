package day12

import (
	"bytes"
	"errors"
	"fmt"

	"go.saser.se/adventofgo/striter"
)

type record struct {
	Springs []byte
	Groups  []int
}

func parseRecord(s string) *record {
	b := []byte(s)
	r := &record{}
	space := bytes.IndexByte(b, ' ')
	r.Springs = b[:space]

	acc := 0
	for _, c := range b[space+1:] {
		if c == ',' {
			r.Groups = append(r.Groups, acc)
			acc = 0
			continue
		}
		acc = acc*10 + int(c-'0')
	}
	r.Groups = append(r.Groups, acc)

	return r
}

func (r *record) validArrangements(springsIndex, groupsIndex, n int) int {
	foundWildcard := false
	for !foundWildcard {
		if springsIndex == len(r.Springs) {
			// We've reached the end of the input. There are two cases where we
			// have a valid arrangement:
			// 1. n > 0, there's one group left and we've exactly matched it.
			// 2. n == 0 and there are no more groups to be matched.
			if (n > 0 && groupsIndex == len(r.Groups)-1 && n == r.Groups[groupsIndex]) || (n == 0 && groupsIndex == len(r.Groups)) {
				return 1
			}
			// All other cases indicate an invalid arrangement.
			return 0
		}

		c := r.Springs[springsIndex]

		if c == '.' {
			if n > 0 {
				// We've found the end of a group.
				if n != r.Groups[groupsIndex] {
					// The group we just found doesn't match the current target,
					// so this arrangement is invalid.
					return 0
				}
				// The group we just found matches the current target, so reset
				// n and move on to the next group.
				groupsIndex++
				n = 0
			}
			springsIndex++
			continue
		}

		if c == '#' {
			if n == 0 && groupsIndex == len(r.Groups) {
				// We've found the start of a group, but there are no more
				// groups to be matched. Therefore, this arrangement is invalid.
				return 0
			}
			if n > r.Groups[groupsIndex] {
				// We've found a group that is larger than the current target.
				// Therefore, this arrangment is invalid.
				return 0
			}
			n++
			springsIndex++
			continue
		}

		// Implied: c == '?'.
		foundWildcard = true
	}

	// springsIndex now points at a '?' character.

	valid := 0
	// First choice: turn '?' => '#' which means we recurse where we are but
	// with n+1. This is only useful if there are more groups to match _and_
	// we're not going to go over the current group.
	if groupsIndex < len(r.Groups) && n < r.Groups[groupsIndex] {
		valid += r.validArrangements(springsIndex+1, groupsIndex, n+1)
	}
	// Second choice: turn '?' => '.'. If n > 0 then we've "created" a group by
	// terminating it with '.', so check whether we matched the current target.
	// If we did we can recurse by starting with the next group and n = 0.
	if n > 0 && n == r.Groups[groupsIndex] {
		valid += r.validArrangements(springsIndex+1, groupsIndex+1, 0)
	}
	if n == 0 {
		valid += r.validArrangements(springsIndex+1, groupsIndex, 0)
	}
	return valid
}

func (r *record) ValidArrangements() int {
	return r.validArrangements(0, 0, 0)
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	lines := striter.OverLines(input)
	sum := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		sum += parseRecord(line).ValidArrangements()
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
