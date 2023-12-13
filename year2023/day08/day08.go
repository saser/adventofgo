package day08

import (
	"fmt"
	"regexp"

	"go.saser.se/adventofgo/striter"
)

type node struct {
	Name        string
	Left, Right string
}

var nodeRE = regexp.MustCompile(`(\w{3}) = \((\w{3}), (\w{3})\)`)

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
	var currents []string
	if part == 1 {
		currents = []string{"AAA"}
	} else {
		for name := range nodes {
			if name[2] == 'A' {
				currents = append(currents, name)
			}
		}
	}
	steps := 0
	i := 0
	for {
		instr := instructions[i]
		i++
		if i == len(instructions) {
			i = 0
		}
		onZ := 0
		for j := range currents {
			switch instr {
			case 'L':
				currents[j] = nodes[currents[j]].Left
			case 'R':
				currents[j] = nodes[currents[j]].Right
			}
			if currents[j][2] == 'Z' {
				onZ++
			}
		}
		steps++
		if onZ == len(currents) {
			break
		}
	}
	return fmt.Sprint(steps), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
