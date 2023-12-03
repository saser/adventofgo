package day03

import (
	"fmt"

	"go.saser.se/adventofgo/geometry"
	"go.saser.se/adventofgo/striter"
)

type number struct {
	// The position of this number in the input, as determined by the leftmost
	// digit. X = column from left to right, Y = row from top to bottom.
	Pos geometry.Pos2
	// How many digits are in this number.
	Len int
	// The actual number.
	V int
}

func (n number) Adjacent() []geometry.Pos2 {
	// The number of adjacent positions are:
	//     (a) n.Len + 2 (above the number, including diagonally)
	//     (b) + 2       (horizontally left of the leftmost digit and horizontally right of the rightmost digit)
	//     (c) n.Len + 2 (below the number, including diagonally)
	// Illustration:
	//    aaaaa
	//    b178b
	//    ccccc
	adj := make([]geometry.Pos2, 0, 2*n.Len+2+4)
	// (a) above the number, including diagonally
	for col := n.Pos.X - 1; col <= n.Pos.X+n.Len; col++ {
		adj = append(adj, geometry.Pos2{X: col, Y: n.Pos.Y - 1})
	}
	// (b) horizontally left of the leftmost digit and horizontally right of the rightmost digit
	adj = append(adj,
		geometry.Pos2{X: n.Pos.X - 1, Y: n.Pos.Y},     // left
		geometry.Pos2{X: n.Pos.X + n.Len, Y: n.Pos.Y}, // right
	)
	// (c) below the number, including diagonally
	for col := n.Pos.X - 1; col <= n.Pos.X+n.Len; col++ {
		adj = append(adj, geometry.Pos2{X: col, Y: n.Pos.Y + 1})
	}
	return adj
}

type symbol struct {
	// The position of this symbol in the input. X = column from left to right,
	// Y = row from top to bottom.
	Pos geometry.Pos2
	// The actual symbol.
	V byte
}

type schematic struct {
	Numbers map[geometry.Pos2]number
	Symbols map[geometry.Pos2]symbol
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func parse(input string) schematic {
	s := schematic{
		Numbers: make(map[geometry.Pos2]number),
		Symbols: make(map[geometry.Pos2]symbol),
	}
	lines := striter.OverLines(input)
	row := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		col := 0
		for col < len(line) {
			c := line[col]
			if c == '.' {
				// Dots are just skipped over.
				col++
				continue
			} else if isDigit(c) {
				// We've found the start of a number. Consume digits for as long
				// as we can and store the final number.
				num := number{
					Pos: geometry.Pos2{X: col, Y: row},
					Len: 0,
					V:   0,
				}
				for col < len(line) && isDigit(line[col]) {
					c := line[col]
					num.Len++
					num.V = 10*num.V + int(c-'0')
					col++
				}
				// col now either is outside the line or points to the first
				// non-digit character. Either way, we've finished constructing
				// the number, so we store it in the schematic and then continue
				// the outer loop.
				s.Numbers[num.Pos] = num
				continue
			} else {
				// We've found a symbol, which we assume is only ever 1
				// character in length. We store it in the schematic, advance
				// col, and continue the outer loop.
				sym := symbol{
					Pos: geometry.Pos2{X: col, Y: row},
					V:   c,
				}
				s.Symbols[sym.Pos] = sym
				col++
				continue
			}
		}
		row++
	}
	return s
}

func Part1(input string) (string, error) {
	schema := parse(input)
	sum := 0
	for _, num := range schema.Numbers {
		for _, pos := range num.Adjacent() {
			if _, ok := schema.Symbols[pos]; ok {
				sum += num.V
				break
			}
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	schema := parse(input)
	numbersOnRow := make(map[int][]number)
	for pos, num := range schema.Numbers {
		numbersOnRow[pos.Y] = append(numbersOnRow[pos.Y], num)
	}
	sum := 0
symLoop:
	for pos, sym := range schema.Symbols {
		if sym.V != '*' {
			continue
		}
		ratio := 1
		// The numbers we've seen adjecent to this gear, so we avoid
		// double-counting numbers that might be considered adjacent in multiple
		// positions.
		// Example:
		//     617
		//      *
		// All of the three digits in 617 are in positions adjacent to the gear.
		// If we didn't keep track of already having 617, we could end up
		// counting it three times.
		seen := make(map[geometry.Pos2]bool)
		for _, delta := range []geometry.Pos2{
			{X: -1, Y: -1}, // upper left
			{X: +0, Y: -1}, // up
			{X: +1, Y: -1}, // upper right

			{X: -1, Y: +0}, // left
			{X: +1, Y: +0}, // right

			{X: -1, Y: +1}, // lower left
			{X: +0, Y: +1}, // down
			{X: +1, Y: +1}, // lower right
		} {
			adj := pos.Add(delta)
			for _, num := range numbersOnRow[adj.Y] {
				if adj.X >= num.Pos.X && adj.X < num.Pos.X+num.Len && !seen[num.Pos] {
					seen[num.Pos] = true
					ratio *= num.V
					// There can be only one number that matches this adjacent
					// position, so we can break out of the loop as soon as
					// we've found one.
					break
				}
			}
			if len(seen) > 2 {
				continue symLoop
			}
		}
		if len(seen) == 2 {
			sum += ratio
		}
	}
	return fmt.Sprint(sum), nil
}
