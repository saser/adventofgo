package day14

import (
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
)

func moveRock(g *asciigrid.Grid, p asciigrid.Pos, dir asciigrid.Direction) {
	if g.Get(p) != 'O' {
		return
	}
	for p2 := p.Step(dir); g.InBounds(p2); p, p2 = p2, p2.Step(dir) {
		if g.Get(p2) != '.' {
			break
		}
		g.Set(p, '.')
		g.Set(p2, 'O')
	}
}

func tiltNorth(g *asciigrid.Grid) {
	for row := 1; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			moveRock(g, asciigrid.Pos{Row: row, Col: col}, asciigrid.Up)
		}
	}
}

func tiltEast(g *asciigrid.Grid) {
	for col := g.NCols() - 2; col >= 0; col-- {
		for row := 0; row < g.NRows(); row++ {
			moveRock(g, asciigrid.Pos{Row: row, Col: col}, asciigrid.Right)
		}
	}
}

func tiltSouth(g *asciigrid.Grid) {
	for row := g.NRows() - 2; row >= 0; row-- {
		for col := 0; col < g.NCols(); col++ {
			moveRock(g, asciigrid.Pos{Row: row, Col: col}, asciigrid.Down)
		}
	}
}

func tiltWest(g *asciigrid.Grid) {
	for col := 1; col < g.NCols(); col++ {
		for row := 0; row < g.NRows(); row++ {
			moveRock(g, asciigrid.Pos{Row: row, Col: col}, asciigrid.Left)
		}
	}
}

func spinCycle(g *asciigrid.Grid) {
	tiltNorth(g)
	tiltWest(g)
	tiltSouth(g)
	tiltEast(g)
}

func spin(g *asciigrid.Grid, n int) {
	var seen []string
	seen = append(seen, g.String())
	idx := make(map[string]int) // grid state -> nr of spins (which is also an index into seen)
	idx[g.String()] = 0
	var cycleStart, cycleLen int
	for i := 1; i <= n; i++ {
		spinCycle(g)
		s := g.String()
		if start, ok := idx[s]; ok {
			cycleStart = start
			cycleLen = i - start
			break
		}
		seen = append(seen, s)
		idx[s] = i
	}
	final := cycleStart + ((n - cycleStart) % cycleLen)
	s := seen[final]
	*g = *asciigrid.MustNew(s)
}

func totalLoad(g *asciigrid.Grid) int {
	load := 0
	for row := 0; row < g.NRows(); row++ {
		f := g.NRows() - row
		it := g.Row(row)
		for _, tile, ok := it.Next(); ok; _, tile, ok = it.Next() {
			if tile == 'O' {
				load += f
			}
		}
	}
	return load
}

func solve(input string, part int) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", err
	}
	if part == 1 {
		tiltNorth(g)
	} else {
		spin(g, 1e9)
	}
	return fmt.Sprint(totalLoad(g)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
