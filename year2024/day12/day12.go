package day12

import (
	"fmt"

	"go.saser.se/adventofgo/asciigrid"
)

// region is a "coordinate set" that represents all the tiles in the same
// region.
type region map[asciigrid.Pos]struct{}

// findRegion performs a sort of BFS/union-find hybrid search over the grid to
// find all the regions contained within.
func findRegions(g *asciigrid.Grid) []region {
	seen := make(map[asciigrid.Pos]struct{})
	var regions []region
	for p := range g.All() {
		if _, ok := seen[p]; ok {
			continue
		}
		r := make(region)
		regions = append(regions, r)
		queue := []asciigrid.Pos{p}
		for len(queue) != 0 {
			p := queue[0]
			queue = queue[1:]
			r[p] = struct{}{}
			for _, neighbor := range p.Neighbors4() {
				if !g.InBounds(neighbor) || g.Get(neighbor) != g.Get(p) {
					continue
				}
				if _, ok := seen[neighbor]; ok {
					continue
				}
				seen[neighbor] = struct{}{}
				queue = append(queue, neighbor)
			}
		}
	}
	return regions
}

// perimeter returns the length of the perimeter of the region.
func (r region) perimeter() int {
	count := 0
	for p := range r {
		for _, n := range p.Neighbors4() {
			if _, ok := r[n]; !ok {
				count++
			}
		}
	}
	return count
}

// sides returns the number of sides of the region.
func (r region) sides() int {
	// The number of sides is equal to the number of corners. We can find
	// corners by iterating over every tile in the grid and counting how many of
	// its four corners are also interior or exterior corners for the region.
	//
	// For example:
	//
	//	XXXX
	//	XAAA
	//	XAYY
	//	XAAY
	//	XXAY
	//
	// Let's focus on just the A region, and on a specific tile (@) in that
	// region.
	//
	//	....
	//	.AAA
	//	.A..
	//	.A@.
	//	..A.
	//
	// The @ tile has four corners, marked here with numbers:
	//
	// 1-2
	// |@|
	// 3-4
	//
	// Corner 2 is an "interior" corner because @ is on the "inside" of the
	// corner:
	//
	// +-+-2
	// +A|@|
	// +-+-+
	// ..|A|
	// ..+-+
	//
	// Similarly, corner 3 is an "exterior" corner because @ is on the "outside"
	// of it:
	//
	// +-+-+
	// +A|@|
	// +-3-+
	// ..|A|
	// ..+-+
	//
	// An interior corner can easily be identified because the "outsides" of the
	// corner are edges against another region, and the "insides" are against @.
	// In this case, the "outsides" are up and right; taking a step either up or
	// right from @ will land you in another region, so there is an interior
	// corner there.
	//
	// Exterior corners are similar: the "outsides" of the corner are also
	// against another region, and only one tile of that region, while the
	// "insides" are against other tiles in the A region. In this case, the
	// "outside" is down-left; taking a step in that direction from @ will land
	// you in another region; and the "insides" are down and left; taking a
	// step in those directions from @ will land you in the same region.
	//
	// Using this logic and generalizing it to the other corners we can easily
	// count how many corners the @ tile contributes to. When we apply this
	// logic across all tiles in a region, we will accumulate the total number
	// of corners, and therefore the total number of sides.
	count := 0
	for p := range r {
		_, up := r[p.Step(asciigrid.Up)]
		_, topRight := r[p.Step(asciigrid.TopRight)]
		_, right := r[p.Step(asciigrid.Right)]
		_, bottomRight := r[p.Step(asciigrid.BottomRight)]
		_, down := r[p.Step(asciigrid.Down)]
		_, bottomLeft := r[p.Step(asciigrid.BottomLeft)]
		_, left := r[p.Step(asciigrid.Left)]
		_, topLeft := r[p.Step(asciigrid.TopLeft)]

		// Interior corners.
		if !up && !right {
			count++
		}
		if !right && !down {
			count++
		}
		if !down && !left {
			count++
		}
		if !left && !up {
			count++
		}

		// Exterior corners.
		if up && right && !topRight {
			count++
		}
		if right && down && !bottomRight {
			count++
		}
		if down && left && !bottomLeft {
			count++
		}
		if left && up && !topLeft {
			count++
		}
	}
	return count
}

func solve(input string, part int) (string, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return "", err
	}
	regions := findRegions(g)
	sum := 0
	for _, r := range regions {
		area := len(r)
		if part == 1 {
			sum += area * r.perimeter()
		} else {
			sum += area * r.sides()
		}
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
