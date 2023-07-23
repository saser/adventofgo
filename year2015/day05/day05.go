package day05

import (
	"fmt"
	"strings"

	"go.saser.se/adventofgo/striter"
)

type checker interface {
	// Check returns true if the password passes all rules.
	Check(password string) bool
}

// part1 implements the rules for part 1.
type part1 struct{}

func checkVowels(password string) bool {
	count := 0
	for _, r := range password {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
		if count == 3 {
			return true
		}
	}
	return false
}

func checkLetterTwice(password string) bool {
	for i := 0; i < len(password)-1; i++ {
		if password[i] == password[i+1] {
			return true
		}
	}
	return false
}

func checkBadStrings(password string) bool {
	for _, bad := range []string{"ab", "cd", "pq", "xy"} {
		if strings.Contains(password, bad) {
			return false
		}
	}
	return true
}

func (part1) Check(password string) bool {
	return checkVowels(password) && checkLetterTwice(password) && checkBadStrings(password)
}

// part2 implements the rules for part 2.
type part2 struct{}

func checkTwoPairs(password string) bool {
	for i := range password[:len(password)-2] {
		search := password[i : i+2]
		if strings.Contains(password[i+2:], search) {
			return true
		}
	}
	return false
}

func checkThreeWithRepeat(password string) bool {
	for i := range password[:len(password)-2] {
		if password[i] == password[i+2] {
			return true
		}
	}
	return false
}

func (part2) Check(password string) bool {
	return checkTwoPairs(password) && checkThreeWithRepeat(password)
}

func solve(input string, part int) (string, error) {
	c := []checker{part1{}, part2{}}[part-1]
	count := 0
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		if c.Check(line) {
			count++
		}
	}
	return fmt.Sprint(count), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
