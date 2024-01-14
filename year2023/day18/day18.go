package day18

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/striter"
)

type instruction struct {
	Direction asciigrid.Direction
	N         int64
}

func parseInstruction(s string, swapped bool) (*instruction, error) {
	parts := strings.Split(s, " ")
	if got, want := len(parts), 3; got != want {
		return nil, fmt.Errorf("parse instruction from %q: splitting by space resulted in %q (len: %d); want len %d", s, parts, got, want)
	}

	i := new(instruction)

	if swapped {
		hexString := strings.Trim(parts[2], "()#")
		var err error
		i.N, err = strconv.ParseInt(hexString[:5], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("parse instruction from %q: parse hex string as step count: %v", s, err)
		}
		switch hexString[5] {
		case '0':
			i.Direction = asciigrid.Right
		case '1':
			i.Direction = asciigrid.Down
		case '2':
			i.Direction = asciigrid.Left
		case '3':
			i.Direction = asciigrid.Up
		}
	} else {
		switch dir := parts[0]; dir {
		case "U":
			i.Direction = asciigrid.Up
		case "D":
			i.Direction = asciigrid.Down
		case "L":
			i.Direction = asciigrid.Left
		case "R":
			i.Direction = asciigrid.Right
		default:
			return nil, fmt.Errorf("parse instruction from %q: invalid direction %q", s, dir)
		}
		var err error
		i.N, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse instruction from %q: invalid step count: %v", s, err)
		}
	}
	return i, nil
}

func (i *instruction) String() string {
	var dir string
	switch i.Direction {
	case asciigrid.Up:
		dir = "U"
	case asciigrid.Down:
		dir = "D"
	case asciigrid.Left:
		dir = "L"
	case asciigrid.Right:
		dir = "R"
	}
	return fmt.Sprintf("%s %d", dir, i.N)
}

// turnFrom returns +1 if turning from a to b is a clockwise turn or -1 if it is
// a counter-clockwise turn. If turning from a to b isn't a 90 degree turn,
// turnFrom panics.
func turnFrom(a, b asciigrid.Direction) int {
	type pair struct{ a, b asciigrid.Direction }
	switch (pair{a, b}) {
	case pair{asciigrid.Up, asciigrid.Right}:
		return +1
	case pair{asciigrid.Up, asciigrid.Left}:
		return -1
	case pair{asciigrid.Down, asciigrid.Right}:
		return -1
	case pair{asciigrid.Down, asciigrid.Left}:
		return +1
	case pair{asciigrid.Left, asciigrid.Up}:
		return +1
	case pair{asciigrid.Left, asciigrid.Down}:
		return -1
	case pair{asciigrid.Right, asciigrid.Up}:
		return -1
	case pair{asciigrid.Right, asciigrid.Down}:
		return +1
	default:
		panic(fmt.Errorf("invalid turn from %v to %v", a, b))
	}
}

func makeLoop(instructions []*instruction) (loop []asciigrid.Pos, clockwise bool) {
	loop = make([]asciigrid.Pos, len(instructions))
	turns := 0
	p := asciigrid.Pos{Row: 0, Col: 0}
	for i := range instructions {
		instr := instructions[i]
		p = p.StepN(instr.Direction, int(instr.N))
		loop[i] = p
		if i > 0 {
			prevInstr := instructions[i-1]
			turns += turnFrom(prevInstr.Direction, instr.Direction)
		}
	}
	return loop, turns > 0
}

func findDirection(from, to asciigrid.Pos) asciigrid.Direction {
	switch {
	case from.Row < to.Row:
		return asciigrid.Down
	case from.Row > to.Row:
		return asciigrid.Up
	case from.Col < to.Col:
		return asciigrid.Right
	case from.Col > to.Col:
		return asciigrid.Left
	default:
		panic(fmt.Errorf("findDirection: invalid combination of from=%v to=%v", from, to))
	}
}

