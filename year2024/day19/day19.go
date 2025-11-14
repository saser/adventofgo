package day19

import (
	"fmt"
	"strings"
)

// arrangementCount returns the number of possible arrangements for the given
// towel.
func arrangementCount(towel string, patterns []string, memo map[string]uint64) (n uint64) {
	if prev, ok := memo[towel]; ok {
		return prev
	}
	defer func() { memo[towel] = n }()
	if towel == "" {
		return 1
	}
	var sum uint64
	for _, p := range patterns {
		if strings.HasPrefix(towel, p) {
			sum += arrangementCount(strings.TrimPrefix(towel, p), patterns, memo)
		}
	}
	return sum
}

// isPossible returns 1 if the given towel is possible to create, and 0
// otherwise.
func isPossible(towel string, patterns []string, memo map[string]uint64) (n uint64) {
	if prev, ok := memo[towel]; ok {
		return prev
	}
	defer func() { memo[towel] = n }()
	if towel == "" {
		return 1
	}
	for _, p := range patterns {
		if strings.HasPrefix(towel, p) && isPossible(strings.TrimPrefix(towel, p), patterns, memo) == 1 {
			return 1
		}
	}
	return 0
}

func solve(input string, part int) (string, error) {
	i := 0
	var patterns []string
	memo := make(map[string]uint64)
	var answer uint64
	countFn := isPossible
	if part == 2 {
		countFn = arrangementCount
	}
	for line := range strings.SplitSeq(input, "\n") {
		switch i {
		case 0:
			patterns = strings.Split(line, ", ")
		case 1:
			// Empty line between patterns and towels.
		default:
			answer += countFn(line, patterns, memo)
		}
		i++
	}
	return fmt.Sprint(answer), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
