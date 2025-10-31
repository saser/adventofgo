package day03

import (
	"fmt"
)

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

func (p *parser) sum1() int {
	sum := 0
	for p.pos < len(p.input) {
		// We can achieve significant speedups by simply continuing early if
		// there's no chance that the current position is the start of one of
		// the tokens we care about.
		if p.input[p.pos] != 'm' {
			p.pos++
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
		sum += a * b
	}
	return sum
}

func (p *parser) sum2() int {
	sum := 0
	enabled := true
	for p.pos < len(p.input) {
		// We can achieve significant speedups by simply continuing early if
		// there's no chance that the current position is the start of one of
		// the tokens we care about.
		if c := p.input[p.pos]; c != 'd' && c != 'm' {
			p.pos++
			continue
		}
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
		if enabled {
			sum += a * b
		}
	}
	return sum
}

func solve(input string, part int) (string, error) {
	p := &parser{input: input}
	var sum int
	if part == 1 {
		sum = p.sum1()
	} else {
		sum = p.sum2()
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
