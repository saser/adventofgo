package day10

import (
	"fmt"
	"maps"

	"go.saser.se/adventofgo/asciigrid"
)

func traverse(g *asciigrid.Grid, p asciigrid.Pos, reachableNines map[asciigrid.Pos]map[asciigrid.Pos]struct{}) {
	if _, ok := reachableNines[p]; ok {
		return
	}
	nines := make(map[asciigrid.Pos]struct{})
	defer func() { reachableNines[p] = nines }()
	if g.Get(p) == '9' {
		nines[p] = struct{}{}
		return
	}
	for _, n := range p.Neighbors4() {
		if !g.InBounds(n) {
			continue
		}
		if g.Get(n) != g.Get(p)+1 {
			continue
		}
		traverse(g, n, reachableNines)
		maps.Copy(nines, reachableNines[n])
	}
}

func findTrails(g *asciigrid.Grid, p asciigrid.Pos, scores map[asciigrid.Pos]int) (sum int) {
	defer func() { scores[p] = sum }()
	if n, ok := scores[p]; ok {
		return n
	}
	if g.Get(p) == '9' {
		return 1
	}
	for _, n := range p.Neighbors4() {
		if !g.InBounds(n) {
			continue
		}
		if g.Get(n) != g.Get(p)+1 {
			continue
		}
		sum += findTrails(g, n, scores)
	}
	return sum
}

func Part1(input string) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %w", err)
	}
	reachableNines := make(map[asciigrid.Pos]map[asciigrid.Pos]struct{}, g.NCols()*g.NRows())
	sum := 0
	for row := 0; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			p := asciigrid.Pos{Row: row, Col: col}
			if g.Get(p) == '0' {
				traverse(g, p, reachableNines)
				sum += len(reachableNines[p])
			}
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %w", err)
	}
	scores := make(map[asciigrid.Pos]int, g.NCols()*g.NRows())
	sum := 0
	for row := 0; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			p := asciigrid.Pos{Row: row, Col: col}
			if g.Get(p) == '0' {
				sum += findTrails(g, p, scores)
			}
		}
	}
	return fmt.Sprint(sum), nil
}
