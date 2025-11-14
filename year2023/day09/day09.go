package day09

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type report []int

func parseReport(line string) (report, error) {
	var r report
	for s := range strings.SplitSeq(line, " ") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("parse report from %q: parse number %q: %v", line, s, err)
		}
		r = append(r, i)
	}
	return r, nil
}

func nextValue(r report) int {
	next := 0
	for {
		n := len(r) - 1
		next += r[n]
		onlyZeroes := true
		for i := 0; i < n; i++ {
			d := r[i+1] - r[i]
			r[i] = d
			onlyZeroes = onlyZeroes && d == 0
		}
		if onlyZeroes {
			break
		}
		r = r[:n]
	}
	return next
}

func parse(input string) ([]report, error) {
	var rs []report
	for line := range strings.SplitSeq(input, "\n") {
		r, err := parseReport(line)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func solve(input string, part int) (string, error) {
	rs, err := parse(input)
	if err != nil {
		return "", err
	}
	sum := 0
	for _, r := range rs {
		if part == 2 {
			slices.Reverse(r)
		}
		sum += nextValue(r)
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
