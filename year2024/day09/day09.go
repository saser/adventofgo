package day09

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func parse(input string) []int64 {
	blocks := make([]int64, 0, len(input)*9)
	isFile := true
	for i, n := range input {
		var fileID int64 = -1 // for '.'
		if isFile {
			fileID = int64(i) / 2
		}
		blocks = append(blocks, slices.Repeat([]int64{fileID}, int(n-'0'))...)
		isFile = !isFile
	}
	return blocks
}

func printBlocks(blocks []int64) {
	var sb strings.Builder
	sb.Grow(5 * len(blocks))
	for _, id := range blocks {
		if id == -1 {
			sb.WriteByte('.')
		} else {
			fmt.Fprintf(&sb, "%d", id)
		}
		sb.WriteByte(' ')
	}
	fmt.Println(sb.String())
	fmt.Println()
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	blocks := parse(input)
	i := 0
	j := len(blocks) - 1
	for {
		for blocks[i] != -1 {
			i++
		}
		for blocks[j] == -1 {
			j--
		}
		if i >= j {
			break
		}
		blocks[i], blocks[j] = blocks[j], blocks[i]
		i++
		j--
	}
	var sum int64 = 0
	for pos, id := range blocks {
		if id == -1 {
			break
		}
		sum += int64(pos) * id
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
