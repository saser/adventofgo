package day23

import (
	"fmt"
	"slices"

	"go.saser.se/adventofgo/asciigrid"
)

func uphill(dir asciigrid.Direction) byte {
	switch dir {
	case asciigrid.Up:
		return 'v'
	case asciigrid.Down:
		return '^'
	case asciigrid.Left:
		return '>'
	case asciigrid.Right:
		return '<'
	default:
		panic("unreachable")
	}
}

type edge struct {
	Src, Dst asciigrid.Index
	Weight   int
}

type graph struct {
	junctions  []asciigrid.Index
	start, end asciigrid.Index
	edges      [][]edge
	isDAG      bool
}

func newGraph2(input string, isDAG bool) (*graph, error) {
	grid, err := asciigrid.New(input)
	if err != nil {
		return nil, err
	}
	junctions := findJunctions(grid)
	start := junctions[0]
	end := junctions[len(junctions)-1]
	edges := findEdges(grid, start, end)
	if !isDAG {
		for _, e := range edges {
			edges = append(edges, edge{
				Src:    e.Dst,
				Dst:    e.Src,
				Weight: e.Weight,
			})
		}
	}
	g := &graph{
		junctions: junctions,
		start:     start,
		end:       end,
		edges:     make([][]edge, grid.NRows()*grid.NCols()),
		isDAG:     isDAG,
	}
	for _, e := range edges {
		g.edges[e.Src] = append(g.edges[e.Src], e)
	}
	return g, nil
}

func findJunctions(grid *asciigrid.Grid) []asciigrid.Index {
	var junctions []asciigrid.Index
	for row := 0; row < grid.NRows(); row++ {
		it := grid.Row(row)
		for pos, tile, ok := it.Next(); ok; pos, tile, ok = it.Next() {
			if tile != '.' {
				continue
			}
			isJunction := false
			if row == 0 || row == grid.NRows()-1 {
				isJunction = true
			} else {
				slopes := 0
				for _, dir := range []asciigrid.Direction{
					asciigrid.Up,
					asciigrid.Down,
					asciigrid.Left,
					asciigrid.Right,
				} {
					switch grid.Get(pos.Step(dir)) {
					case '^', 'v', '<', '>':
						slopes++
					}
					if slopes >= 2 {
						isJunction = true
						break
					}
				}
			}
			if isJunction {
				grid.Set(pos, 'X')
				junctions = append(junctions, grid.Index(pos))
			}
		}
	}
	return junctions
}

func findEdgesRecursive(acc *[]edge, grid *asciigrid.Grid, visited map[asciigrid.Index]bool, src, end asciigrid.Index) {
	if src == end {
		return
	}
	if visited[src] {
		return
	}
	visited[src] = true
	var next []asciigrid.Index
	for _, dir := range []asciigrid.Direction{
		asciigrid.Up,
		asciigrid.Down,
		asciigrid.Left,
		asciigrid.Right,
	} {
		n := grid.Pos(src).Step(dir)
		if !grid.InBounds(n) {
			continue
		}
		if c := grid.Get(n); c == '#' || c == uphill(dir) {
			continue
		}
		dst, steps := nextJunction2(grid, src, dir)
		next = append(next, dst)
		*acc = append(*acc, edge{
			Src:    src,
			Dst:    dst,
			Weight: steps,
		})
	}
	for _, dst := range next {
		findEdgesRecursive(acc, grid, visited, dst, end)
	}
}

func findEdges(grid *asciigrid.Grid, start, end asciigrid.Index) []edge {
	var edges []edge
	findEdgesRecursive(&edges, grid, make(map[asciigrid.Index]bool), start, end)
	return edges
}

func nextJunction2(grid *asciigrid.Grid, junction asciigrid.Index, towards asciigrid.Direction) (asciigrid.Index, int) {
	p := grid.Pos(junction).Step(towards)
	cameFrom := towards.Inverse() // If we headed Left, we came from Right, etc.
	steps := 1
	for {
		if grid.Get(p) == 'X' {
			return grid.Index(p), steps
		}
		dirs := []asciigrid.Direction{
			// Potential optimization: much of the maze consists of long streaks
			// of walking in the same direction. By checking that direction
			// first we save a bit of execution time.
			cameFrom.Inverse(),
			// At least one of the below is going to be the same as
			// cameFrom.Inverse(). That's okay -- trying to find it and delete,
			// or do other tricks, is empirically slower than just checking it
			// again.
			asciigrid.Up,
			asciigrid.Down,
			asciigrid.Left,
			asciigrid.Right,
		}
		for _, dir := range dirs {
			if dir == cameFrom {
				// If we came from Left, we can't go Left, etc.
				continue
			}
			n := p.Step(dir)
			if !grid.InBounds(n) {
				continue
			}
			c := grid.Get(n)
			if c == '#' {
				continue
			}
			if c == uphill(dir) {
				continue
			}
			p = n
			cameFrom = dir.Inverse() // If we headed Left, we came from Right, etc.
			steps++
			break
		}
	}
}

func (g *graph) TopologicalOrder() []asciigrid.Index {
	// This is a slight variation of the algorithm described in
	// https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search. The
	// major changes are:
	// - We assume there are no cycles, so we do away with the temporary mark.
	// - We know there is exactly one node without incoming edges -- the start
	//   node -- so we can skip the "while" loop at the top and just visit the
	//   start node.
	l := make([]asciigrid.Index, 0, len(g.junctions))
	permanent := make(map[asciigrid.Index]bool)
	var visit func(n asciigrid.Index)
	visit = func(n asciigrid.Index) {
		if permanent[n] {
			return
		}
		for _, e := range g.edges[n] {
			visit(e.Dst)
		}
		permanent[n] = true
		l = append(l, n)
	}
	visit(g.start)
	slices.Reverse(l)
	return l
}

func (g *graph) longestHikeDAG() int {
	// This solution is based on
	// https://en.wikipedia.org/wiki/Topological_sorting#Application_to_shortest_path_finding,
	// retrieved on 2023-12-25. It is a linear time algorithm.
	shortest := map[asciigrid.Index]int{g.start: 0}
	for _, u := range g.TopologicalOrder() {
		for _, e := range g.edges[u] {
			v := e.Dst
			w := e.Weight
			d := shortest[u] - w // Note subtraction: this is where the negation of weights happens.
			if prev, ok := shortest[v]; !ok || prev > d {
				shortest[v] = d
			}
		}
	}
	return -shortest[g.end]
}

func (g *graph) LongestHike() int {
	if g.isDAG {
		return g.longestHikeDAG()
	}
	longest := -1
	// We could have used map[asciigrid.Index]bool here, but preallocating this
	// large slice is, empirically, 10x faster than using a slice.
	visited := make([]bool, len(g.edges))
	var visit func(junction asciigrid.Index, n int)
	visit = func(junction asciigrid.Index, n int) {
		if junction == g.end {
			longest = max(longest, n)
			return
		}
		for _, e := range g.edges[junction] {
			next := e.Dst
			if visited[next] {
				continue
			}
			visited[next] = true
			visit(next, n+e.Weight)
			visited[next] = false
		}
	}
	visited[g.start] = true
	visit(g.start, 0)
	return longest
}

func solve(input string, part int) (string, error) {
	isDAG := part == 1
	g2, err := newGraph2(input, isDAG)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(g2.LongestHike()), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
