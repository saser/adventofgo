package day14

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"go.saser.se/adventofgo/striter"
)

const (
	cols = 101
	rows = 103
)

type robot struct {
	X, Y   int
	DX, DY int
}

func parse(line string) (robot, error) {
	fields := strings.FieldsFunc(line, func(r rune) bool { return !(unicode.IsDigit(r) || r == '-') })
	var r robot
	var err error
	r.X, err = strconv.Atoi(fields[0])
	if err != nil {
		return robot{}, fmt.Errorf("parse robot from %q: parse X coordinate: %v", line, err)
	}
	r.Y, err = strconv.Atoi(fields[1])
	if err != nil {
		return robot{}, fmt.Errorf("parse robot from %q: parse Y coordinate: %v", line, err)
	}
	r.DX, err = strconv.Atoi(fields[2])
	if err != nil {
		return robot{}, fmt.Errorf("parse robot from %q: parse X velocity: %v", line, err)
	}
	if r.DX < 0 {
		r.DX += cols
	}
	r.DY, err = strconv.Atoi(fields[3])
	if err != nil {
		return robot{}, fmt.Errorf("parse robot from %q: parse Y velocity: %v", line, err)
	}
	if r.DY < 0 {
		r.DY += rows
	}
	return r, nil
}

func (r robot) Step(n int) robot {
	r2 := r
	r2.X = (r.X + (n%cols)*r.DX) % cols
	r2.Y = (r.Y + (n%rows)*r.DY) % rows
	return r2
}

func Part1(input string) (string, error) {
	lines := striter.OverLines(input)
	var upperLeft, upperRight, lowerLeft, lowerRight int
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		r, err := parse(line)
		if err != nil {
			return "", fmt.Errorf("parse line: %v", err)
		}
		r = r.Step(100)
		divX := cols / 2
		divY := rows / 2
		switch {
		case r.X < divX && r.Y < divY:
			upperLeft++
		case r.X > divX && r.Y < divY:
			upperRight++
		case r.X < divX && r.Y > divY:
			lowerLeft++
		case r.X > divX && r.Y > divY:
			lowerRight++
		}
	}
	return fmt.Sprint(upperLeft * upperRight * lowerLeft * lowerRight), nil
}

func printRobots(robots []robot) {
	grid := make([][]byte, rows)
	for y := range grid {
		grid[y] = slices.Repeat([]byte{'.'}, cols)
	}
	for _, r := range robots {
		grid[r.Y][r.X] = '#'
	}
	for _, row := range grid {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func hasChristmasTree(robots []robot) bool {
	// Heuristic: if any row has 31 adjacent robots, then there is a Christmas
	// tree picture. This heuristic was created by solving the problem first
	// (with a different, slower heuristic), looking at the Christmas tree
	// picture, and counting the number of robots in the horizontal parts of the
	// "frame" around the tree.
	grid := make([][]byte, rows)
	for y := range grid {
		grid[y] = slices.Repeat([]byte{'.'}, cols)
	}
	for _, r := range robots {
		grid[r.Y][r.X] = '#'
	}
	for _, row := range grid {
		adjacent := 0
		for _, b := range row {
			if b == '.' {
				adjacent = 0
				continue
			}
			if b == '#' {
				adjacent++
				if adjacent == 31 {
					return true
				}
			}
		}
	}
	return false
}

func Part2(input string) (string, error) {
	lines := striter.OverLines(input)
	var robots []robot
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		r, err := parse(line)
		if err != nil {
			return "", fmt.Errorf("parse line: %v", err)
		}
		robots = append(robots, r)
	}
	const limit = 100_000
	for step := range limit {
		if hasChristmasTree(robots) {
			// printRobots(robots)
			return fmt.Sprint(step), nil
		}
		for i := range robots {
			robots[i] = robots[i].Step(1)
		}
	}
	return "", fmt.Errorf("no christmas tree found in %v steps", limit)
}
