package day16

import (
	"fmt"
	"log"
	"slices"

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

func solve(input string, part int) (string, error) {
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
		Cost  int
		Tiles []asciigrid.Pos
	}
	seen := make(map[key]state)
	pq := priorityqueue.NewFunc(func(x, y state) bool { return x.Cost < y.Cost })
	pq.Push(state{
		key: key{
			Pos: start,
			Dir: asciigrid.Right,
		},
		Cost:  0,
		Tiles: []asciigrid.Pos{start},
	})
	lowestScore := -1
	bestTiles := make(map[asciigrid.Pos]struct{})
	for pq.Len() > 0 {
		s := pq.Pop()
		if prev, ok := seen[s.key]; ok && s.Cost > prev.Cost {
			continue
		}
		if s.key.Pos == end {
			if part == 1 {
				return fmt.Sprint(s.Cost), nil
			}
			if lowestScore == -1 {
				lowestScore = s.Cost
			}
			if s.Cost == lowestScore {
				log.Printf("found new solution with cost %d", s.Cost)
				for _, tile := range s.Tiles {
					bestTiles[tile] = struct{}{}
				}
				log.Printf("len(bestTiles): %#+v\n", len(bestTiles))
				continue
			}
		}
		seen[s.key] = s
		pq.Push(state{
			key: key{
				Pos: s.key.Pos,
				Dir: turnCW(s.key.Dir),
			},
			Cost:  s.Cost + 1000,
			Tiles: slices.Clone(s.Tiles),
		})
		pq.Push(state{
			key: key{
				Pos: s.key.Pos,
				Dir: turnCCW(s.key.Dir),
			},
			Cost:  s.Cost + 1000,
			Tiles: slices.Clone(s.Tiles),
		})
		// Assumption: next is within bounds due to the surrounding walls.
		if next := s.key.Pos.Step(s.key.Dir); g.Get(next) != '#' {
			pq.Push(state{
				key: key{
					Pos: next,
					Dir: s.key.Dir,
				},
				Cost:  s.Cost + 1,
				Tiles: append(slices.Clone(s.Tiles), next),
			})
		}
	}
	return fmt.Sprint(len(bestTiles)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
