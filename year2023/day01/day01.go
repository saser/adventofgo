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

// wordDigits acts like a map, where wordDigit[i] (a string) has a value of i
// (an integer).
var wordDigits = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

func Part2(input string) (string, error) {
	sum := 0
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		// To find the first digit we create longer and longer prefixes of the
		// line, until we find a digit at the end of the prefix. That is
		// guaranteed to be the first digit.
		//
		// Examples:
		// "sjv8":
		//     "s"    => nothing
		//     "sj"   => nothing
		//     "sjv"  => nothing
		//     "sjv8" => 8
		//
		// "abcone2three":
		//     "a"       => nothing
		//     "ab"      => nothing
		//     "abc"     => nothing
		//     "abco"    => nothing
		//     "abcon"   => nothing
		//     "abcone"  => one
		first := 0
	firstLoop:
		for i := 1; i <= len(line); i++ {
			prefix := line[:i]
			if c := prefix[len(prefix)-1]; c >= '0' && c <= '9' {
				first = int(c - '0')
				break
			}
			for value, word := range wordDigits {
				if strings.HasSuffix(prefix, word) {
					first = value
					break firstLoop
				}
			}
		}
		// Similarly, we find the last digit by creating suffixes (instead of
		// prefixes) and look at the beginning of the suffix (instead of the end
		// of the prefix).
		//
		// Examples:
		// "sjv8fpr":
		//     "r"    => nothing
		//     "pr"   => nothing
		//     "fpr"  => nothing
		//     "8fpr" => 8
		//
		// "one2threefoo":
		//     "o"        => nothing
		//     "oo"       => nothing
		//     "foo"      => nothing
		//     "efoo"     => nothing
		//     "eefoo"    => nothing
		//     "reefoo"   => nothing
		//     "hreefoo"  => nothing
		//     "threefoo" => 3
		last := 0
	lastLoop:
		for i := len(line) - 1; i >= 0; i-- {
			suffix := line[i:]
			if c := suffix[0]; c >= '0' && c <= '9' {
				last = int(c - '0')
				break
			}
			for value, word := range wordDigits {
				if strings.HasPrefix(suffix, word) {
					last = value
					break lastLoop
				}
			}
		}
		sum += first*10 + last
	}
	return fmt.Sprint(sum), nil
}
