package day03

import (
	"fmt"
	"regexp"
)

var instructionRE = regexp.MustCompile(`(?:do\(\)|don't\(\)|mul\((\d{1,3}),(\d{1,3})\))`)

type parser struct {
	input string
	pos   int
}

func (p *parser) literal(lit string) bool {
	for i, c := range []byte(lit) {
		if p.input[p.pos+i] != c {
			return false
		}
	}
	p.pos += len(lit)
	return true
}

func (p *parser) number() (int, bool) {
	n := 0
	digits := 0
	for range 3 {
		c := p.input[p.pos]
		if c < '0' || c > '9' {
			break
		}
		p.pos++
		digits++
		n = n*10 + int(c-'0')
	}
	return n, digits > 0
}

func (p *parser) sum(part int) int {
	sum := 0
	enabled := true
	for p.pos < len(p.input) {
		if p.literal("do()") {
			enabled = true
			continue
		}
		if p.literal("don't()") {
			enabled = false
			continue
		}
		if !p.literal("mul(") {
			p.pos++
			continue
		}
		a, ok := p.number()
		if !ok {
			p.pos++
			continue
		}
		if !p.literal(",") {
			p.pos++
			continue
		}
		b, ok := p.number()
		if !ok {
			p.pos++
			continue
		}
		if !p.literal(")") {
			p.pos++
			continue
		}
		if part == 1 || enabled {
			sum += a * b
		}
	}
	return sum
}

func solve(input string, part int) (string, error) {
	// input = `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`
	// input = `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`
	return fmt.Sprint((&parser{input: input}).sum(part)), nil
	// matches := instructionRE.FindAllStringSubmatch(input, -1)
	// sum := 0
	// enabled := true
	// for _, submatches := range matches {
	// 	if s := submatches[0]; s == "do()" || s == "don't()" {
	// 		if part == 1 {
	// 			continue
	// 		}
	// 		if s == "do()" {
	// 			enabled = true
	// 			continue
	// 		}
	// 		if s == "don't()" {
	// 			enabled = false
	// 			continue
	// 		}
	// 	}
	// 	if !enabled {
	// 		continue
	// 	}
	// 	a, err := strconv.Atoi(submatches[1])
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	b, err := strconv.Atoi(submatches[2])
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	sum += a * b
	// }
	// return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
