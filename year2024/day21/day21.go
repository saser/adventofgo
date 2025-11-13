package day21

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

func solve(input string, part int) (string, error) {
	dirpads := 1 /*me*/ + 2 /*robots*/
	if part == 2 {
		dirpads = 1 /*me*/ + 25 /*robots*/
	}
	layers := append(slices.Repeat([]pad{dirpad}, dirpads), keypad)
	m := make(memo)
	var sum uint64 = 0
	for _, target := range strings.Split(strings.TrimSpace(input), "\n") {
		var n uint64 = 0
		for _, r := range target[:3] {
			n = n*10 + uint64(r-'0')
		}
		sum += n * m.shortest(target, layers)
	}
	return fmt.Sprint(sum), nil
}

type pos struct {
	Row, Col int
}

type pad map[rune]pos

var keypad = pad{
	'7': {Row: 0, Col: 0},
	'8': {Row: 0, Col: 1},
	'9': {Row: 0, Col: 2},

	'4': {Row: 1, Col: 0},
	'5': {Row: 1, Col: 1},
	'6': {Row: 1, Col: 2},

	'1': {Row: 2, Col: 0},
	'2': {Row: 2, Col: 1},
	'3': {Row: 2, Col: 2},

	'.': {Row: 3, Col: 0},
	'0': {Row: 3, Col: 1},
	'A': {Row: 3, Col: 2},
}

var dirpad = pad{
	'.': {Row: 0, Col: 0},
	'^': {Row: 0, Col: 1},
	'A': {Row: 0, Col: 2},

	'<': {Row: 1, Col: 0},
	'v': {Row: 1, Col: 1},
	'>': {Row: 1, Col: 2},
}

type key struct {
	Src, Dst pos
	Layer    int
}

type memo map[key]uint64

// shortestBetween returns the number of keypresses required at the top (first)
// layer to move from src to dst and press the button at dst on the bottom
// (last) layer.
func (m memo) shortestBetween(src, dst pos, layers []pad) (cost uint64) {
	k := key{src, dst, len(layers)}
	if v, seen := m[k]; seen {
		return v
	}
	defer func() { m[k] = cost }()

	pad := layers[len(layers)-1]
	cost = math.MaxUint64

	type state struct {
		Pos     pos
		Presses string
	}
	queue := []state{{Pos: src}}
	for len(queue) != 0 {
		s := queue[0]
		queue = queue[1:]
		if s.Pos == pad['.'] {
			continue
		}
		if s.Pos == dst {
			cost = min(cost, m.shortest(s.Presses+"A", layers[:len(layers)-1]))
			continue
		}
		// Only ever move in the direction of the destination; no other moves
		// are going to be beneficial.
		r, c := s.Pos.Row, s.Pos.Col
		if r < dst.Row {
			queue = append(queue, state{Pos: pos{Row: r + 1, Col: c}, Presses: s.Presses + "v"})
		}
		if r > dst.Row {
			queue = append(queue, state{Pos: pos{Row: r - 1, Col: c}, Presses: s.Presses + "^"})
		}
		if c < dst.Col {
			queue = append(queue, state{Pos: pos{Row: r, Col: c + 1}, Presses: s.Presses + ">"})
		}
		if c > dst.Col {
			queue = append(queue, state{Pos: pos{Row: r, Col: c - 1}, Presses: s.Presses + "<"})
		}
	}
	return cost
}

// shortest returns the number of keypresses required on the top (first) layer
// to punch in the given target sequence at the bottom (last) layer.
func (m memo) shortest(target string, layers []pad) (cost uint64) {
	if len(layers) == 1 {
		return uint64(len(target))
	}
	pad := layers[len(layers)-1]
	src := pad['A']
	for _, r := range target {
		dst := pad[r]
		cost += m.shortestBetween(src, dst, layers)
		src = dst
	}
	return cost
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
