package day22

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

// This was heavily inspired by
// https://github.com/goggle/AdventOfCode2024.jl/blob/3d7f58e2ba7f108d4b8e0d8faf4debab08825579/src/day22.jl.

func parse(input string) []uint64 {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	secrets := make([]uint64, len(lines))
	for i, line := range lines {
		var n uint64 = 0
		for _, r := range line {
			n = n*10 + uint64(r-'0')
		}
		secrets[i] = n
	}
	return secrets
}

// sequence is a type for the sequences of 4 price differences. Each element is
// in the range [-9; +9], i.e. one of 19 values.
type sequence [4]int8

// Index returns an int that uniquely identifies this sequence and that can be
// used as an index into a slice or similar.
func (s *sequence) Index() int {
	const (
		pow0 = 1
		pow1 = pow0 * 19
		pow2 = pow1 * 19
		pow3 = pow2 * 19
	)
	d0 := int(s[0] + 9)
	d1 := int(s[1] + 9)
	d2 := int(s[2] + 9)
	d3 := int(s[3] + 9)
	return d0*pow0 + d1*pow1 + d2*pow2 + d3*pow3
}

// Add adds the diff to the sequence as the latest price difference.
func (s *sequence) Add(diff int8) {
	s[0], s[1], s[2], s[3] = s[1], s[2], s[3], diff
}

// evolve returns the next number from n.
func evolve(n uint64) uint64 {
	// As it turns out, all of the described operations can be implemented by
	// direct manipulation of bits.
	//
	// - Mixing is simply bitwise XOR, as stated in the problem description.
	// - Pruning n results in n % 16777216, but that number happens to be 2^24,
	//   and the result of n % 2^p is the same as keeping the p lowest bits of
	//   n, which can be achieved with bitwise AND.
	// - Multiplying by 64 is the same as multiplying by 2^6, which is the same
	//   as left-shifting by 6.
	// - Dividing by 32 is the same as dividing by 2^5, which is the same as
	//   right-shifting by 5.
	// - Multiplying by 2048 is the same as multiplying by 2^11, which is the
	//   same as left-shifting by 11.
	//
	// However, after implementing this function with both bit manipulation and
	// regular math functions, there is no meaningful difference in runtime
	// between the two approaches. Comparing the generated amd64 assembly they
	// are actually identical and are implemented with bit shifts and bitwise
	// AND operators. So the compiler has recognized that these operations can
	// be implemented with bit fiddling -- good job compiler!
	n = (n ^ (n * 64)) % 16777216
	n = (n ^ (n / 32)) % 16777216
	n = (n ^ (n * 2048)) % 16777216
	return n
}

func Part1(input string) (string, error) {
	var sum atomic.Uint64
	secrets := parse(input)
	const workers = 16
	var wg sync.WaitGroup
	for chunk := range slices.Chunk(secrets, len(secrets)/workers) {
		wg.Go(func() {
			for _, n := range chunk {
				for range 2000 {
					n = evolve(n)
				}
				sum.Add(n)
			}
		})
	}
	wg.Wait()
	return fmt.Sprint(sum.Load()), nil
	// for line := range strings.SplitSeq(input, "\n") {
	// 	n, err := strconv.ParseUint(line, 10, 64)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	for range 2000 {
	// 		n = evolve(n)
	// 	}
	// 	sum.Add(n)
	// }
	// return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	// There are k = 19^4 = 130321 different
	// possible sequences.
	const k = 19 * 19 * 19 * 19

	// Profits maps from a sequence's index to the total amount of profits
	// given by that sequence, across all buyers.
	var profits [k]*atomic.Uint32
	for i := range profits {
		profits[i] = new(atomic.Uint32)
	}

	secrets := parse(input)
	const workers = 16
	var wg sync.WaitGroup
	for chunk := range slices.Chunk(secrets, len(secrets)/workers) {
		wg.Go(func() {
			var seen [k]bool
			for _, secret := range chunk {
				clear(seen[:])
				price := int8(secret % 10)
				seq := new(sequence)
				for i := range 2000 {
					secret = evolve(secret)
					nextPrice := int8(secret % 10)
					seq.Add(nextPrice - price)
					price = nextPrice

					if i >= 3 { // i.e. we have seen at least 4 price changes (i = 0, 1, 2, 3, ...)
						idx := seq.Index()
						if !seen[idx] {
							seen[idx] = true
							profits[idx].Add(uint32(price))
						}
					}
				}
			}
		})
	}
	wg.Wait()

	// Seen tracks which sequences have been seen for a given buyer. It's
	// cleared out between buyers, to save memory.
	// var seen [k]bool

	// For each buyer, generate all its 2000 sequences and record the
	// profits.
	// for line := range strings.SplitSeq(strings.TrimSpace(input), "\n") {
	// 	secret, err := strconv.ParseUint(line, 10, 64)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	price := int8(secret % 10)
	// 	seq := new(sequence)
	// 	for i := range 2000 {
	// 		secret = evolve(secret)
	// 		nextPrice := int8(secret % 10)
	// 		seq.Add(nextPrice - price)
	// 		price = nextPrice

	// 		if i >= 3 { // i.e. we have seen at least 4 price changes (i = 0, 1, 2, 3, ...)
	// 			idx := seq.Index()
	// 			if !seen[idx] {
	// 				seen[idx] = true
	// 				profits[idx].Add(uint32(price))
	// 			}
	// 		}
	// 	}
	// 	seen = [k]bool{}
	// }
	var answer uint32 = 0
	for _, p := range profits {
		answer = max(answer, p.Load())
	}
	return fmt.Sprint(answer), nil
}
