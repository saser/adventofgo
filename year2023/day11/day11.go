package day11

import (
	"fmt"
	"slices"

	"go.saser.se/adventofgo/asciigrid"
)

type image struct {
	emptyRows, emptyCols []int
	galaxies             []asciigrid.Pos
	expansionFactor      int64
}

func parse(input string, expansionFactor int64) (*image, error) {
	g, err := asciigrid.New(input)
	if err != nil {
		return nil, err
	}
	img := &image{
		expansionFactor: expansionFactor,
	}
	for row := 0; row < g.NRows(); row++ {
		it := g.Row(row)
		hasGalaxy := false
		for pos, tile, ok := it.Next(); ok; pos, tile, ok = it.Next() {
			if tile == '#' {
				hasGalaxy = true
				img.galaxies = append(img.galaxies, pos)
			}
		}
		if !hasGalaxy {
			img.emptyRows = append(img.emptyRows, row)
		}
	}
	for col := 0; col < g.NCols(); col++ {
		it := g.Col(col)
		hasGalaxy := false
		for _, tile, ok := it.Next(); ok; _, tile, ok = it.Next() {
			if tile == '#' {
				hasGalaxy = true
				// We don't append to img.galaxies here -- we already saw this
				// exact galaxy in the loop over rows.
			}
		}
		if !hasGalaxy {
			img.emptyCols = append(img.emptyCols, col)
		}
	}
	return img, nil
}

func (i *image) SumShortestPaths() int64 {
	sum := int64(0)
	for j := 0; j < len(i.galaxies); j++ {
		for k := j + 1; k < len(i.galaxies); k++ {
			sum += i.shortestPath(i.galaxies[j], i.galaxies[k])
		}
	}
	return sum
}

func (i *image) shortestPath(src, dst asciigrid.Pos) int64 {
	// Without considering expanding empty space, the shortest path between two
	// galaxies is simply the Manhattan distance between them, since the
	// shortest path is allowed to pass through other galaxies.
	//
	// When adding expanding empty spaces into the mix, we can observe that a
	// shortest path consists of 0 or more steps in each of one horizontal and
	// one vertical direction. For each empty column we pass through, we count
	// its distance as expansionFactor rather than 1. The same thing happens for
	// each empty row.
	//
	// We can (ab)use the fact that the emptyRows and emptyCols slices are
	// ordered in ascending order. Let's consider columns first (the logic is
	// equivalent for rows). Looking at the example in the problem description,
	// we'd get:
	//
	//     emptyCols = [2, 5, 8]
	//
	// If we're finding the shortest path between galaxies 1 and 6 (again, from
	// the problem description), we'd see that the minimum and maximum columns
	// of them are 3 and 9. If we imagine where they'd be placed in emptyCols, we see:
	//
	//     emptyCols = [2, (3) 5, 8 (9)]
	//
	// so the shortest path will pass through empty columns 5 and 8. These
	// columns will count as expansionFactor rather than 1 in the length of the
	// shortest path.
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	// First calculate the manhattan distance, and then add the expanding empty
	// space.
	pathLen := int64(abs(src.Col-dst.Col) + abs(src.Row-dst.Row))

	minCol := min(src.Col, dst.Col)
	maxCol := max(src.Col, dst.Col)
	// slices.IndexFunc iterates through emptyCols and finds the first column
	// between minCol and maxCol, if any.
	first := slices.IndexFunc(i.emptyCols, func(col int) bool { return col > minCol })
	if first != -1 {
		// We then find the subslice of emptyCols that contains values between
		// minCol and maxCol.
		rest := i.emptyCols[first:]
		// n will eventually hold the number of values between minCol and
		// maxCol.
		n := slices.IndexFunc(rest, func(col int) bool { return col > maxCol })
		if n == -1 {
			n = len(rest)
		}
		// We multiply by expansionFactor-1 because these columns have already
		// been counted once.
		pathLen += int64(n) * (i.expansionFactor - 1)
	}

	// Now we just do the same for rows.
	minRow := min(src.Row, dst.Row)
	maxRow := max(src.Row, dst.Row)
	first = slices.IndexFunc(i.emptyRows, func(row int) bool { return row > minRow })
	if first != -1 {
		rest := i.emptyRows[first:]
		n := slices.IndexFunc(rest, func(row int) bool { return row > maxRow })
		if n == -1 {
			n = len(rest)
		}
		pathLen += int64(n) * (i.expansionFactor - 1)
	}
	return pathLen
}

func solve(input string, part int) (string, error) {
	expansionFactor := int64(2)
	if part == 2 {
		expansionFactor = 1e6
	}
	img, err := parse(input, expansionFactor)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(img.SumShortestPaths()), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
