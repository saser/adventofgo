package day20

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestPart1(t *testing.T) {
	aoctest.Test(t, 2024, 20, 1, Part1)
}

func TestPart2(t *testing.T) {
	aoctest.Test(t, 2024, 20, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2024, 20, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2024, 20, 2, Part2)
}
