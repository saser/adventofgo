package day10

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestLookSay(t *testing.T) {
	for input, want := range map[string]string{
		"1":      "11",
		"11":     "21",
		"21":     "1211",
		"1211":   "111221",
		"111221": "312211",
	} {
		if got := lookSay([]byte(input)); string(got) != want {
			t.Errorf("lookSay(%q) = %q; want %q", input, got, want)
			continue
		}
	}
}

func TestPart1(t *testing.T) {
	aoctest.Test(t, 2015, 10, 1, Part1)
}

func TestPart2(t *testing.T) {
	aoctest.Test(t, 2015, 10, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2015, 10, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2015, 10, 2, Part2)
}
