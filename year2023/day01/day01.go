package day01

import (
	"fmt"
	"strings"
	"unicode"

	"go.saser.se/adventofgo/striter"
)

func Part1(input string) (string, error) {
	sum := 0
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		first := line[strings.IndexFunc(line, unicode.IsDigit)]
		last := line[strings.LastIndexFunc(line, unicode.IsDigit)]
		sum += int(first-'0')*10 + int(last-'0')
	}
	return fmt.Sprint(sum), nil
}

var digitValue = map[string]int{
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func Part2(input string) (string, error) {
	sum := 0
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		firstIndex := len(line)
		firstDigit := 0
		lastIndex := -1
		lastDigit := 0
		for digit := range digitValue {
			if i := strings.Index(line, digit); i != -1 && i < firstIndex {
				firstIndex = i
				firstDigit = digitValue[digit]
			}
			if i := strings.LastIndex(line, digit); i != -1 && i > lastIndex {
				lastIndex = i
				lastDigit = digitValue[digit]
			}
		}
		add := firstDigit*10 + lastDigit
		sum += add
	}
	return fmt.Sprint(sum), nil
}
