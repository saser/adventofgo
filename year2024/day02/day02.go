package day02

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func isSafeIncreasing(report []int) bool {
	for i := range len(report) - 1 {
		d := report[i+1] - report[i]
		if 1 <= d && d <= 3 {
			continue
		}
		return false
	}
	return true
}

func isSafeIncreasingWithSkip(report []int) bool {
	// Any slice of length 0, 1, or 2 can always be
	// considered safe.
	n := len(report)
	if n <= 2 {
		return true
	}
	for i := range n - 1 {
		x := report[i]
		y := report[i+1]
		if d := y - x; 1 <= d && d <= 3 {
			continue
		}
		// Skip either the current element or the next element. The skip is
		// implemented by swapping elements around in the slice, so that we can
		// always take a subslice and call [isSafeIncreasing] on it. This should
		// (in theory) be more efficient than creating new slices, which will
		// allocate and copy elements on the heap.
		//
		// We have to be careful to restore the order after swapping, since
		// we're mutating the input argument.
		if i == 0 {
			// The slice looks like this:
			//     [x, y, ...]
			// First try skipping x.
			//     initial state: [x, y, ...]
			//     check [1:]:       [y, ...]
			if isSafeIncreasing(report[1:]) {
				return true
			}
			// If that doesn't work, try skipping y:
			//     initial state: [x, y, ...]
			//     swap x and y:  [y, x, ...]
			report[i], report[i+1] = report[i+1], report[i]
			//     check [1:]:       [x, ...]
			ok := isSafeIncreasing(report[1:])
			//     swap x and y:  [x, y, ...]
			report[i], report[i+1] = report[i+1], report[i]
			if ok {
				return true
			}
			// Neither skipping x or y helped; this report is not safe.
			return false
		}
		if i == len(report)-2 {
			// The slice looks like this:
			//     [..., x, y]
			// First try skipping x.
			//     initial state: [..., x, y]
			//     swap x and y:  [..., y, x]
			report[i], report[i+1] = report[i+1], report[i]
			//     check [:n-2]:  [..., y]
			ok := isSafeIncreasing(report[:n-2])
			//     swap x and y:  [..., x, y]
			report[i], report[i+1] = report[i+1], report[i]
			if ok {
				return true
			}
			// If that doesn't work, try skipping y:
			//     initial state: [..., x, y]
			//     check [:n-2]:  [..., x]
			ok = isSafeIncreasing(report[:n-2])
			if ok {
				return true
			}
			// Neither skipping x or y helped; this report is not safe.
			return false
		}
		// The slice looks like this:
		//     [..., w, x, y, z, ...]
		// Since we iterate from the left to the right, we know that the report
		// [..., w, x] is safe. Therefore, we don't need to check anything
		// before w again; we only need to check from w and forward.
		//
		// First try skipping x.
		//                             i
		//     initial state: [..., w, x, y, z, ...]
		//     swap w and x:  [..., x, w, y, z, ...]
		report[i-1], report[i] = report[i], report[i-1]
		//     check from w:          [w, y, z, ...]
		ok := isSafeIncreasing(report[i:])
		//     swap w and x:  [..., w, x, y, z, ...]
		report[i-1], report[i] = report[i], report[i-1]
		if ok {
			return true
		}
		// If that doesn't work, try skipping y.
		//                             i
		//     initial state: [..., w, x, y, z, ...]
		//     swap x and y:  [..., w, y, x, z, ...]
		report[i], report[i+1] = report[i+1], report[i]
		// check from x:                 [x, z, ...]
		ok = isSafeIncreasing(report[i+1:])
		//     swap x and y:  [..., w, x, y, z, ...]
		report[i], report[i+1] = report[i+1], report[i]
		if ok {
			return true
		}
		// Neither skipping x or y helped; this report is not safe.
		return false
	}
	// The report is safe without skipping anything.
	return true
}

func solve(input string, part int) (string, error) {
	lines := striter.OverLines(input)
	count := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		fields := strings.Fields(line)
		report := make([]int, len(fields))
		for i, f := range fields {
			var err error
			report[i], err = strconv.Atoi(f)
			if err != nil {
				return "", fmt.Errorf("parse integer from line %q: %w", line, err)
			}
		}
		var isSafe func(report []int) bool
		if part == 1 {
			isSafe = isSafeIncreasing
		} else {
			isSafe = isSafeIncreasingWithSkip
		}
		// Try both directions, if necessary.
		if isSafe(report) {
			count++
			continue
		}
		slices.Reverse(report)
		if isSafe(report) {
			count++
			continue
		}
	}
	return fmt.Sprint(count), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
