package day06

import (
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
)

type traveler struct {
	Pos       asciigrid.Pos
	Direction asciigrid.Direction
}

func solve(input string, part int) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	start := traveler{
		Direction: asciigrid.Up,
	}
outer:
	for row := 0; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			p := asciigrid.Pos{Row: row, Col: col}
			if g.Get(p) == '^' {
				start.Pos = p
				break outer
			}
		}
	}
	g.Set(start.Pos, '.')
	path, _ := walk(g, start)
	if part == 1 {
		return fmt.Sprint(len(path)), nil
	}
	loops := 0
	for p := range path {
		if p == start.Pos {
			continue
		}
		g.Set(p, '#')
		if _, isLoop := walk(g, start); isLoop {
			loops++
		}
		g.Set(p, '.')
	}
	return fmt.Sprint(loops), nil
}

func walk(g *asciigrid.Grid, t traveler) (positions map[asciigrid.Pos]struct{}, isLoop bool) {
	states := make(map[traveler]struct{})
	positions = make(map[asciigrid.Pos]struct{})
	for {
		if _, seen := states[t]; seen {
			return positions, true
		}
		states[t] = struct{}{}
		positions[t.Pos] = struct{}{}
		next := t.Pos.Step(t.Direction)
		if !g.InBounds(next) {
			return positions, false
		}
		if g.Get(next) == '#' {
			t.Direction = t.Direction.Turn(asciigrid.TurnClockwise90)
			continue
		}
		t.Pos = next
	}
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
