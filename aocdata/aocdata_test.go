package aocdata

import "testing"

func TestInput(t *testing.T) {
	// There should be data available for all existing problems. This test will
	// go out of date over time as new events are run. It should ideally be
	// updated after December 25th each year.
	for year := 2015; year <= 2022; year++ {
		for day := 1; day <= 25; day++ {
			input, ok := Input(year, day)
			if !ok {
				t.Errorf("Input(%d, %d) ok = false", year, day)
				continue
			}
			if input == "" {
				t.Errorf(`Input(%d, %d) input = ""; want non-empty`, year, day)
				continue
			}
		}
	}
}

func TestAnswer(t *testing.T) {
	// The testing strategy here is to specify a rather large subset of problems
	// that I know there should be an answer for, at the time of writing. It's
	// not comprehensive, but it should find egregious and/or obvious bugs.
	for year := 2015; year <= 2022; year++ {
		for day := 1; day <= 15; day++ {
			for part := 1; part <= 2; part++ {
				answer, ok := Answer(year, day, part)
				if !ok {
					t.Errorf("Answer(%d, %d, %d) ok = false", year, day, part)
					continue
				}
				if answer == "" {
					t.Errorf(`Answer(%d, %d, %d) input = ""; want non-empty`, year, day, part)
					continue
				}
			}
		}
	}
}
