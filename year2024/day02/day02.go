package day02

import (
	"fmt"
	"iter"
	"slices"
)

func parse(input string) iter.Seq[[]int] {
	// We can make assumptions to efficiently parse a line:
	// 1. Each line looks like "x y z ..." where x, y, z are base 10 integers
	//    with no signs, separators, etc.
	// 2. Estimate how many integers each line contains, and allocating memory
	//    up front for them.
	// In my input, the frequency of report lengths looked like this:
	//     length   frequency
	//     8        250
	//     7        264
	//     6        238
	//     5        248
	// So we simply allocate space for 8 elements and hope that no line contains
	// more integers, so that we never have to grow the slice.
	//
	// Furthermore, to avoid allocating a new slice for every report, we use an
	// iterator and have it retain ownership over the memory backing the slice
	// that is yielded to the caller. This allows us to only ever need to
	// allocate once for all reports. If a caller wants to retain access to a
	// given report, they will have to clone the yielded slice.
	return func(yield func([]int) bool) {
		n := 0
		report := make([]int, 0, 8)
		for _, r := range input {
			if r == '\n' {
				report = append(report, n)
				n = 0
				if !yield(report) {
					return
				}
				report = report[:0]
				continue
			}
			if r == ' ' {
				report = append(report, n)
				n = 0
				continue
			}
			n = n*10 + int(r-'0')
		}
		report = append(report, n)
		yield(report)
	}
}

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
			// We iterate from left to right and therefore know that this is
			// safe:
			//     [..., x]
			// So, we can just skip y.
			return true
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
	count := 0
	for report := range parse(input) {
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
