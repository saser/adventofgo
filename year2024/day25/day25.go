package day25

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/striter"
)

type pinHeight [5]int

func fits(key, lock pinHeight) bool {
	for i := range key {
		if key[i]+lock[i] > 5 {
			return false
		}
	}
	return true
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	fragments := striter.OverSplit(input, "\n\n")
	var locks, keys []pinHeight
	for fragment, ok := fragments.Next(); ok; fragment, ok = fragments.Next() {
		g, err := asciigrid.New(fragment)
		if err != nil {
			return "", fmt.Errorf("parse fragment as grid: %v", err)
		}
		isLock := g.Get(asciigrid.Pos{Row: 0, Col: 0}) == '#'
		var sig [5]int
		if isLock {
			for row := 1; row <= 5; row++ {
				done := true
				for col := 0; col < 5; col++ {
					p := asciigrid.Pos{Row: row, Col: col}
					if g.Get(p) == '#' {
						sig[col]++
						done = false
					}
				}
				if done {
					break
				}
			}
			locks = append(locks, sig)
		} else {
			for row := 5; row >= 1; row-- {
				done := true
				for col := 0; col < 5; col++ {
					p := asciigrid.Pos{Row: row, Col: col}
					if g.Get(p) == '#' {
						sig[col]++
						done = false
					}
				}
				if done {
					break
				}
			}
			keys = append(keys, sig)
		}
	}
	answer := 0
	// Not sure why, but sorting the `keys` slice first actually improves the
	// runtime by ~25%! From ~500us to ~380us on my laptop. I don't know why,
	// but I suspect it has to do with cache locality and stuff to do. Since the
	// slice is sorted, there will be large "classes" of keys that will be
	// skipped at once due to sharing similar "prefixes". This is handwave-y but
	// I don't have the energy to think more about it right now.
	slices.SortFunc(keys, func(a, b pinHeight) int {
		return cmp.Or(
			cmp.Compare(a[0], b[0]),
			cmp.Compare(a[1], b[1]),
			cmp.Compare(a[2], b[2]),
			cmp.Compare(a[3], b[3]),
			cmp.Compare(a[4], b[4]),
		)
	})
	for _, lock := range locks {
		for _, key := range keys {
			if fits(key, lock) {
				answer++
			}
		}
	}
	return fmt.Sprint(answer), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}
