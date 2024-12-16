package day16

import (
	"fmt"
	"slices"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/container/priorityqueue"
)

func solve(input string, part int) (string, error) {
	// Parse the grid and pick out the starting and ending positions.
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

	// Run Dijkstra's algorithm, potentially many times, to find the/all
	// shortest path(s) from the start to the end.
	type key struct {
		Pos asciigrid.Pos
		Dir asciigrid.Direction
	}
	type state struct {
		Key   key
		Cost  int
		Tiles []asciigrid.Pos
	}
	seen := make(map[key]int) // key -> cost of shortest path there
	pq := priorityqueue.NewFunc(func(x, y state) bool { return x.Cost < y.Cost })
	pq.Push(state{
		Key:   key{Pos: start, Dir: asciigrid.Right},
		Cost:  0,
		Tiles: []asciigrid.Pos{start},
	})
	bestTiles := make(map[asciigrid.Pos]struct{})
	for pq.Len() > 0 {
		s := pq.Pop()
		if s.Key.Pos == end {
			if part == 1 {
				return fmt.Sprint(s.Cost), nil
			}
			for _, tile := range s.Tiles {
				bestTiles[tile] = struct{}{}
			}
			continue
		}
		for _, dir := range []asciigrid.Direction{
			s.Key.Dir, // keep current direction
			s.Key.Dir.Turn(asciigrid.TurnClockwise90),        // turn right relative to current direction
			s.Key.Dir.Turn(asciigrid.TurnCounterClockwise90), // turn left relative to current direction
		} {
			// Assumption: next is within bounds due to the surrounding walls.
			if next := s.Key.Pos.Step(dir); g.Get(next) != '#' {
				s2 := state{
					Key:   key{Pos: next, Dir: dir},
					Cost:  s.Cost + 1,
					Tiles: append(slices.Clone(s.Tiles), next),
				}
				if s2.Key.Dir != s.Key.Dir { // we turned
					s2.Cost += 1000
				}
				if prev, ok := seen[s2.Key]; ok && s2.Cost > prev {
					continue
				}
				seen[s2.Key] = s2.Cost
				pq.Push(s2)
			}
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
