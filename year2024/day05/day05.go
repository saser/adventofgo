package day05

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"
)

func parsePageOrderings(lines iter.Seq[string]) (map[int][]int, error) {
	m := make(map[int][]int)
	for line := range lines {
		if line == "" {
			break
		}
		beforeStr, afterStr, ok := strings.Cut(line, "|")
		if !ok {
			return nil, fmt.Errorf("invalid line %q", line)
		}
		before, err := strconv.Atoi(beforeStr)
		if err != nil {
			return nil, fmt.Errorf("parse first number: %v", err)
		}
		after, err := strconv.Atoi(afterStr)
		if err != nil {
			return nil, fmt.Errorf("parse second number: %v", err)
		}
		m[before] = append(m[before], after)
	}
	return m, nil
}

func parseUpdates(lines iter.Seq[string]) ([][]int, error) {
	var updates [][]int
	for line := range lines {
		raw := strings.Split(line, ",")
		u := make([]int, len(raw))
		for i, s := range strings.Split(line, ",") {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			u[i] = n
		}
		updates = append(updates, u)
	}
	return updates, nil
}

func solve(input string, part int) (string, error) {
	fragments := strings.Split(input, "\n\n")
	orderings, err := parsePageOrderings(strings.SplitSeq(fragments[0], "\n"))
	if err != nil {
		return "", fmt.Errorf("parse page orderings: %v", err)
	}
	updates, err := parseUpdates(strings.SplitSeq(fragments[1], "\n"))
	if err != nil {
		return "", fmt.Errorf("parse updates: %v", err)
	}
	sortFunc := func(a, b int) int {
		if slices.Contains(orderings[a], b) {
			return -1
		}
		if slices.Contains(orderings[b], a) {
			return +1
		}
		return 0
	}
	sum := 0
	for _, u := range updates {
		valid := slices.IsSortedFunc(u, sortFunc)
		add := false
		switch {
		case valid && part == 1:
			add = true
		case !valid && part == 2:
			slices.SortFunc(u, sortFunc)
			add = true
		}
		if add {
			sum += u[len(u)/2]
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