func trenchArea(loop []asciigrid.Pos, clockwise bool) int64 {
	// We create an "outline" of the loop and then use a well-known algorithm
	// for calculating the area inside this outline. Algorithm is from
	// https://web.archive.org/web/20100405070507/http://valis.cs.uiuc.edu/~sariel/research/CG/compgeom/msg00831.html
	// which in turn I found through
	// https://stackoverflow.com/questions/451426/how-do-i-calculate-the-area-of-a-2d-polygon.

	// We're going to modify loop, so we make a clone of it.
	loop = slices.Clone(loop)
	// Add the first and second vertex to the end of the loop. This makes it so
	// that we can easily refer to the vertex before and after a given vertex.
	// It's not necessary, but it makes the code a little bit cleaner.
	loop = append(loop, loop[0], loop[1])
	// Each coordinate in the loop will become a vertex of the outline.
	outline := make([]asciigrid.Pos, 0, len(loop))
	for i := 1; i < len(loop)-1; i++ {
		start, mid, end := loop[i-1], loop[i], loop[i+1]
		a := findDirection(start, mid)
		b := findDirection(mid, end)

		// This enormous if-switch combination implements a table for mapping a
		// box coordinate to a point coordinate. Here's a quick illustration:
		//
		// Box coordinates:
		//     0   1   2
		//   +---+---+---+
		//   |   |   |   |
		// 0 | x | x | x |
		//   |   |   |   |
		//   +---+---+---+
		//   |   |   |   |
		// 1 | x | x | x |
		//   |   |   |   |
		//   +---+---+---+
		//
		// Point coordinates:
		//   0   1   2   3
		// 0 x---x---x---x
		//   |   |   |   |
		//   |   |   |   |
		//   |   |   |   |
		// 1 x---x---x---x
		//   |   |   |   |
		//   |   |   |   |
		//   |   |   |   |
		// 2 x---x---x---x
		//
		// The loop is initially interpreted as box coordinates, but we convert
		// them to the right point coordinates to build the outline.
		//
		// To understand why this table is the way it is, I recommend drawing it
		// out for yourself. It's too complex to effectively communicate it only
		// with ASCII/Unicode art.
		type ft struct{ from, to asciigrid.Direction }
		if clockwise {
			switch (ft{from: a, to: b}) {
			case ft{from: asciigrid.Up, to: asciigrid.Left}:
				mid.Row++
			case ft{from: asciigrid.Up, to: asciigrid.Right}:
				// No diff.

			case ft{from: asciigrid.Down, to: asciigrid.Left}:
				mid.Row++
				mid.Col++
			case ft{from: asciigrid.Down, to: asciigrid.Right}:
				mid.Col++

			case ft{from: asciigrid.Left, to: asciigrid.Up}:
				mid.Row++
			case ft{from: asciigrid.Left, to: asciigrid.Down}:
				mid.Row++
				mid.Col++

			case ft{from: asciigrid.Right, to: asciigrid.Up}:
				// No diff.
			case ft{from: asciigrid.Right, to: asciigrid.Down}:
				mid.Col++
			}
		} else {
			switch (ft{from: a, to: b}) {
			case ft{from: asciigrid.Up, to: asciigrid.Left}:
				mid.Col++
			case ft{from: asciigrid.Up, to: asciigrid.Right}:
				mid.Row++
				mid.Col++

			case ft{from: asciigrid.Down, to: asciigrid.Left}:
				// No diff.
			case ft{from: asciigrid.Down, to: asciigrid.Right}:
				mid.Row++

			case ft{from: asciigrid.Left, to: asciigrid.Up}:
				mid.Col++
			case ft{from: asciigrid.Left, to: asciigrid.Down}:
				// No diff.

			case ft{from: asciigrid.Right, to: asciigrid.Up}:
				mid.Row++
				mid.Col++
			case ft{from: asciigrid.Right, to: asciigrid.Down}:
				mid.Row++
			}
		}
		outline = append(outline, mid)
	}

	// Now that we have the outline, we can finally calculate the area.
	var area int64 = 0
	for i := range outline {
		pi := outline[i]
		j := i + 1
		if j == len(outline) {
			j = 0
		}
		pj := outline[j]

		area += int64(pi.Col) * int64(pj.Row)
		area -= int64(pi.Row) * int64(pj.Col)
	}
	area /= 2
	// The area may be negative, so we take the absolute value of it.
	if area < 0 {
		area *= -1
	}
	return area
}

func solve(input string, part int) (string, error) {
	swapped := part == 2
	lines := striter.OverLines(input)
	var instructions []*instruction
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		i, err := parseInstruction(line, swapped)
		if err != nil {
			return "", err
		}
		instructions = append(instructions, i)
	}
	area := trenchArea(makeLoop(instructions))
	return fmt.Sprint(area), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
