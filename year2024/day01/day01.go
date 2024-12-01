package day01

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func parse(input string) (left, right []int, err error) {
	fields := strings.Fields(input)
	left = make([]int, len(fields)/2)
	right = make([]int, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		var err error
		j := i / 2
		left[j], err = strconv.Atoi(fields[i])
		if err != nil {
			return nil, nil, err
		}
		right[j], err = strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, nil, err
		}
	}
	return left, right, nil
}

func Part1(input string) (string, error) {
	left, right, err := parse(input)
	if err != nil {
		return "", err
	}
	slices.Sort(left)
	slices.Sort(right)
	sum := 0
	for i := range left {
		a := left[i]
		b := right[i]
		if a > b {
			sum += a - b
		} else {
			sum += b - a
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	left, right, err := parse(input)
	if err != nil {
		return "", err
	}
	frequency := make(map[int]int, len(right))
	for _, n := range right {
		frequency[n]++
	}
	sum := 0
	for _, n := range left {
		sum += n * frequency[n]
	}
	return fmt.Sprint(sum), nil
}
