package day07

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestPart1(t *testing.T) {
	const example = `123 -> x
456 -> y
x AND y -> d
x OR y -> e
x LSHIFT 2 -> f
y RSHIFT 2 -> g
NOT x -> a
NOT y -> i`
	got, err := Part1(example)
	if err != nil {
		t.Fatalf("Part1(<example>) err = %v", err)
	}
	if want := "65412"; got != want {
		t.Fatalf("Part1(<example>) = %q; want %q", got, want)
	}
	aoctest.Test(t, 2015, 7, 1, Part1)
}

func TestPart2(t *testing.T) {
	aoctest.Test(t, 2015, 7, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2015, 7, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2015, 7, 2, Part2)
}
