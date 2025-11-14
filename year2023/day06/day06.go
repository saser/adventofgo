package day06

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type game struct {
	Time   int
	Record int
}

func parse(input string, joinNumbers bool) ([]game, error) {
	var times, records []int
	for line := range strings.SplitSeq(input, "\n") {
		fields := strings.Fields(line)
		var acc *[]int
		if fields[0] == "Time:" {
			acc = &times
		} else {
			acc = &records
		}
		if joinNumbers {
			joined := strings.Join(fields[1:], "")
			i, err := strconv.Atoi(joined)
			if err != nil {
				return nil, fmt.Errorf("parse input: parse line %q with joined numbers: parse %q: %v", line, joined, err)
			}
			*acc = append(*acc, i)
		} else {
			for _, f := range fields[1:] {
				i, err := strconv.Atoi(f)
				if err != nil {
					return nil, fmt.Errorf("parse input: parse line %q: parse %q: %v", line, f, err)
				}
				*acc = append(*acc, i)
			}
		}
	}
	if a, b := len(times), len(records); a != b {
		return nil, fmt.Errorf("INTERNAL BUG: parse input: found %d times and %d records; they should be the same number", a, b)
	}
	games := make([]game, len(times))
	for i := range times {
		games[i] = game{Time: times[i], Record: records[i]}
	}
	return games, nil
}

func winningDurations(g game) int {
	// Let:
	// - x be the amount of time the button is held for;
	// - t be the amount of time the game goes on for
	// - r be the record.
	// Then the total distance travelled will be x(t - x), and to beat the
	// record it has to hold that x(t - x) > r. This is a quadratic equation that can be rewritten as:
	//            x(t - x) > r
	//     => x(t - x) - r > 0
	//     => xt - x^2 - r > 0
	//     => x^2 - tx + r < 0
	// While this is an inequality, we can use the well known formula for the
	// corresponding equality equation to solve it:
	//        x^2 + (-t)x + r = 0
	//     => x = t/2 +- sqrt( (-t/2)^2 - r )
	// This will give us solutions x1 and x2 such that:
	//     0 = f(x1) > f(x) < f(x2) = 0
	//     for values x1 < x < x2
	// This means that the values we are looking for, for this AoC problem, are
	// the integers v1 and v2 such that v1 is the smallest integer greater than
	// x1, and v2 is the smallest integer less than x2.
	t := float64(g.Time)
	r := float64(g.Record)
	x1 := t/2 - math.Sqrt(math.Pow(-t/2, 2)-r)
	x2 := t/2 + math.Sqrt(math.Pow(-t/2, 2)-r)
	// It's possible that x1 and x2 are integer solutions to the equation, and
	// we need v1 > x1 and v2 < x2. To achieve that, we add or subtract 1,
	// respectively, to ensure that if x1 is an integer, then v1 will be the
	// smallest integer greater than x1, and vice versa for x2.
	v1 := int(math.Floor(x1 + 1))
	v2 := int(math.Ceil(x2 - 1))
	return v2 - v1 + 1
}

func solve(input string, part int) (string, error) {
	joinNumbers := part == 2
	games, err := parse(input, joinNumbers)
	if err != nil {
		return "", fmt.Errorf("part %d: %v", part, err)
	}
	product := 1
	for _, g := range games {
		product *= winningDurations(g)
	}
	return fmt.Sprint(product), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
