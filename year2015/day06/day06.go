package day06

import (
	"fmt"
	"regexp"
	"strconv"

	"go.saser.se/adventofgo/geometry"
	"go.saser.se/adventofgo/striter"
)

type op int

const (
	opTurnOn op = iota
	opToggle
	opTurnOff
)

type instruction struct {
	Operation op
	From, To  geometry.Pos2
}

var instructionRE = regexp.MustCompile(`(turn on|toggle|turn off) (\d+),(\d+) through (\d+),(\d+)`)

func parse(line string) (instruction, error) {
	matches := instructionRE.FindStringSubmatch(line)
	if matches == nil {
		return instruction{}, fmt.Errorf("invalid line: %q", line)
	}
	var in instruction
	switch matches[1] {
	case "turn on":
		in.Operation = opTurnOn
	case "toggle":
		in.Operation = opToggle
	case "turn off":
		in.Operation = opTurnOff
	}
	var err error
	in.From.X, err = strconv.Atoi(matches[2])
	if err != nil {
		return instruction{}, fmt.Errorf("invalid line %q: parse from's X coordinate: %v", line, err)
	}
	in.From.Y, err = strconv.Atoi(matches[3])
	if err != nil {
		return instruction{}, fmt.Errorf("invalid line %q: parse from's Y coordinate: %v", line, err)
	}
	in.To.X, err = strconv.Atoi(matches[4])
	if err != nil {
		return instruction{}, fmt.Errorf("invalid line %q: parse to's X coordinate: %v", line, err)
	}
	in.To.Y, err = strconv.Atoi(matches[5])
	if err != nil {
		return instruction{}, fmt.Errorf("invalid line %q: parse to's Y coordinate: %v", line, err)
	}
	return in, nil
}

const (
	side = 1000
	size = side * side
)

func apply(in instruction, l lights) {
	for x := in.From.X; x <= in.To.X; x++ {
		for y := in.From.Y; y <= in.To.Y; y++ {
			i := side*y + x
			switch in.Operation {
			case opTurnOn:
				l.TurnOn(i)
			case opToggle:
				l.Toggle(i)
			case opTurnOff:
				l.TurnOff(i)
			}
		}
	}
}

type lights interface {
	TurnOn(i int)
	Toggle(i int)
	TurnOff(i int)
	Brightness() int
}

type binaryLights [size]bool

func (b *binaryLights) TurnOn(i int)  { b[i] = true }
func (b *binaryLights) Toggle(i int)  { b[i] = !b[i] }
func (b *binaryLights) TurnOff(i int) { b[i] = false }
func (b *binaryLights) Brightness() int {
	n := 0
	for _, l := range *b {
		if l {
			n++
		}
	}
	return n
}

type numericLights [size]int

func (b *numericLights) TurnOn(i int)  { b[i] += 1 }
func (b *numericLights) Toggle(i int)  { b[i] += 2 }
func (b *numericLights) TurnOff(i int) { b[i] = max(0, b[i]-1) }
func (b *numericLights) Brightness() int {
	n := 0
	for _, l := range *b {
		n += l
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(input string, part int) (string, error) {
	lines := striter.OverLines(input)
	l := []lights{&binaryLights{}, &numericLights{}}[part-1]
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		in, err := parse(line)
		if err != nil {
			return "", err
		}
		apply(in, l)
	}
	return fmt.Sprint(l.Brightness()), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
