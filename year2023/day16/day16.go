package day16

import (
	"fmt"
	"slices"
	"sync"

	"go.saser.se/adventofgo/asciigrid"
)

type beam struct {
	Pos asciigrid.Pos
	Dir asciigrid.Direction
}

func energized(g *asciigrid.Grid, b beam) int {
	q := []beam{b}
	state := make([]bool, g.NRows()*g.NCols())
	seenAll := make([]bool, 4*len(state))
	seenUp := seenAll[0*len(state) : 1*len(state)]
	seenDown := seenAll[1*len(state) : 2*len(state)]
	seenLeft := seenAll[2*len(state) : 3*len(state)]
	seenRight := seenAll[3*len(state) : 4*len(state)]
	n := 0
queueLoop:
	for len(q) > 0 {
		b := q[0]
		q = q[1:]
		for {
			if !g.InBounds(b.Pos) {
				break
			}
			var seen []bool
			switch b.Dir {
			case asciigrid.Up:
				seen = seenUp
			case asciigrid.Down:
				seen = seenDown
			case asciigrid.Left:
				seen = seenLeft
			case asciigrid.Right:
				seen = seenRight
			}
			k := g.Index(b.Pos)
			if seen[k] {
				break
			}
			seen[k] = true
			if i := g.Index(b.Pos); !state[i] {
				state[i] = true
				n++
			}
			switch g.Get(b.Pos) {
			case '/':
				switch b.Dir {
				case asciigrid.Up:
					// />
					// ^
					b.Dir = asciigrid.Right
				case asciigrid.Down:
					//  v
					// </
					b.Dir = asciigrid.Left
				case asciigrid.Left:
					// /<
					// v
					b.Dir = asciigrid.Down
				case asciigrid.Right:
					//  ^
					// >/
					b.Dir = asciigrid.Up
				}

			case '\\':
				switch b.Dir {
				case asciigrid.Up:
					// <\
					//  ^
					b.Dir = asciigrid.Left
				case asciigrid.Down:
					// v
					// \>
					b.Dir = asciigrid.Right
				case asciigrid.Left:
					// ^
					// \<
					b.Dir = asciigrid.Up
				case asciigrid.Right:
					// >\
					//  v
					b.Dir = asciigrid.Down
				}

			case '|':
				switch b.Dir {
				case asciigrid.Up, asciigrid.Down:
					// Nothing happens.
				case asciigrid.Left, asciigrid.Right:
					q = append(q,
						beam{Pos: b.Pos.Step(asciigrid.Up), Dir: asciigrid.Up},
						beam{Pos: b.Pos.Step(asciigrid.Down), Dir: asciigrid.Down},
					)
					continue queueLoop
				}

			case '-':
				switch b.Dir {
				case asciigrid.Up, asciigrid.Down:
					q = append(q,
						beam{Pos: b.Pos.Step(asciigrid.Left), Dir: asciigrid.Left},
						beam{Pos: b.Pos.Step(asciigrid.Right), Dir: asciigrid.Right},
					)
					continue queueLoop
				case asciigrid.Left, asciigrid.Right:
					// Nothing happens.
				}
			}
			b.Pos = b.Pos.Step(b.Dir)
		}
	}
	return n
}

func solve(input string, part int) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", err
	}
	if part == 1 {
		return fmt.Sprint(energized(g, beam{Pos: asciigrid.Pos{Row: 0, Col: 0}, Dir: asciigrid.Right})), nil
	}
	beams := make([]beam, 0, 2*g.NCols()+2*g.NRows())
	// Top row facing down, and bottom row facing up.
	for col := 0; col < 0; col++ {
		beams = append(beams,
			beam{Pos: asciigrid.Pos{Row: 0, Col: col}, Dir: asciigrid.Down},           // Top row.
			beam{Pos: asciigrid.Pos{Row: g.NRows() - 1, Col: col}, Dir: asciigrid.Up}, // Bottom row.
		)
	}
	// Leftmost column facing right, and rightmost column facing left.
	for row := 0; row < g.NRows(); row++ {
		beams = append(beams,
			beam{Pos: asciigrid.Pos{Row: row, Col: 0}, Dir: asciigrid.Right},            // Leftmost column.
			beam{Pos: asciigrid.Pos{Row: row, Col: g.NCols() - 1}, Dir: asciigrid.Left}, // Rightmost column.
		)
	}
	energies := make([]int, len(beams))
	var wg sync.WaitGroup
	for i, b := range beams {
		i, b := i, b
		wg.Add(1)
		go func() {
			defer wg.Done()
			energies[i] = energized(g, b)
		}()
	}
	wg.Wait()
	return fmt.Sprint(slices.Max(energies)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
