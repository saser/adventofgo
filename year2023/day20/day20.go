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
	Name() string
	Next() []string
	Recv(mod string, sig bool) bool
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

func (m broadcastModule) Name() string                 { return "broadcaster" }
func (m broadcastModule) Next() []string               { return m.next }
func (m broadcastModule) Recv(_ string, sig bool) bool { return sig }

type flipflopModule struct {
	name string
	next []string
}

func parseFlipFlopModule(s string) (flipflopModule, error) {
	parts, err := splitModule(s)
	if err != nil {
		return flipflopModule{}, fmt.Errorf("parse flip-flop module: %v", err)
	}
	name, _ := parts.Next()
	if name[0] != '%' {
		return flipflopModule{}, fmt.Errorf("parse flip-flop module: %q doesn't being with '%%'", name)
	}
	var next []string
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		next = append(next, part)
	}
	return flipflopModule{
		name: name[1:], // Skip over '%'.
		next: next,
	}, nil
}

func (m flipflopModule) Name() string                 { return m.name }
func (m flipflopModule) Next() []string               { return m.next }
func (m flipflopModule) Recv(_ string, sig bool) bool { return !sig }

type conjunctionModule struct {
	name   string
	next   []string
	inputs map[string]bool
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
		name:   name[1:], // Skip over '&'.
		next:   next,
		inputs: make(map[string]bool),
	}, nil
}

func (m conjunctionModule) Name() string   { return m.name }
func (m conjunctionModule) Next() []string { return m.next }
func (m *conjunctionModule) Recv(mod string, sig bool) bool {
	m.inputs[mod] = sig
	out := true
	for _, v := range m.inputs {
		out = out && v
	}
	return out
}

type circuit struct {
	Modules           map[string]module
	SentHigh, SentLow int
}

func parse(input string) (*circuit, error) {
	c := &circuit{Modules: make(map[string]module)}
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
		sig := mod.Recv(s.From, s.Sig)
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
