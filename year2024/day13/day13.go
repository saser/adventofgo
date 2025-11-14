package day13

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	buttonARegex = regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	buttonBRegex = regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	prizeRegex   = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
)

type game struct {
	P, Q int64 // button A
	R, S int64 // button B
	X, Y int64 // prize
}

func parseGame(fragment string) (game, error) {
	var g game
	var err error
	lines := strings.Split(fragment, "\n")

	matches := buttonARegex.FindStringSubmatch(lines[0])
	if matches == nil {
		return game{}, fmt.Errorf("invalid button A line: %q", lines[0])
	}
	g.P, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse AX: %v", err)
	}
	g.Q, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse AY: %v", err)
	}

	matches = buttonBRegex.FindStringSubmatch(lines[1])
	if matches == nil {
		return game{}, fmt.Errorf("invalid button B line: %q", lines[1])
	}
	g.R, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse BX: %v", err)
	}
	g.S, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse BY: %v", err)
	}

	matches = prizeRegex.FindStringSubmatch(lines[2])
	if matches == nil {
		return game{}, fmt.Errorf("invalid prize line: %q", lines[2])
	}
	g.X, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse X: %v", err)
	}
	g.Y, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse Y: %v", err)
	}

	return g, nil
}

func solveGame(g game) (a, b int64, ok bool) {
	// This looks very cryptic at first sight, but it's not too difficult.
	//
	// Each game essentially represents a system of equations of the form:
	//
	//     a*p + b*r = x
	//     a*q + b*s = y
	//
	// Since there are two unknowns and two equations, this system has either
	// one solution or no solution at all. The solution, if it exists, can be
	// found by doing some arithmetic: define a in terms of b, then find the
	// solution for b, then put that back into a.
	//
	// The full derivation of these expressions is not included here because I'm
	// too lazy to type it out, but they are not hard to do.
	//
	// The boolean expression below simply checks if the solution (1) exists and
	// (b) is integer.
	b = (g.Y*g.P - g.X*g.Q) / (g.P*g.S - g.Q*g.R)
	a = (g.X - b*g.R) / g.P
	return a, b, a*g.P+b*g.R == g.X && a*g.Q+b*g.S == g.Y
}

func parse(input string) ([]game, error) {
	var games []game
	for fragment := range strings.SplitSeq(input, "\n\n") {
		g, err := parseGame(fragment)
		if err != nil {
			return nil, fmt.Errorf("parse game: %v", err)
		}
		games = append(games, g)
	}
	return games, nil
}

func solve(input string, part int) (string, error) {
	games, err := parse(input)
	if err != nil {
		return "", err
	}
	var sum int64 = 0
	for _, g := range games {
		if part == 2 {
			const extra = 10000000000000
			g.X += extra
			g.Y += extra
		}
		a, b, ok := solveGame(g)
		if !ok {
			continue
		}
		if part == 1 && max(a, b) > 100 {
			continue
		}
		sum += 3*a + b
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
