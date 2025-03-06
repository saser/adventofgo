package day02

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func isSafeDecreasing(report []int, skipped bool) bool {
	for i := range len(report) - 1 {
		d := report[i] - report[i+1]
		if 1 <= d && d <= 3 {
			continue
		}
		if !skipped {
			// Skip either the current and the next element, and then check
			// again, setting the skipped flag to true.
			var beforeCurrent []int
			if i > 0 {
				beforeCurrent = report[i-1 : i]
			}
			afterCurrent := report[i+1:]
			if isSafeDecreasing(slices.Concat(beforeCurrent, afterCurrent), true) {
				return true
			}
			beforeNext := report[i : i+1]
			var afterNext []int
			if i < len(report)-2 {
				afterNext = report[i+2:]
			}
			return isSafeDecreasing(slices.Concat(beforeNext, afterNext), true)
		}
		return false
	}
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
		skipped := true
		if part == 2 {
			skipped = false
		}
		// Try both directions, if necessary.
		if isSafeDecreasing(report, skipped) {
			count++
			continue
		}
		slices.Reverse(report)
		if isSafeDecreasing(report, skipped) {
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
