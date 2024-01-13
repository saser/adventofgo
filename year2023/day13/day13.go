package day13

import (
	"fmt"

	"go.saser.se/adventofgo/striter"
)

type pattern [][]bool

func parsePattern(fragment string) pattern {
	var g pattern
	lines := striter.OverLines(fragment)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		row := make([]bool, len(line))
		for i, c := range line {
			row[i] = c == '#'
		}
		g = append(g, row)
	}
	return g
}

func parse(input string) []pattern {
	var ps []pattern
	fragments := striter.OverSplit(input, "\n\n")
	for fragment, ok := fragments.Next(); ok; fragment, ok = fragments.Next() {
		ps = append(ps, parsePattern(fragment))
	}
	return ps
}

func mirroredOverRow(g pattern, row int, fixSmudge bool) bool {
	hasFixedSmudge := false
	for j, k := row, row+1; j >= 0 && k < len(g); j, k = j-1, k+1 {
		lo, hi := g[j], g[k]
		diff := 0
		for col := range lo {
			if lo[col] != hi[col] {
				diff++
				if !fixSmudge || diff > 1 {
					return false
				}
			}
		}
		if diff == 0 {
			continue
		}
		if diff == 1 && !hasFixedSmudge {
			hasFixedSmudge = true
			continue
		}
		return false
	}
	if fixSmudge {
		return hasFixedSmudge
	}
	return true
}

func mirroredOverColumn(g pattern, col int, fixSmudge bool) bool {
	hasFixedSmudge := false
	for j, k := col, col+1; j >= 0 && k < len(g[0]); j, k = j-1, k+1 {
		diff := 0
		for row := range g {
			if g[row][j] != g[row][k] {
				diff++
				if !fixSmudge || diff > 1 {
					return false
				}
			}
		}
		if diff == 0 {
			continue
		}
		if diff == 1 && !hasFixedSmudge {
			hasFixedSmudge = true
			continue
		}
		return false
	}
	if fixSmudge {
		return hasFixedSmudge
	}
	return true
}

func solve(input string, part int) (string, error) {
	fixSmudge := part == 2
	patterns := parse(input)
	sum := 0
	for _, g := range patterns {
		for row := 0; row < len(g)-1; row++ {
			if mirroredOverRow(g, row, fixSmudge) {
				// The row used for summation is 1-indexed, not 0-indexed.
				sum += 100 * (row + 1)
				break // remove?
			}
		}
		for col := 0; col < len(g[0])-1; col++ {
			if mirroredOverColumn(g, col, fixSmudge) {
				// The column used for summation is 1-indexed, not 0-indexed.
				sum += col + 1
			}
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
