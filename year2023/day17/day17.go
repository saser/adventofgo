package day17

import (
	"cmp"
	"errors"
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/container/priorityqueue"
)

type direction rune

const (
	dirNone  direction = '-'
	dirUp    direction = 'U'
	dirDown  direction = 'D'
	dirLeft  direction = 'L'
	dirRight direction = 'R'
)

func (d direction) Inverse() direction {
	switch d {
	case dirNone:
		return dirNone
	case dirUp:
		return dirDown
	case dirDown:
		return dirUp
	case dirLeft:
		return dirRight
	case dirRight:
		return dirLeft
	default:
		panic(fmt.Errorf("invalid direction %v", d))
	}
}

type crucible struct {
	Pos       asciigrid.Pos
	Loss      int
	Direction direction
	Steps     int // Consecutive steps in Direction. Resets on changing direction. Is 0 if Direction is '-'.
}

func solve(input string, part int) (string, error) {
	minSteps := 1
	maxSteps := 3
	if part == 2 {
		minSteps = 4
		maxSteps = 10
	}
	g, err := asciigrid.New(input)
	if err != nil {
		return "", err
	}
	start := asciigrid.Pos{Row: 0, Col: 0}
	end := asciigrid.Pos{Row: g.NRows() - 1, Col: g.NCols() - 1}
	q := priorityqueue.NewFunc[crucible](func(x, y crucible) bool { return cmp.Less(x.Loss, y.Loss) })
	type state struct {
		Pos       asciigrid.Pos
		Direction direction
		Steps     int
	}
	seen := make(map[state]int) // State -> lowest seen heat loss for that state.
	q.Push(crucible{
		Pos:       start,
		Loss:      0,
		Direction: dirNone,
		Steps:     0,
	})
	for q.Len() > 0 {
		c := q.Pop()
		if c.Pos == end {
			return fmt.Sprint(c.Loss), nil
		}
		deltas := map[direction]asciigrid.Pos{
			dirUp:    {Row: -1, Col: 0},
			dirDown:  {Row: +1, Col: 0},
			dirLeft:  {Row: 0, Col: -1},
			dirRight: {Row: 0, Col: +1},
		}
		delete(deltas, c.Direction.Inverse())
		if c.Steps == maxSteps {
			delete(deltas, c.Direction)
		}
		for dir, delta := range deltas {
			c2 := crucible{
				Pos:       c.Pos,
				Loss:      c.Loss,
				Direction: dir,
				Steps:     0,
			}
			if dir == c.Direction {
				c2.Steps = c.Steps
			}
			rem := max(minSteps-c2.Steps, 1)
			if p := (asciigrid.Pos{Row: c2.Pos.Row + delta.Row*rem, Col: c2.Pos.Col + delta.Col*rem}); !g.InBounds(p) {
				continue
			}
			for i := 0; i < rem; i++ {
				c2.Pos.Row += delta.Row
				c2.Pos.Col += delta.Col
				c2.Loss += int(g.Get(c2.Pos)) - '0'
				c2.Steps++
			}
			s := state{
				Pos:       c2.Pos,
				Direction: c2.Direction,
				Steps:     c2.Steps,
			}
			if _, ok := seen[s]; ok {
				continue
			}
			seen[s] = c2.Loss
			q.Push(c2)
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
