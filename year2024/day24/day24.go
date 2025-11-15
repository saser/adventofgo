package day24

import (
	"errors"
	"fmt"
	"iter"
	"regexp"
	"strings"
)

var (
	initialRE = regexp.MustCompile(`^([xy]\d{2}): ([01])$`)
	gateRE    = regexp.MustCompile(`^(\w{3}) (AND|OR|XOR) (\w{3}) -> (\w{3})$`)
)

func ptr(v uint64) *uint64 { return &v }

type wire struct {
	Name  string  // e.g. "x00" or "pqr"
	Value *uint64 // 0 or 1
}

func parseWire(line string) (*wire, error) {
	matches := initialRE.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid line %q", line)
	}
	w := &wire{
		Name:  matches[1],
		Value: ptr(0),
	}
	if matches[2] == "1" {
		w.Value = ptr(1)
	}
	return w, nil
}

func (w *wire) String() string {
	if w.Value == nil {
		return fmt.Sprintf("%v(.)", w.Name)
	} else {
		return fmt.Sprintf("%v(%v)", w.Name, *w.Value)
	}
}

type gateDesc struct {
	W1  string
	Op  string
	W2  string
	Out string
}

func parseGateDesc(line string) (gateDesc, error) {
	matches := gateRE.FindStringSubmatch(line)
	if matches == nil {
		return gateDesc{}, fmt.Errorf("invalid line %q", line)
	}
	var g gateDesc
	g.W1, g.Op, g.W2, g.Out = matches[1], matches[2], matches[3], matches[4]
	return g, nil
}

type gate struct {
	W1  *wire
	Op  string // AND, OR, XOR
	W2  *wire
	Out *wire
}

func (g *gate) String() string { return fmt.Sprintf("%v %v %v -> %v", g.W1, g.W2, g.Op, g.Out) }

type system struct {
	wires map[string]*wire // wire name -> wire
	gates map[string]*gate // output wire -> gate
}

func parseSystem(input string) (*system, error) {
	fragments := strings.Split(input, "\n\n")

	wires := make(map[string]*wire)
	for line := range strings.SplitSeq(fragments[0], "\n") {
		w, err := parseWire(line)
		if err != nil {
			return nil, err
		}
		wires[w.Name] = w
	}

	gates := make(map[string]*gate)
	for line := range strings.SplitSeq(fragments[1], "\n") {
		d, err := parseGateDesc(line)
		if err != nil {
			return nil, err
		}
		for _, name := range []string{d.W1, d.W2, d.Out} {
			if _, exists := wires[name]; !exists {
				wires[name] = &wire{Name: name}
			}
		}
		gates[d.Out] = &gate{
			W1:  wires[d.W1],
			Op:  d.Op,
			W2:  wires[d.W2],
			Out: wires[d.Out],
		}
	}

	return &system{
		wires: wires,
		gates: gates,
	}, nil
}

func (s *system) evaluate(w *wire) uint64 {
	if w.Value != nil {
		return *w.Value
	}
	g, ok := s.gates[w.Name]
	if !ok {
		panic(fmt.Errorf("%v has no value and is also not the output of any gate", w))
	}
	v1 := s.evaluate(g.W1)
	v2 := s.evaluate(g.W2)
	var out uint64
	switch g.Op {
	case "AND":
		out = v1 & v2
	case "OR":
		out = v1 | v2
	case "XOR":
		out = v1 ^ v2
	}
	w.Value = ptr(out)
	return out
}

func (s *system) evaluateAll() {
	for _, w := range s.wires {
		s.evaluate(w)
	}
}

func (s *system) reset(x, y uint64) {
	for _, w := range s.wires {
		w.Value = nil
	}
	for i, v := range bitsOf(x) {
		// My puzzle input has x00 through x44
		if i > 44 {
			break
		}
		name := fmt.Sprintf("x%02d", i)
		w, ok := s.wires[name]
		if !ok {
			break
		}
		w.Value = ptr(v)
	}
	for i, v := range bitsOf(y) {
		// My puzzle input has y00 through y44
		if i > 44 {
			break
		}
		name := fmt.Sprintf("y%02d", i)
		w, ok := s.wires[name]
		if !ok {
			break
		}
		w.Value = ptr(v)
	}
}

func (s *system) z() uint64 {
	var z uint64
	for i := range 64 {
		name := fmt.Sprintf("z%02d", i)
		w, ok := s.wires[name]
		if !ok {
			break
		}
		z |= (*w.Value << i)
	}
	return z
}

func (s *system) run(x, y uint64) uint64 {
	s.reset(x, y)
	s.evaluateAll()
	return s.z()
}

func bitsOf(v uint64) iter.Seq2[int, uint64] {
	return func(yield func(int, uint64) bool) {
		for i := range 64 {
			if !yield(i, v&1) {
				return
			}
			v >>= 1
		}
	}
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	s, err := parseSystem(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	s.evaluateAll()
	return fmt.Sprint(s.z()), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
