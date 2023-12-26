package day10

import (
	"errors"
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
)

func findStart(g *asciigrid.Grid) asciigrid.Pos {
	for row := 0; row < g.NRows(); row++ {
		it := g.Row(row)
		for pos, tile, ok := it.Next(); ok; pos, tile, ok = it.Next() {
			if tile == 'S' {
				return pos
			}
		}
	}
	panic("unreachable")
}

func enterGoing(c byte, dir asciigrid.Direction) bool {
	switch c {
	case '|':
		return dir == asciigrid.Up || dir == asciigrid.Down
	case '-':
		return dir == asciigrid.Left || dir == asciigrid.Right
	case 'L':
		return dir == asciigrid.Down || dir == asciigrid.Left
	case 'J':
		return dir == asciigrid.Right || dir == asciigrid.Down
	case '7':
		return dir == asciigrid.Right || dir == asciigrid.Up
	case 'F':
		return dir == asciigrid.Up || dir == asciigrid.Left
	default:
		panic("unreachable")
	}
}

func goThrough(pipe byte, towards asciigrid.Direction) (asciigrid.Direction, bool) {
	switch pipe {
	case '|':
		return towards, towards == asciigrid.Up || towards == asciigrid.Down
	case '-':
		return towards, towards == asciigrid.Left || towards == asciigrid.Right
	case 'L':
		if towards == asciigrid.Down {
			return asciigrid.Right, true
		}
		if towards == asciigrid.Left {
			return asciigrid.Up, true
		}
	case 'J':
		if towards == asciigrid.Right {
			return asciigrid.Up, true
		}
		if towards == asciigrid.Down {
			return asciigrid.Left, true
		}
	case '7':
		if towards == asciigrid.Right {
			return asciigrid.Down, true
		}
		if towards == asciigrid.Up {
			return asciigrid.Left, true
		}
	case 'F':
		if towards == asciigrid.Up {
			return asciigrid.Right, true
		}
		if towards == asciigrid.Left {
			return asciigrid.Down, true
		}
	}
	return asciigrid.None, false
}

func walkLoop(g *asciigrid.Grid, start asciigrid.Pos) int {
	var p asciigrid.Pos
	var towards asciigrid.Direction
	for _, dir := range []asciigrid.Direction{
		asciigrid.Up,
		asciigrid.Down,
		asciigrid.Left,
		asciigrid.Right,
	} {
		p = start.Step(dir)
		// Assumption: first is going to be in bounds, or in other words, S is
		// not on an edge.
		c := g.Get(p)
		if c == '.' {
			continue
		}
		if newDir, ok := goThrough(c, dir); ok {
			towards = newDir
			break
		}
	}
	steps := 1
	for p != start {
		p = p.Step(towards)
		towards, _ = goThrough(g.Get(p), towards)
		steps++
	}
	return steps / 2
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	g, err := asciigrid.New(input)
	if err != nil {
		return "", err
	}
	start := findStart(g)
	return fmt.Sprint(walkLoop(g, start)), nil
	// return "", errors.New("unimplemented")
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
