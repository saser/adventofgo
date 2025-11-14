package day02

import (
	"fmt"
	"strconv"
	"strings"
)

type grab struct {
	Red, Green, Blue int
}

type game struct {
	ID    int
	Grabs []grab
}

func (g game) IsPossibleInPart1() bool {
	const (
		redLimit   = 12
		greenLimit = 13
		blueLimit  = 14
	)
	for _, gr := range g.Grabs {
		if gr.Red > redLimit || gr.Green > greenLimit || gr.Blue > blueLimit {
			return false
		}
	}
	return true
}

func (g game) RequiredCubes() (red, green, blue int) {
	for _, gr := range g.Grabs {
		red = max(red, gr.Red)
		green = max(green, gr.Green)
		blue = max(blue, gr.Blue)
	}
	return red, green, blue
}

func parseLine(line string) (game, error) {
	var g game

	idString, rest, ok := strings.Cut(line, ": ")
	if !ok {
		return game{}, fmt.Errorf("invalid line %q", line)
	}
	var err error
	g.ID, err = strconv.Atoi(strings.TrimPrefix(idString, "Game "))
	if err != nil {
		return game{}, fmt.Errorf("invalid line %q: parse ID from %q: %v", line, idString, err)
	}

	for grabString := range strings.SplitSeq(rest, "; ") {
		var gr grab
		for chunk := range strings.SplitSeq(grabString, ", ") {
			nString, color, ok := strings.Cut(chunk, " ")
			if !ok {
				return game{}, fmt.Errorf(`invalid line %q: parse grab %q: invalid chunk %q`, line, grabString, chunk)
			}
			n, err := strconv.Atoi(nString)
			if err != nil {
				return game{}, fmt.Errorf(`invalid line %q: parse grab %q: parse chunk %q: %v`, line, grabString, chunk, err)
			}
			switch color {
			case "red":
				gr.Red = n
			case "green":
				gr.Green = n
			case "blue":
				gr.Blue = n
			}
		}
		g.Grabs = append(g.Grabs, gr)
	}

	return g, nil
}

func Part1(input string) (string, error) {
	sum := 0
	for line := range strings.SplitSeq(input, "\n") {
		g, err := parseLine(line)
		if err != nil {
			return "", err
		}
		if g.IsPossibleInPart1() {
			sum += g.ID
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	sum := 0
	for line := range strings.SplitSeq(input, "\n") {
		g, err := parseLine(line)
		if err != nil {
			return "", err
		}
		red, green, blue := g.RequiredCubes()
		sum += red * green * blue
	}
	return fmt.Sprint(sum), nil
}
