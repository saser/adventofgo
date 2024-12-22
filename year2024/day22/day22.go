package day22

import (
	"fmt"
	"strconv"

	"go.saser.se/adventofgo/striter"
)

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
	// between the two approaches. (Maybe the compiler recognizes that the
	// operations can be optimized into bit operations?) So I'm keeping the one
	// with math operations because it maps more clearly to the problem
	// description.

	n = (n ^ (n * 64)) % 16777216
	n = (n ^ (n / 32)) % 16777216
	n = (n ^ (n * 2048)) % 16777216
	return n
}

func Part1(input string) (string, error) {
	lines := striter.OverLines(input)
	var sum uint64 = 0
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		n, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			return "", err
		}
		for range 2000 {
			n = evolve(n)
		}
		sum += n
	}
	return fmt.Sprint(sum), nil
}

// signature returns the signature produced by the given sequence of five
// numbers.
func signature(n0, n1, n2, n3, n4 uint64) uint32 {
	d0 := uint8(n0 % 10)
	d1 := uint8(n1 % 10)
	d2 := uint8(n2 % 10)
	d3 := uint8(n3 % 10)
	d4 := uint8(n4 % 10)
	return 0 |
		uint32(d1-d0)<<24 |
		uint32(d2-d1)<<16 |
		uint32(d3-d2)<<8 |
		uint32(d4-d3)
}

// nextSignature takes an existing signature, the last number included in that
// signature, and the new number, and returns the new signature incorporating
// the new number.
func nextSignature(sig uint32, n3, n4 uint64) uint32 {
	d3 := uint8(n3 % 10)
	d4 := uint8(n4 % 10)
	return sig<<8 | uint32(d4-d3)
}

func Part2(input string) (string, error) {
	// The price changes range from -18 to +18, so they all fit in a signed
	// 8-bit integer each (1 bit for sign + 6 bits for magnitude < 8 bits). This
	// means that a sequence of four changes can be represented by a 32-bit
	// integer, where the most significant 8 bits hold the "oldest" change.
	// Let's call that 32-bit integer the "signature" of the four changes. The
	// `wonByBuyer` slice holds a map for each buyer, where the map key is the
	// change signature, and the value is the amount of bananas won by selling
	// the first time a change with that signature is seen.
	//
	// Using the example in the problem description:
	//
	// The signature for `-1,-1,0,2` is:
	// -1 => 11111111
	// -1 => 11111111
	//  0 => 00000000 (0 for positive, 000 for 0 in binary)
	//  2 => 00000010 (0 for negative, 010 for 2 in binary)
	//  => 11111111 11111111 00000000 00000010
	//  => 11111111111111110000000000000010
	//  => 4294901762 when interpreted as uint32.
	//
	// The first time the sequence is seen, the amount of bananas won is 6. So
	// if the first buyer starts with 123, then wonByBuyer[0][4294901762] = 6.
	//
	// The overall idea is to first build up the wonByBuyer data structure as
	// well as the set of all seen signatures. Then, go through all the seen
	// signatures, calculate the profit for each signature, and pick the best
	// profit.
	var wonByBuyer []map[uint32]uint8
	allSignatures := make(map[uint32]struct{})
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		n0, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			return "", err
		}
		bananas := make(map[uint32]uint8)
		// Get the first 5 numbers so we can build up the first signature.
		n1 := evolve(n0)
		n2 := evolve(n1)
		n3 := evolve(n2)
		n4 := evolve(n3)
		sig := signature(n0, n1, n2, n3, n4)
		allSignatures[sig] = struct{}{}
		bananas[sig] = uint8(n4 % 10)
		// Then go through the remaining numbers and find all the signatures.
		// Only store the amount of bananas the first time we see a signature.
		for range 2000 - 4 {
			n3, n4 = n4, evolve(n4)
			sig = nextSignature(sig, n3, n4)
			allSignatures[sig] = struct{}{}
			if _, seen := bananas[sig]; !seen {
				bananas[sig] = uint8(n4 % 10)
			}
		}
		wonByBuyer = append(wonByBuyer, bananas)
	}
	var best uint64 = 0
	for sig := range allSignatures {
		var sum uint64 = 0
		for _, won := range wonByBuyer {
			sum += uint64(won[sig])
		}
		best = max(best, sum)
	}
	return fmt.Sprint(best), nil
}
