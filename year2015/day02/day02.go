package day02

import (
	"fmt"
	"strings"
)

func min(i int, is ...int) int {
	r := i
	for _, ii := range is {
		if ii < r {
			r = ii
		}
	}
	return r
}

func wrappingPaper(l, w, h int) int {
	return 2*l*w + 2*w*h + 2*h*l + min(l*w, w*h, l*h)
}

func ribbon(l, w, h int) int {
	return l*w*h + min(2*(l+w), 2*(w+h), 2*(l+h))
}

func solve(input string, part int) (string, error) {
	sum := 0
	for line := range strings.SplitSeq(input, "\n") {
		var l, w, h int
		if _, err := fmt.Sscanf(line, "%dx%dx%d", &l, &w, &h); err != nil {
			return "", fmt.Errorf("parse line %q: %v", line, err)
		}
		if part == 1 {
			sum += wrappingPaper(l, w, h)
		} else {
			sum += ribbon(l, w, h)
		}
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
