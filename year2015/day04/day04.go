package day04

import (
	"crypto/md5"
	"fmt"
)

func solve(input string, part int) (string, error) {
	const max = 1 << 25 // Upper limit so that we don't loop forever.
	for i := 0; i < max; i++ {
		h := md5.Sum([]byte(fmt.Sprintf("%s%d", input, i)))
		// In the hexadecimal representation of the hash, each byte makes up two
		// characters. A prefix of 5 or 6 zeroes will therefore be contained in
		// the first 3 bytes of the hash. The strategy here is therefore to
		// create an integer out of these three bytes and check whether that
		// integer is zero.
		b := uint32(h[0])<<16 | uint32(h[1])<<8 | uint32(h[2])
		// In the case of 5 bytes, the last 4 bits of b are irrelevant. We
		// therefore shift right to get rid of them.
		if part == 1 {
			b >>= 4
		}
		if b == 0 {
			return fmt.Sprint(i), nil
		}
	}
	return "", fmt.Errorf("no solution found in %d = 0x%x hashes", max, max)
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
