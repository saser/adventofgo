package day10

import "fmt"

func lookSay(seq []byte) []byte {
	if len(seq) == 0 {
		return nil
	}
	var res []byte
	curr := seq[0]
	var n byte = 1
	for i := 1; i < len(seq); i++ {
		if c := seq[i]; c == curr {
			n++
		} else {
			res = append(res, n+'0', curr)
			curr = c
			n = 1
		}
	}
	res = append(res, n+'0', curr)
	return res
}

func solve(input string, part int) (string, error) {
	seq := []byte(input)
	rounds := 40
	if part == 2 {
		rounds = 50
	}
	for i := 0; i < rounds; i++ {
		seq = lookSay(seq)
	}
	return fmt.Sprint(len(seq)), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
