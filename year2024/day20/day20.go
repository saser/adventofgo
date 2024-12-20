package day20

import (
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/container/priorityqueue"
)

func costsFrom(g *asciigrid.Grid, p asciigrid.Pos) map[asciigrid.Pos]int {
	type state struct {
		Pos  asciigrid.Pos
		Cost int
	}
	pq := priorityqueue.NewFunc(func(x, y state) bool { return x.Cost < y.Cost })
	costs := make(map[asciigrid.Pos]int)
	pq.Push(state{Pos: p})
	for pq.Len() > 0 {
		s := pq.Pop()
		if _, seen := costs[s.Pos]; seen {
			continue
		}
		costs[s.Pos] = s.Cost
		for _, n := range s.Pos.Neighbors4() {
			if g.Get(n) != '#' {
				pq.Push(state{
					Pos:  n,
					Cost: s.Cost + 1,
				})
			}
		}
	}
	return costs
}

func reachableWithCheat(g *asciigrid.Grid, p asciigrid.Pos, radius int) []asciigrid.Pos {
	reachable := make([]asciigrid.Pos, 0, radius*radius/2)
	for ahead := 1; ahead <= radius; ahead++ {
		for sideways := 0; ahead+sideways <= radius; sideways++ {
			for _, q := range []asciigrid.Pos{
				p.StepN(asciigrid.Up, ahead).StepN(asciigrid.Right, sideways),
				p.StepN(asciigrid.Right, ahead).StepN(asciigrid.Down, sideways),
				p.StepN(asciigrid.Down, ahead).StepN(asciigrid.Left, sideways),
				p.StepN(asciigrid.Left, ahead).StepN(asciigrid.Up, sideways),
			} {
				if !g.InBounds(q) {
					continue
				}
				if g.Get(q) == '#' {
					continue
				}
				reachable = append(reachable, q)
			}
		}
	}
	return reachable
}

func manhattan(p, q asciigrid.Pos) int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	return abs(p.Row-q.Row) + abs(p.Col-q.Col)
}

func findSavings(g *asciigrid.Grid, start, end asciigrid.Pos, cheatDuration int) map[int]int {
	countBySavings := make(map[int]int)
	fromStart := costsFrom(g, start)
	toEnd := costsFrom(g, end)
	baseline := fromStart[end]
	for cheatStart, initialCost := range fromStart {
		if initialCost > baseline {
			break
		}
		for _, cheatEnd := range reachableWithCheat(g, cheatStart, cheatDuration) {
			steps := manhattan(cheatStart, cheatEnd)
			totalCost := initialCost + steps + toEnd[cheatEnd]
			savings := baseline - totalCost
			if savings > 0 {
				countBySavings[savings]++
			}
		}
	}
	return countBySavings
}

func solve(input string, part int) (string, error) {
	cheatDuration := 2
	if part == 2 {
		cheatDuration = 20
	}
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	var start, end asciigrid.Pos // zero-values are never valid positions in this graph.
	for pos, b := range g.All() {
		switch b {
		case 'S':
			start = pos
		case 'E':
			end = pos
		}
		if start != (asciigrid.Pos{}) && end != (asciigrid.Pos{}) {
			break
		}
	}
	sum := 0
	for savings, count := range findSavings(g, start, end, cheatDuration) {
		if savings >= 100 {
			sum += count
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
