package day04

import (
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
)

// findXMAS returns the number of "XMAS" words found in all possible directions
// starting from pos.
func findXMAS(g *asciigrid.Grid, pos asciigrid.Pos) int {
	sum := 0
	for _, dir := range []asciigrid.Direction{
		asciigrid.Right,
		asciigrid.BottomRight,
		asciigrid.Down,
		asciigrid.BottomLeft,
		asciigrid.Left,
		asciigrid.TopLeft,
		asciigrid.Up,
		asciigrid.TopRight,
	} {
		sum += findXMASInDirection(g, pos, dir)
	}
	return sum
}

// findXMASInDirection returns 1 if "XMAS" can be starting in pos and going in
// the given direction, and 0 otherwise. It is safe to call with positions that
// may be out of bounds in the grid.
func findXMASInDirection(g *asciigrid.Grid, pos asciigrid.Pos, dir asciigrid.Direction) int {
	const xmas = "XMAS"
	for i := range xmas {
		if !g.InBounds(pos) {
			return 0
		}
		if g.Get(pos) != xmas[i] {
			return 0
		}
		pos = pos.Step(dir)
	}
	return 1
}

// findXDashMAS returns the number of "X-MAS" found that includes pos as the A.
// It assumes that pos and all of its neighbors are in bounds of the grid.
func findXDashMAS(g *asciigrid.Grid, pos asciigrid.Pos) int {
	if g.Get(pos) != 'A' {
		return 0
	}
	// There are only four possible combinations:
	//   M.M   S.M   S.S   M.S
	//   .A.   .A.   .A.   .A.
	//   S.S   S.M   M.M   M.S
	// If pos is the A, then we can do direct checks for what the other
	// positions should be.
	type pattern struct {
		TopLeft     byte
		TopRight    byte
		BottomRight byte
		BottomLeft  byte
	}
	patterns := []pattern{
		// M.M
		// .A.
		// S.S
		{
			TopLeft:     'M',
			TopRight:    'M',
			BottomRight: 'S',
			BottomLeft:  'S',
		},
		// S.M
		// .A.
		// S.M
		{
			TopLeft:     'S',
			TopRight:    'M',
			BottomRight: 'M',
			BottomLeft:  'S',
		},
		// S.S
		// .A.
		// M.M
		{
			TopLeft:     'S',
			TopRight:    'S',
			BottomRight: 'M',
			BottomLeft:  'M',
		},
		// M.S
		// .A.
		// M.S
		{
			TopLeft:     'M',
			TopRight:    'S',
			BottomRight: 'S',
			BottomLeft:  'M',
		},
	}
	sum := 0
	for _, want := range patterns {
		got := pattern{
			TopLeft:     g.Get(pos.Step(asciigrid.TopLeft)),
			TopRight:    g.Get(pos.Step(asciigrid.TopRight)),
			BottomRight: g.Get(pos.Step(asciigrid.BottomRight)),
			BottomLeft:  g.Get(pos.Step(asciigrid.BottomLeft)),
		}
		if got == want {
			sum++
		}
	}
	return sum
}

func Part1(input string) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	sum := 0
	for row := 0; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			sum += findXMAS(g, asciigrid.Pos{Row: row, Col: col})
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	sum := 0
	for row := 1; row < g.NRows()-1; row++ {
		for col := 1; col < g.NCols()-1; col++ {
			sum += findXDashMAS(g, asciigrid.Pos{Row: row, Col: col})
		}
	}
	return fmt.Sprint(sum), nil
}
