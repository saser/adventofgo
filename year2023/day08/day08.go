package day08

import (
	"errors"
	"fmt"
	"regexp"

	"go.saser.se/adventofgo/striter"
)

type node struct {
	Name        string
	Left, Right string
}

var nodeRE = regexp.MustCompile(`([A-Z]{3}) = \(([A-Z]{3}), ([A-Z]{3})\)`)

func parseNode(s string) (node, error) {
	matches := nodeRE.FindStringSubmatch(s)
	if matches == nil {
		return node{}, fmt.Errorf("parse node: %q doesn't match regexp %s", s, nodeRE.String())
	}
	return node{
		Name:  matches[1],
		Left:  matches[2],
		Right: matches[3],
	}, nil
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	lines := striter.OverLines(input)
	instructions, _ := lines.Next() // Assume correct input.
	_, _ = lines.Next()             // Skip over empty line.
	nodes := make(map[string]node)  // Node name -> node
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		n, err := parseNode(line)
		if err != nil {
			return "", err
		}
		nodes[n.Name] = n
	}
	current := "AAA"
	steps := 0
	i := 0
	for current != "ZZZ" {
		instr := instructions[i]
		i++
		if i == len(instructions) {
			i = 0
		}
		switch instr {
		case 'L':
			current = nodes[current].Left
		case 'R':
			current = nodes[current].Right
		}
		steps++
	}
	return fmt.Sprint(steps), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
