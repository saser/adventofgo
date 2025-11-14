package day11

import (
	"fmt"
	"strconv"
	"strings"
)

func digits(v uint64) uint64 {
	var n uint64 = 0
	for v != 0 {
		n++
		v /= 10
	}
	return n
}

func blink2(stones []uint64, blinks int) uint64 {
	type key struct {
		Stone     uint64
		Remaining int
	}
	memo := make(map[key]uint64)
	var blinkAux func(stone uint64, blinks int) uint64
	blinkAux = func(stone uint64, blinks int) (n uint64) {
		k := key{Stone: stone, Remaining: blinks}
		if prev, seen := memo[k]; seen {
			return prev
		}
		defer func() { memo[k] = n }()
		if blinks == 0 {
			return 1
		}
		if stone == 0 {
			return blinkAux(1, blinks-1)
		}
		if n := digits(stone); n%2 == 0 {
			var p uint64 = 1
			for range n / 2 {
				p *= 10
			}
			return blinkAux(stone/p, blinks-1) + blinkAux(stone%p, blinks-1)
		}
		return blinkAux(stone*2024, blinks-1)
	}
	var sum uint64 = 0
	for _, s := range stones {
		sum += blinkAux(s, blinks)
	}
	return sum
}

func solve(input string, part int) (string, error) {
	var stones []uint64
	for field := range strings.SplitSeq(input, " ") {
		n, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return "", err
		}
		stones = append(stones, n)
	}
	blinks := 25
	if part == 2 {
		blinks = 75
	}
	return fmt.Sprint(blink2(stones, blinks)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
