package day19

import (
	"fmt"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/math/span"
	"go.saser.se/adventofgo/striter"
)

const (
	none     = ""
	rejected = "R"
	accepted = "A"
)

type rule interface {
	// Check evaluates p and returns the name of the next workflow p should be
	// sent to.  Check returns an empty string if p doesn't match the rule.
	// Check returns "R" if the part is immediately rejected and "A" if the part
	// is immediately accepted.
	Check(p part) string
}

// cmpRule implements a "x<1337:foo" style rule.
type cmpRule struct {
	// Which value to check. Will be one of 'x', 'm', 'a', 's'.
	XMAS byte
	// Whether the comparison is '<'. If Less is false it is assumed that the
	// comparison is '>'.
	Less bool
	// Which value to compare against.
	Value int
	// The name of the next workflow, or "R" or "A".
	Next string
}

func parseCmpRule(s string) (cmpRule, error) {
	var c cmpRule
	c.XMAS = s[0]
	if s[1] == '<' {
		c.Less = true
	}
	colon := strings.IndexByte(s, ':')
	var err error
	c.Value, err = strconv.Atoi(s[2:colon])
	if err != nil {
		return cmpRule{}, fmt.Errorf("parse compare rule from %q: %v", s, err)
	}
	c.Next = s[colon+1:]
	return c, nil
}

func (c cmpRule) Check(p part) string {
	v := map[byte]int{
		'x': p.X,
		'm': p.M,
		'a': p.A,
		's': p.S,
	}[c.XMAS]
	var match bool
	if c.Less {
		match = v < c.Value
	} else {
		match = v > c.Value
	}
	if match {
		return c.Next
	}
	return none
}

// constRule is a rule that just immediately sends the part to workflow.
type constRule struct {
	Next string
}

func (c constRule) Check(_ part) string { return c.Next }

type workflow struct {
	Name  string
	Rules []rule
}

func parseWorkflow(s string) (workflow, error) {
	var w workflow
	openBrace := strings.IndexByte(s, '{')
	w.Name = s[:openBrace]
	parts := striter.OverSplit(s[openBrace+1:len(s)-1], ",")
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		if len(part) >= 2 && (part[1] == '<' || part[1] == '>') {
			r, err := parseCmpRule(part)
			if err != nil {
				return workflow{}, fmt.Errorf("parse workflow: %v", err)
			}
			w.Rules = append(w.Rules, r)
		} else {
			w.Rules = append(w.Rules, constRule{Next: part})
		}
	}
	return w, nil
}

// Run goes through the workflows rules, one by one, and returns the next
// workflow for the first matching rule.
func (w workflow) Run(p part) string {
	for _, r := range w.Rules {
		next := r.Check(p)
		switch next {
		case none:
			continue
		case accepted, rejected:
			return next
		default:
			return next
		}
	}
	panic("unreachable")
}

type downstreamConstraint struct {
	Next       string
	Constraint constraint
}

// Constraints returns, for each rule, the next workflow and the constraint that
// workflow will be evaluated with.
func (w workflow) Constraints(c constraint) []downstreamConstraint {
	dc := make([]downstreamConstraint, len(w.Rules))
	for i, r := range w.Rules {
		switch r := r.(type) {
		case cmpRule:
			dc[i] = downstreamConstraint{
				Next:       r.Next,
				Constraint: c.And(constraintFromCmpRule(r)),
			}
			// Create rNeg which is a negation of r, and use that to construct the
			// constraint for the next rule in the workflow.
			rNeg := r
			if rNeg.Less {
				// r  == x<1337
				// r2 == !(x<1337) ---> x>1336
				rNeg.Value--
				rNeg.Less = false
			} else {
				// r  == x>1337
				// r2 == !(x>1337) ---> x<1338
				rNeg.Value++
				rNeg.Less = true
			}
			c = c.And(constraintFromCmpRule(rNeg))
		case constRule:
			dc[i] = downstreamConstraint{
				Next:       r.Next,
				Constraint: c,
			}
		}
	}
	return dc
}

type part struct {
	X, M, A, S int
}

