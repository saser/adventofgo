package day07

import (
	"fmt"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
)

type expression interface {
	FromString(s string) error
	isExpression()
}

type nonary struct {
	integer uint16
	wire    string
}

func (e *nonary) FromString(s string) error {
	if strings.ContainsRune(s, ' ') {
		return fmt.Errorf("parse nonary: %q contains spaces", s)
	}
	i, err := strconv.ParseUint(s, 10, 16)
	if err == nil { // if NO error: is integer
		e.integer = uint16(i)
	} else { // if error: is wire
		e.wire = s
	}
	return nil
}

func (e *nonary) isExpression()   {}
func (e *nonary) IsWire() bool    { return e.wire != "" }
func (e *nonary) IsInteger() bool { return !e.IsWire() }

type unary struct {
	e *nonary
	f func(v uint16) uint16
}

func (u *unary) FromString(s string) error {
	var ns string
	if _, err := fmt.Sscanf(s, "NOT %s", &ns); err != nil {
		return fmt.Errorf(`parse unary: %q doesn't have the form "NOT <nonary>"`, s)
	}
	var n nonary
	if err := n.FromString(ns); err != nil {
		return fmt.Errorf("parse unary from %q: %v", s, err)
	}
	u.e = &n
	u.f = func(v uint16) uint16 { return ^v }
	return nil
}

func (*unary) isExpression() {}

type binary struct {
	n1, n2 *nonary
	f      func(v1, v2 uint16) uint16
}

func (b *binary) FromString(s string) error {
	parts := strings.SplitN(s, " ", 3)
	if len(parts) != 3 {
		return fmt.Errorf("parse binary: %q doesn't consist of three space-separated parts", s)
	}
	var n1, n2 nonary
	if err := n1.FromString(parts[0]); err != nil {
		return fmt.Errorf("parse binary from %q: %v", s, err)
	}
	if err := n2.FromString(parts[2]); err != nil {
		return fmt.Errorf("parse binary from %q: %v", s, err)
	}
	ops := map[string]func(v1, v2 uint16) uint16{
		"AND":    func(v1, v2 uint16) uint16 { return v1 & v2 },
		"OR":     func(v1, v2 uint16) uint16 { return v1 | v2 },
		"LSHIFT": func(v1, v2 uint16) uint16 { return v1 << v2 },
		"RSHIFT": func(v1, v2 uint16) uint16 { return v1 >> v2 },
	}
	op, ok := ops[parts[1]]
	if !ok {
		return fmt.Errorf("parse binary from %q: unknown operation %q", s, parts[1])
	}
	b.n1 = &n1
	b.n2 = &n2
	b.f = op
	return nil
}

func (*binary) isExpression() {}

type circuit struct {
	wires map[string]expression
	env   map[string]uint16
}

func newCircuit(wires map[string]expression) *circuit {
	return &circuit{
		wires: wires,
		env:   make(map[string]uint16),
	}
}

func (c *circuit) evaluateNonary(n *nonary) uint16 {
	if n.IsInteger() {
		return n.integer
	}
	return c.Evaluate(n.wire)
}

func (c *circuit) Evaluate(wire string) (value uint16) {
	expr := c.wires[wire]
	if cached, ok := c.env[wire]; ok {
		return cached
	}
	defer func() { c.env[wire] = value }()
	switch v := expr.(type) {
	case *nonary:
		return c.evaluateNonary(v)
	case *unary:
		return v.f(c.evaluateNonary(v.e))
	case *binary:
		return v.f(c.evaluateNonary(v.n1), c.evaluateNonary(v.n2))
	default:
		panic(fmt.Errorf("expression for wire %q has invalid type %T: %#v", wire, v, v))
	}
}

func solve(input string, part int) (string, error) {
	lines := striter.OverLines(input)
	wires := make(map[string]expression)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		exprString, wire, found := strings.Cut(line, " -> ")
		if !found {
			return "", fmt.Errorf(`invalid input line %q: doesn't have form "<expr> -> <wire>"`, line)
		}
		for _, expr := range []expression{
			new(nonary),
			new(unary),
			new(binary),
		} {
			if err := expr.FromString(exprString); err == nil { // if NO error
				wires[wire] = expr
				break
			}
		}
		if _, ok := wires[wire]; !ok {
			return "", fmt.Errorf("invalid input line %q", line)
		}
	}
	c := newCircuit(wires)
	a := c.Evaluate("a")
	if part == 2 {
		wires["b"] = &nonary{integer: a}
		c = newCircuit(wires)
		a = c.Evaluate("a")
	}
	return fmt.Sprint(a), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
