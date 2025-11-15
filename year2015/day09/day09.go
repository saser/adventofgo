package day09

// I have a few ideas for memoization that could potentially make the solution
// faster. Basically, the idea is that the cost of a trip is:
//
//     d(a, b)
//     + (d(b, c) + d(c, d) + ...)
//
// and it's very likely that we will calculate the cost of the second part many
// times over the course of trying all possible trips. Memoization could help
// here. Currently, I don't have the energy for figuring out the details here.

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type edge struct {
	From, To string
}

func parse(input string) (map[edge]int, error) {
	distances := make(map[edge]int)
	for line := range strings.SplitSeq(input, "\n") {
		var cityA, cityB string
		var d int
		if _, err := fmt.Sscanf(line, "%s to %s = %d", &cityA, &cityB, &d); err != nil {
			return nil, fmt.Errorf("parse %q: %v", line, err)
		}
		distances[edge{From: cityA, To: cityB}] = d
		distances[edge{From: cityB, To: cityA}] = d
	}
	return distances, nil
}

func permutations(ss []string) [][]string {
	if len(ss) == 1 {
		return [][]string{ss}
	}
	var perms [][]string
	for i, s := range ss {
		ss[0], ss[i] = ss[i], ss[0]
		for _, p := range permutations(ss[1:]) {
			perms = append(perms, append([]string{s}, p...))
		}
		ss[0], ss[i] = ss[i], ss[0]
	}
	return perms
}

type cmpFunc func(a, b int) int

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func optimalTrip(distances map[edge]int, cmp cmpFunc) int {
	citySet := make(map[string]struct{})
	for e := range distances {
		citySet[e.From] = struct{}{}
	}
	cities := slices.Collect(maps.Keys(citySet))
	var best *int
	for _, itinerary := range permutations(cities) {
		trip := 0
		for i := 0; i < len(itinerary)-1; i++ {
			trip += distances[edge{From: itinerary[i], To: itinerary[i+1]}]
		}
		if best == nil {
			best = &trip
		} else {
			*best = cmp(*best, trip)
		}
	}
	return *best
}

func solve(input string, part int) (string, error) {
	distances, err := parse(input)
	if err != nil {
		return "", err
	}
	cmp := []cmpFunc{min, max}[part-1]
	return fmt.Sprint(optimalTrip(distances, cmp)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
