package day01

import (
	"errors"
	"fmt"
)

func Part1(input string) (string, error) {
	floor := 0
	for _, r := range input {
		if r == '(' {
			floor++
		} else {
			floor--
		}
	}
	return fmt.Sprint(floor), nil
}

func Part2(input string) (string, error) {
	floor := 0
	for i, r := range input {
		if r == '(' {
			floor++
		} else {
			floor--
			if floor < 0 {
				return fmt.Sprint(i + 1), nil
			}
		}
	}
	return "", errors.New("basement never entered")
}
