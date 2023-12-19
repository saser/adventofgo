package day19

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestPart1(t *testing.T) {
	aoctest.Test(t, 2023, 19, 1, Part1)
}

// 115849194428000 = 0x695d3e14d660
// 167409079868000

func TestPart2(t *testing.T) {
	aoctest.Test(t, 2023, 19, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2023, 19, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2023, 19, 2, Part2)
}
