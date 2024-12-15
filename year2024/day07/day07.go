package day07

import (
	"fmt"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func parse(line string) (target int, parts []int, err error) {
	fields := strings.FieldsFunc(line, func(r rune) bool { return r == ':' || r == ' ' })
	numbers := make([]int, len(fields))
	for i, s := range fields {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, nil, fmt.Errorf("parse %q: %v", line, err)
		}
		numbers[i] = n
	}
	return numbers[0], numbers[1:], nil
}

func digits(x int) int {
	n := 0
	for x > 0 {
		n++
		x /= 10
	}
	return n
}

func concat(a, b int) int {
	// Let N be the number of digits in b. Then a || b == a*10^N + b.
	//
	// Example: a = 12, b = 345, so N = 3.
	// =>     a * 10^3 = 12000
	// => a * 10^3 + b = 12345
	//
	// As a special case, if b == 0, then a || b is just a*10.
	if b == 0 {
		return a * 10
	}
	for range digits(b) {
		a *= 10
	}
	return a + b
}

func solvableAuxiliary(target int, current int, numbers []int, operators []string) bool {
	if len(numbers) == 0 {
		return target == current
	}
	next := numbers[0]
	var rest []int
	if len(numbers) > 1 {
		rest = numbers[1:]
	}
	for _, op := range operators {
		b := false
		switch op {
		case "+":
			b = solvableAuxiliary(target, current+next, rest, operators)
		case "*":
			b = solvableAuxiliary(target, current*next, rest, operators)
		case "||":
			b = solvableAuxiliary(target, concat(current, next), rest, operators)
		}
		if b {
			return true
		}
	}
	return false
}

func solvable(target int, parts []int, operators []string) bool {
	return solvableAuxiliary(target, 0, parts, operators)
}

func solve(input string, part int) (string, error) {
	operators := []string{"+", "*"}
	if part == 2 {
		operators = append(operators, "||")
	}
	lines := striter.OverLines(input)
	sum := 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		target, parts, err := parse(line)
		if err != nil {
			return "", err
		}
		if solvable(target, parts, operators) {
			sum += target
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
