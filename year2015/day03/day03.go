package day03

import (
	"fmt"

	"go.saser.se/adventofgo/geometry"
)

func solve(input string, part int) (string, error) {
	presents := make(map[geometry.Pos2]struct{})
	// In part 1 there is one santa (Santa)
	// In part 2 there are two (Santa and Robo-Santa).
	santas := make([]geometry.Pos2, part)
	presents[geometry.Pos2{X: 0, Y: 0}] = struct{}{}
	for i, r := range input {
		n := i % len(santas)
		switch r {
		case '^':
			santas[n].Y++
		case '>':
			santas[n].X++
		case 'v':
			santas[n].Y--
		case '<':
			santas[n].X--
		}
		presents[santas[n]] = struct{}{}
	}
	return fmt.Sprint(len(presents)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
