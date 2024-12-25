package day25

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestPart1(t *testing.T) {
	aoctest.Test(t, 2024, 25, 1, Part1)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2024, 25, 1, Part1)
}
