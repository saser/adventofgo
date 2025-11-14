package day23

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"go.saser.se/adventofgo/container/set"
)

func parse(input string) (map[string]set.Set[string], error) {
	edges := make(map[string]set.Set[string])
	for line := range strings.SplitSeq(input, "\n") {
		from, to, ok := strings.Cut(line, "-")
		if !ok {
			return nil, fmt.Errorf("malformed line: %q", line)
		}
		if _, ok := edges[from]; !ok {
			edges[from] = make(set.Set[string])
		}
		edges[from].Add(to)
		if _, ok := edges[to]; !ok {
			edges[to] = make(set.Set[string])
		}
		edges[to].Add(from)
	}
	return edges, nil
}

func Part1(input string) (string, error) {
	edges, err := parse(input)
	if err != nil {
		return "", err
	}
	type component struct{ N1, N2, N3 string }
	seen := make(set.Set[component])
	answer := 0
	for n1 := range edges {
		for n2 := range edges[n1].All() {
			for n3 := range edges[n2].All() {
				connected := edges[n1].Contains(n3)
				hasT := n1[0] == 't' || n2[0] == 't' || n3[0] == 't'
				if !(connected && hasT) {
					continue
				}
				nodes := []string{n1, n2, n3}
				slices.Sort(nodes)
				c := component{
					N1: nodes[0],
					N2: nodes[1],
					N3: nodes[2],
				}
				if seen.Contains(c) {
					continue
				}
				seen.Add(c)
				answer++
			}
		}
	}
	return fmt.Sprint(answer), nil
}

// runBronKerboschWithPivot implements the algorithm described at
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm#With_pivoting
// retrieved on 2024-12-29. The method for choosing the pivot element is to just
// pick a random element in P ⋃ X. I experimented with picking the element with
// the maximum degree, but that didn't have any meaningful impact on my runtime.
func runBronKerboschWithPivot(edges map[string]set.Set[string], r, p, x set.Set[string], report func(set.Set[string])) {
	if p.Len() == 0 && x.Len() == 0 {
		report(r)
		return
	}
	var u string
	for n := range set.Union(p, x).All() {
		u = n
		break
	}
	nu := edges[u]
	for v := range set.Minus(p, nu).All() {
		nv := edges[v]
		runBronKerboschWithPivot(edges, set.Union(r, set.Of(v)), set.Intersection(p, nv), set.Intersection(x, nv), report)
		p.Delete(v)
		x.Add(v)
	}
}

// maximalCliques runs the Bron-Kerbosch algorithm
// (https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm) to gather a
// list of all the maximal cliques. The maximum clique will be the largest
// clique in this list.
func maximalCliques(edges map[string]set.Set[string]) []set.Set[string] {
	var cliques []set.Set[string]
	report := func(clique set.Set[string]) {
		cliques = append(cliques, clique.Clone())
	}
	r := make(set.Set[string])
	p := make(set.Set[string], len(edges))
	for n := range edges {
		p.Add(n)
	}
	x := make(set.Set[string])
	runBronKerboschWithPivot(edges, r, p, x, report)
	return cliques
}

func Part2(input string) (string, error) {
	edges, err := parse(input)
	if err != nil {
		return "", err
	}
	cliques := maximalCliques(edges)
	largest := slices.MaxFunc(cliques, func(a, b set.Set[string]) int { return cmp.Compare(a.Len(), b.Len()) })
	return strings.Join(slices.Sorted(largest.All()), ","), nil
}
