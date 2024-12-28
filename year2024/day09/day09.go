package day09

import (
	"fmt"
	"slices"
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

func defrag(blocks []int64) {
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		for i < j && blocks[i] != -1 {
			i++
		}
		for i < j && blocks[j] == -1 {
			j--
		}
		if i == j {
			return
		}
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
}

func moveFiles(blocks []int64) {
	seen := make(map[int64]bool)
	for j := len(blocks) - 1; j >= 0; j-- {
		id := blocks[j]
		if id == -1 {
			continue
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		fileEnd := j
		fileStart := fileEnd - 1
		for fileStart >= 0 && blocks[fileStart] == id {
			fileStart--
		}
		fileStart++
		fileLen := fileEnd - fileStart + 1
		j = fileStart

		for i := 0; i < fileStart; i++ {
			if blocks[i] != -1 {
				continue
			}
			emptyStart := i
			emptyEnd := i
			for emptyEnd < fileStart && blocks[emptyEnd] == -1 {
				emptyEnd++
			}
			emptyEnd--
			emptyLen := emptyEnd - emptyStart + 1
			if emptyLen < fileLen {
				i = emptyEnd
				continue
			}
			for i := range fileLen {
				blocks[emptyStart+i] = id
				blocks[fileStart+i] = -1
			}
			break
		}
	}
}

func checksum(blocks []int64) uint64 {
	var sum uint64 = 0
	for pos, id := range blocks {
		if id == -1 {
			continue
		}
		sum += uint64(pos) * uint64(id)
	}
	return sum
}

func solve(input string, part int) (string, error) {
	blocks := parse(input)
	switch part {
	case 1:
		defrag(blocks)
	case 2:
		moveFiles(blocks)
	}
	return fmt.Sprint(checksum(blocks)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
