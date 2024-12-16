package day08

import (
	"fmt"
	"iter"

	"go.saser.se/adventofgo/asciigrid"
)

// antinodesWithoutResonance iterates over the single antinode created with
// antennae a and b, in the direction a->b.
func antinodesWithoutResonance(a, b asciigrid.Pos) iter.Seq[asciigrid.Pos] {
	// Let c = b-a, i.e. a vector from a to b.
	// The antinode occurs at b+c, i.e. b+(b-a) = 2b-a.
	return func(yield func(asciigrid.Pos) bool) {
		antinode := asciigrid.Pos{
			Row: 2*b.Row - a.Row,
			Col: 2*b.Col - a.Col,
		}
		if !yield(antinode) {
			return
		}
	}
}

// antinodesWithResonance returns all antinodes on the line that goes through a
// and b in the direction a->b.
func antinodesWithResonance(a, b asciigrid.Pos) iter.Seq[asciigrid.Pos] {
	// Let c = b-a, i.e. a vector from a to b.
	// The antinodes occurs at b+kc, where k >= 0.
	return func(yield func(asciigrid.Pos) bool) {
		c := asciigrid.Pos{
			Row: b.Row - a.Row,
			Col: b.Col - a.Col,
		}
		antinode := b
		for yield(antinode) {
			antinode.Row += c.Row
			antinode.Col += c.Col
		}
	}
}

func solve(input string, part int) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", fmt.Errorf("parse input as grid: %v", err)
	}
	antennae := make(map[byte][]asciigrid.Pos)
	for p, b := range g.All() {
		if b == '.' {
			continue
		}
		antennae[b] = append(antennae[b], p)
	}
	antinodes := make(map[asciigrid.Pos]struct{})
	findAntinodes := antinodesWithoutResonance
	if part == 2 {
		findAntinodes = antinodesWithResonance
	}
	for _, ps := range antennae {
		for i := 0; i < len(ps); i++ {
			for j := i + 1; j < len(ps); j++ {
				a := ps[i]
				b := ps[j]
				for an := range findAntinodes(a, b) {
					if !g.InBounds(an) {
						break
					}
					antinodes[an] = struct{}{}
				}
				for an := range findAntinodes(b, a) {
					if !g.InBounds(an) {
						break
					}
					antinodes[an] = struct{}{}
				}
			}
		}
	}
	return fmt.Sprint(len(antinodes)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
