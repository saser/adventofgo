package day08

import (
	"fmt"
	"strconv"

	"go.saser.se/adventofgo/striter"
)

func solve(input string, part int) (string, error) {
	lines := striter.OverLines(input)
	sum := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		line := line
		switch part {
		case 1:
			memory, err := strconv.Unquote(line)
			if err != nil {
				panic(fmt.Errorf("invalid input line: %v", err))
			}
			sum += len(line) - len(memory)
		case 2:
			sum += len(strconv.Quote(line)) - len(line)
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
