// Package aocdata embeds all the stored data about problem inputs and answers.
// It provides convenience functions to return the stored data.
package aocdata

import (
	"embed"
	"fmt"
	"strings"
	"testing"
)

//go:embed *_input *_output
var data embed.FS

// Input returns the puzzle input for the given year and day. If no input was
// found, it returns false. The input, if any, is returned with any trailing
// newlines removed.
func Input(year int, day int) (string, bool) {
	input, err := data.ReadFile(fmt.Sprintf("year%d_day%02d_input", year, day))
	if err != nil {
		return "", false
	}
	return strings.TrimRight(string(input), "\n"), true
}

// InputT is like Input but fails the test if the input is not found.
func InputT(tb testing.TB, year int, day int) string {
	tb.Helper()
	input, ok := Input(year, day)
	if !ok {
		tb.Fatalf("No input found for year %d, day %d.", year, day)
	}
	return input
}

// Answer returns the known answer for the given year, day, and part. If no
// answer was found, it returns false. The answer, if any, is returned with any
// trailing newlines removed.
func Answer(year int, day int, part int) (string, bool) {
	answer, err := data.ReadFile(fmt.Sprintf("year%d_day%02d_part%d_output", year, day, part))
	if err != nil {
		return "", false
	}
	return strings.TrimRight(string(answer), "\n"), true
}

// AnswerT is like Answer but fails the test if the answer is not found.
func AnswerT(tb testing.TB, year int, day int, part int) string {
	tb.Helper()
	answer, ok := Answer(year, day, part)
	if !ok {
		tb.Fatalf("No answer found for year %d, day %d, part %d.", year, day, part)
	}
	return answer
}
