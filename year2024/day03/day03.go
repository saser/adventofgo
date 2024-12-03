package day03

import (
	"fmt"
	"regexp"
	"strconv"
)

var instructionRE = regexp.MustCompile(`(?:do\(\)|don't\(\)|mul\((\d{1,3}),(\d{1,3})\))`)

func solve(input string, part int) (string, error) {
	matches := instructionRE.FindAllStringSubmatch(input, -1)
	sum := 0
	enabled := true
	for _, submatches := range matches {
		if s := submatches[0]; s == "do()" || s == "don't()" {
			if part == 1 {
				continue
			}
			if s == "do()" {
				enabled = true
				continue
			}
			if s == "don't()" {
				enabled = false
				continue
			}
		}
		if !enabled {
			continue
		}
		a, err := strconv.Atoi(submatches[1])
		if err != nil {
			return "", err
		}
		b, err := strconv.Atoi(submatches[2])
		if err != nil {
			return "", err
		}
		sum += a * b
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
