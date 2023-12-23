package day20

import (
	"errors"
	"fmt"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func splitModule(s string) (*striter.Joined, error) {
	name, rest, ok := strings.Cut(s, " -> ")
	if !ok {
		return nil, fmt.Errorf(`invalid module string %q: no "->" separator found`, s)
	}
	return striter.Join(
		striter.Of(name),
		striter.OverSplit(rest, ", "),
	), nil
}

type module interface {
	// Name is the name of the module.
	Name() string
	// Next is the ordered list of downstream modules to which pulses are sent.
	Next() []string
	// Recv tells the module to process the given in pulse from the named
	// module. The out pulse is valid only if ok is true; ok is false if no
	// pulse was sent downstream from this module.
	Recv(mod string, in bool) (out bool, ok bool)
}

type broadcastModule struct {
	next []string
}

func parseBroadcastModule(s string) (broadcastModule, error) {
	parts, err := splitModule(s)
	if err != nil {
		return broadcastModule{}, fmt.Errorf("parse broadcaster module: %v", err)
	}
	name, _ := parts.Next()
	if want := "broadcaster"; name != want {
		return broadcastModule{}, fmt.Errorf("parse broadcaster module: name isn't exactly %q", want)
	}
	var next []string
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		next = append(next, part)
	}
	return broadcastModule{
		next: next,
	}, nil
}

func (m broadcastModule) Name() string                        { return "broadcaster" }
func (m broadcastModule) Next() []string                      { return m.next }
func (m broadcastModule) Recv(_ string, in bool) (bool, bool) { return in, true }

type flipflopModule struct {
	name  string
	next  []string
	state bool
}

func parseFlipFlopModule(s string) (*flipflopModule, error) {
	parts, err := splitModule(s)
	if err != nil {
		return nil, fmt.Errorf("parse flip-flop module: %v", err)
	}
	name, _ := parts.Next()
	if name[0] != '%' {
		return nil, fmt.Errorf("parse flip-flop module: %q doesn't being with '%%'", name)
	}
	var next []string
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		next = append(next, part)
	}
	return &flipflopModule{
		name:  name[1:], // Skip over '%'.
		next:  next,
		state: false,
	}, nil
}

func (m *flipflopModule) Name() string   { return m.name }
func (m *flipflopModule) Next() []string { return m.next }
func (m *flipflopModule) Recv(_ string, in bool) (bool, bool) {
	if !in { // Low pulse -> flip and send pulse.
		m.state = !m.state
		return m.state, true
	}
	// High pulse -> nothing happens.
	return false, false
}

type conjunctionModule struct {
	name     string
	next     []string
	inputs   map[string]bool
	lastSent bool
}

func parseConjunctionModule(s string) (*conjunctionModule, error) {
	parts, err := splitModule(s)
	if err != nil {
		return nil, fmt.Errorf("parse conjunction module: %v", err)
	}
	name, _ := parts.Next()
	if name[0] != '&' {
		return nil, fmt.Errorf("parse conjunction module: %q doesn't being with '&'", name)
	}
	var next []string
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		next = append(next, part)
	}
	return &conjunctionModule{
		name:     name[1:], // Skip over '&'.
		next:     next,
		inputs:   make(map[string]bool),
		lastSent: true, // Since all inputs are considered low initially, the "last sent" pulse would have been high.
	}, nil
}

func (m *conjunctionModule) Name() string   { return m.name }
func (m *conjunctionModule) Next() []string { return m.next }
func (m *conjunctionModule) Recv(mod string, sig bool) (bool, bool) {
	prev, ok := m.inputs[mod]
	if !ok {
		panic(fmt.Errorf("conjunction module %q has no input named %q", m.Name(), mod))
	}
	if sig == prev {
		return m.lastSent, true
	}
	m.inputs[mod] = sig
	allHigh := true
	for _, high := range m.inputs {
		allHigh = allHigh && high
	}
	// All high pulses -> send low pulse.
	m.lastSent = !allHigh
	return m.lastSent, true
}

// AddInput defines mod as an input to this module and sets its value to false
// (low pulse).
func (m *conjunctionModule) AddInput(mod string) { m.inputs[mod] = false }

type circuit struct {
	Modules           map[string]module
	SentHigh, SentLow int
}

func parse(input string) (*circuit, error) {
	c := &circuit{
		Modules: make(map[string]module),
	}
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		var (
			m   module
			err error
		)
		switch line[0] {
		case '%':
			m, err = parseFlipFlopModule(line)
		case '&':
			m, err = parseConjunctionModule(line)
		default:
			m, err = parseBroadcastModule(line)
		}
		if err != nil {
			return nil, err
		}
		c.Modules[m.Name()] = m
	}
	for srcName, srcMod := range c.Modules {
		for _, dstName := range srcMod.Next() {
			dstMod := c.Modules[dstName]
			if c, ok := dstMod.(*conjunctionModule); ok {
				c.AddInput(srcName)
			}
		}
	}
	return c, nil
}

func (c *circuit) PressButton() {
	type send struct {
		From, To string
		Sig      bool
	}
	queue := []send{
		{From: "button", To: "broadcaster", Sig: false},
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		if s.Sig {
			c.SentHigh++
		} else {
			c.SentLow++
		}
		mod, ok := c.Modules[s.To]
		if !ok {
			// Sent to output module, so we just skip it for now.
			continue
		}
		sig, ok := mod.Recv(s.From, s.Sig)
		if !ok {
			// No pulse was sent downstream.
			continue
		}
		for _, dst := range mod.Next() {
			s2 := send{
				From: mod.Name(),
				To:   dst,
				Sig:  sig,
			}
			queue = append(queue, s2)
		}
	}
}

func solve(input string, part int) (string, error) {
	if part == 2 {
		return "", errors.New("unimplemented")
	}
	c, err := parse(input)
	if err != nil {
		return "", err
	}
	for i := 0; i < 1000; i++ {
		c.PressButton()
	}
	return fmt.Sprint(c.SentHigh * c.SentLow), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
