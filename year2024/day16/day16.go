package day16

import (
	"errors"
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/container/priorityqueue"
)

func turnCW(dir asciigrid.Direction) asciigrid.Direction {
	switch dir {
	case asciigrid.Up:
		return asciigrid.Right
	case asciigrid.Right:
		return asciigrid.Down
	case asciigrid.Down:
		return asciigrid.Left
	case asciigrid.Left:
		return asciigrid.Up
	default:
		panic(fmt.Errorf("invalid direction: %v", dir))
	}
}

func turnCCW(dir asciigrid.Direction) asciigrid.Direction {
	switch dir {
	case asciigrid.Up:
		return asciigrid.Left
	case asciigrid.Left:
		return asciigrid.Down
	case asciigrid.Down:
		return asciigrid.Right
	case asciigrid.Right:
		return asciigrid.Up
	default:
		panic(fmt.Errorf("invalid direction: %v", dir))
	}
}

func printMap(g *asciigrid.Grid, pos asciigrid.Pos, dir asciigrid.Direction) {
	defer g.Set(pos, g.Get(pos))
	g.Set(pos, map[asciigrid.Direction]byte{
		asciigrid.Up:    '^',
		asciigrid.Right: '>',
		asciigrid.Down:  'v',
		asciigrid.Left:  '<',
	}[dir])
	fmt.Println(g)
	fmt.Println()
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	var (
		start      asciigrid.Pos
		foundStart bool
		end        asciigrid.Pos
		foundEnd   bool
	)
	for p, b := range g.All() {
		if b == 'S' {
			start = p
			foundStart = true
		}
		if b == 'E' {
			end = p
			foundEnd = true
		}
		if foundStart && foundEnd {
			break
		}
	}
	type key struct {
		Pos asciigrid.Pos
		Dir asciigrid.Direction
	}
	type state struct {
		key
		Cost int
	}
	seen := make(map[key]state)
	pq := priorityqueue.NewFunc(func(x, y state) bool { return x.Cost < y.Cost })
	pq.Push(state{
		key: key{
			Pos: start,
			Dir: asciigrid.Left,
		},
		Cost: 0,
	})
	for pq.Len() > 0 {
		s := pq.Pop()
		if _, ok := seen[s.key]; ok {
			continue
		}
		printMap(g, s.key.Pos, s.key.Dir)
		if s.key.Pos == end {
			return fmt.Sprint(s.Cost), nil
		}
		seen[s.key] = s
		pq.Push(state{
			key: key{
				Pos: s.key.Pos,
				Dir: turnCW(s.key.Dir),
			},
			Cost: s.Cost + 1000,
		})
		pq.Push(state{
			key: key{
				Pos: s.key.Pos,
				Dir: turnCCW(s.key.Dir),
			},
			Cost: s.Cost + 1000,
		})
		// Assumption: next is within bounds due to the surrounding walls.
		if next := s.key.Pos.Step(s.key.Dir); g.Get(next) != '#' {
			pq.Push(state{
				key: key{
					Pos: next,
					Dir: s.key.Dir,
				},
				Cost: s.Cost + 1,
			})
		}
	}
	return "", errors.New("no solution found")
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
