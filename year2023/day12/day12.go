package day12

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func solve(input string, part int) (string, error) {
	lines := striter.OverLines(input)
	answer := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		pattern, groupsStr, _ := strings.Cut(line, " ")
		var groups []int
		for _, s := range strings.Split(groupsStr, ",") {
			i, _ := strconv.Atoi(s)
			groups = append(groups, i)
		}
		count := 1
		if part == 2 {
			count = 5
		}
		pattern = strings.Join(slices.Repeat([]string{pattern}, count), "?")
		groups = slices.Repeat(groups, count)
		answer += countPossible(pattern, groups)
	}
	return fmt.Sprint(answer), nil
}

func countPossible(pattern string, groups []int) int {
	// This solution was heavily inspired by
	// https://github.com/ConcurrentCrab/AoC/blob/aba36645b18566bbb7028437a0929d4b6af0e5f5/solutions/12-2.go
	// which I found through Reddit. I rewrote it a little bit to make it easier
	// to understand (for me).

	// Simulate an NFA (Non-deterministic Finite Automaton) that matches the
	// pattern against the groups.
	type state struct {
		GroupsPos  int  // Index into 'groups' of currently matched group.
		Broken     int  // Length of run of consecutive broken springs ('#'s).
		RequireDot bool // Whether the next character *must* be a working spring ('.').
	}
	// An NFA can exist in multiple states at once. These maps, corresponding to
	// the current and the next state respectively, count how many "instances"
	// of the NFA are in each state.
	curr := map[state]int{{0, 0, false}: 1}
	next := map[state]int{}
	for _, c := range pattern {
		chars := []rune{c}
		if c == '?' {
			chars = []rune{'.', '#'}
		}
		for _, c := range chars {
			for s, count := range curr {
				// The cases of this switch statement takes a little while to
				// understand, but what they're essentially doing is modeling
				// the
				switch c {
				case '#':
					if s.RequireDot {
						// We found a '#' but we needed to find a '.', so we've
						// entered an invalid state.
						continue
					}
					if s.GroupsPos == len(groups) {
						// This broken springs cannot be matched against a group
						// because there are no more groups left to match, so
						// we've entered an invalid state.
						continue
					}
					s.Broken++
					if s.Broken == groups[s.GroupsPos] {
						// We've found a sufficient number of consecutive broken
						// springs, so transition to a state where we will match
						// against the next group, reset the count of broken
						// springs, and require that the next character is a
						// dot.
						s.GroupsPos++
						s.Broken = 0
						s.RequireDot = true
					}
					next[state{s.GroupsPos, s.Broken, s.RequireDot}] += count
				case '.':
					if s.Broken != 0 {
						// We were on a run of consecutive broken springs that
						// ended too early, so we've entered an invalid state.
						continue
					}
					next[state{s.GroupsPos, 0, false}] += count
				}
			}
		}
		curr, next = next, curr
		clear(next)
	}

	sum := 0
	for s, count := range curr {
		if s.GroupsPos == len(groups) {
			sum += count
		}
	}
	return sum
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
