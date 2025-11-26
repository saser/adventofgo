// Package aoctest provides some convenience functions to run test and benchmarks on real puzzle inputs. It is intended to be used in unit tests for specific solutions, like so:
//
//	package day01
//
//	func TestPart1(t *testing.T) {
//		aoctest.Test(t, 2015, 1, 1, Part1)
//	}
//
//	func TestPart2(t *testing.T) {
//		aoctest.Test(t, 2015, 1, 2, Part2)
//	}
//
//	func BenchmarkPart1(b *testing.B) {
//		aoctest.Benchmark(b, 2015, 1, 1, Part1)
//	}
//
//	func BenchmarkPart2(b *testing.B) {
//		aoctest.Benchmark(b, 2015, 1, 2, Part2)
//	}
package aoctest

import (
	"testing"

	"go.saser.se/adventofgo/aocdata"
)

// SolveFunc is the canonical form of a solver function.
type SolveFunc func(input string) (string, error)

// Test tests the given solver function against the real input for the specified
// puzzle.
func Test(t *testing.T, year int, day int, part int, fn SolveFunc) {
	t.Helper()
	input := aocdata.InputT(t, year, day)
	want := aocdata.AnswerT(t, year, day, part)
	got, err := fn(input)
	if err != nil {
		t.Fatalf("Part%d(<real input>) err = %v", part, err)
	}
	if got != want {
		t.Fatalf("Part%d(<real input>) = %q; want %q", part, got, want)
	}
}

// Benchmark benchmarks the given solver function against the real input for the
// specified puzzle.
func Benchmark(b *testing.B, year int, day int, part int, solve SolveFunc) {
	b.Helper()
	input := aocdata.InputT(b, year, day)
	want := aocdata.AnswerT(b, year, day, part)
	got, err := solve(input)
	if err != nil {
		b.Fatalf("Part%d(<real input>) err = %v", part, err)
	}
	if got != want {
		b.Fatalf("Part%d(<real input>) = %q; want %q", part, got, want)
	}
	b.ResetTimer()
	for b.Loop() {
		solve(input)
	}
}