func parsePart(s string) (part, error) {
	var p part
	xmas := map[byte]*int{
		'x': &p.X,
		'm': &p.M,
		'a': &p.A,
		's': &p.S,
	}
	s = s[1 : len(s)-1] // Strip braces.
	splits := striter.OverSplit(s, ",")
	for split, ok := splits.Next(); ok; split, ok = splits.Next() {
		v, err := strconv.Atoi(split[2:])
		if err != nil {
			return part{}, fmt.Errorf("parse part: %v", err)
		}
		*(xmas[split[0]]) = v
	}
	return p, nil
}

type system struct {
	Workflows map[string]workflow // Workflow name -> workflow
	Parts     []part
}

func parse(input string, parseParts bool) (system, error) {
	s := system{
		Workflows: make(map[string]workflow),
	}
	fragments := striter.OverSplit(input, "\n\n")

	// Parse workflows.
	fragment, _ := fragments.Next()
	lines := striter.OverLines(fragment)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		w, err := parseWorkflow(line)
		if err != nil {
			return system{}, fmt.Errorf("parse: %v", err)
		}
		s.Workflows[w.Name] = w
	}

	if !parseParts {
		return s, nil
	}
	// Parse parts.
	fragment, _ = fragments.Next()
	lines = striter.OverLines(fragment)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		p, err := parsePart(line)
		if err != nil {
			return system{}, fmt.Errorf("parse: %v", err)
		}
		s.Parts = append(s.Parts, p)
	}
	return s, nil
}

func (s system) SumAcceptedParts() int {
	sum := 0
partLoop:
	for _, p := range s.Parts {
		w := s.Workflows["in"]
		for {
			next := w.Run(p)
			switch next {
			case accepted:
				sum += p.X + p.M + p.A + p.S
				continue partLoop
			case rejected:
				continue partLoop
			default:
				w = s.Workflows[next]
			}
		}
	}
	return sum
}

func (s system) CountAcceptedCombinations() uint64 {
	type elem struct {
		Workflow   string
		Constraint constraint
	}
	queue := []elem{
		{Workflow: "in", Constraint: universalConstraint()},
	}
	var n uint64 = 0
	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]
		w := e.Workflow
		c := e.Constraint
		switch e.Workflow {
		case accepted:
			n += c.Combinations()
		case rejected:
			continue
		default:
			for _, dc := range s.Workflows[w].Constraints(c) {
				queue = append(queue, elem{Workflow: dc.Next, Constraint: dc.Constraint})
			}
		}
	}
	return n
}

// constraint models a subset of all possible combinations of part ratings.
type constraint struct {
	X, M, A, S span.Span[int]
}

// universalConstraint returns a constraint that models all possible
// combinations of part ratings.
func universalConstraint() constraint {
	s := span.Span[int]{Start: 1, End: 4000}
	return constraint{
		X: s,
		M: s,
		A: s,
		S: s,
	}
}

// constraintFromCmpRule returns a constraint that models the subset of
// combinations that match the given cmpRule. For example, a rule like "x<1337" would return:
//
//	    {
//			X: [1, 1336],
//			M: [1, 4000],
//			A: [1, 4000],
//			S: [1, 4000],
//	    }
func constraintFromCmpRule(r cmpRule) constraint {
	c := universalConstraint()
	s := map[byte]*span.Span[int]{
		'x': &c.X,
		'm': &c.M,
		'a': &c.A,
		's': &c.S,
	}[r.XMAS]
	if r.Less {
		s.End = r.Value - 1
	} else {
		s.Start = r.Value + 1
	}
	return c
}

// And constructs the conjunction of two constraints, modeling all numbers that
// match both constraints.
func (c constraint) And(other constraint) constraint {
	c2 := c
	c2.X = span.Intersection(c2.X, other.X)
	c2.M = span.Intersection(c2.M, other.M)
	c2.A = span.Intersection(c2.A, other.A)
	c2.S = span.Intersection(c2.S, other.S)
	return c2
}

// Combinations returns the number of distinct part ratings that match this
// constraint.
func (c constraint) Combinations() uint64 {
	return uint64(c.X.Len()) * uint64(c.M.Len()) * uint64(c.A.Len()) * uint64(c.S.Len())
}

func solve(input string, part int) (string, error) {
	parseParts := part == 1
	s, err := parse(input, parseParts)
	if err != nil {
		return "", err
	}
	if part == 1 {
		return fmt.Sprint(s.SumAcceptedParts()), nil
	} else {
		return fmt.Sprint(s.CountAcceptedCombinations()), nil
	}
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
