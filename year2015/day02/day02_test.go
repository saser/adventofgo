package day02

import (
	"testing"

	"go.saser.se/adventofgo/aoctest"
)

func TestPart1(t *testing.T) {
	for _, test := range []struct {
		input string
		want  string
	}{
		{input: "2x3x4", want: "58"},
		{input: "1x1x10", want: "43"},
	} {
		got, err := Part1(test.input)
		if err != nil {
			t.Errorf("Part1(%q) err = %v", test.input, err)
			continue
		}
		if got != test.want {
			t.Errorf("Part1(%q) = %q; want %q", test.input, got, test.want)
			continue
		}
	}
	aoctest.Test(t, 2015, 2, 1, Part1)
}

func TestPart2(t *testing.T) {
	for _, test := range []struct {
		input string
		want  string
	}{
		{input: "2x3x4", want: "34"},
		{input: "1x1x10", want: "14"},
	} {
		got, err := Part2(test.input)
		if err != nil {
			t.Errorf("Part2(%q) err = %v", test.input, err)
			continue
		}
		if got != test.want {
			t.Errorf("Part2(%q) = %q; want %q", test.input, got, test.want)
			continue
		}
	}
	aoctest.Test(t, 2015, 2, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2015, 2, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2015, 2, 2, Part2)
}
