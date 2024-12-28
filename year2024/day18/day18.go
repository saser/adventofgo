package day18

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/asciigrid"
	"go.saser.se/adventofgo/container/priorityqueue"
)

const coordmax = 70

func fallingBytes(input string, limit int) iter.Seq[asciigrid.Pos] {
	return func(yield func(asciigrid.Pos) bool) {
		fields := strings.FieldsFunc(input, func(r rune) bool { return r == ',' || r == '\n' })
		if limit > 0 {
			fields = fields[:min(len(fields), limit*2)]
		}
		for i := 0; i < len(fields)-1; i += 2 {
			x, err := strconv.Atoi(fields[i])
			if err != nil {
				panic(fmt.Errorf("parse X: %v", err))
			}
			y, err := strconv.Atoi(fields[i+1])
			if err != nil {
				panic(fmt.Errorf("parse Y: %v", err))
			}
			if !yield(asciigrid.Pos{Row: y, Col: x}) {
				return
			}
		}
	}
}

func Part1(input string) (string, error) {
	g, err := asciigrid.New(strings.Repeat(strings.Repeat(".", coordmax+1)+"\n", coordmax+1))
	if err != nil {
		return "", fmt.Errorf("construct grid: %v", err)
	}
	for p := range fallingBytes(input, 1024) {
		g.Set(p, '#')
	}

	type state struct {
		Pos  asciigrid.Pos
		Cost int
	}
	seen := make(map[asciigrid.Pos]struct{})
	start := asciigrid.Pos{Row: 0, Col: 0}
	end := asciigrid.Pos{Row: coordmax, Col: coordmax}
	pq := priorityqueue.NewFunc(func(a, b state) bool { return a.Cost < b.Cost })
	pq.Push(state{
		Pos:  start,
		Cost: 0,
	})
	for pq.Len() > 0 {
		s := pq.Pop()
		if _, ok := seen[s.Pos]; ok {
			continue
		}
		if s.Pos == end {
			return fmt.Sprint(s.Cost), nil
		}
		seen[s.Pos] = struct{}{}
		for _, n := range s.Pos.Neighbors4() {
			if !g.InBounds(n) {
				continue
			}
			if g.Get(n) != '.' {
				continue
			}
			pq.Push(state{
				Pos:  n,
				Cost: s.Cost + 1,
			})
		}
	}
	return "", errors.New("no solution found")
}

type connections struct {
	DownLeft bool
	UpRight  bool
}

func visit(connectionsByPos map[asciigrid.Pos]connections, start asciigrid.Pos, modify func(*connections)) {
	seen := make(map[asciigrid.Pos]struct{})
	q := []asciigrid.Pos{start}
	for len(q) > 0 {
		var p asciigrid.Pos
		p, q = q[0], q[1:]
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		conn := connectionsByPos[p]
		modify(&conn)
		connectionsByPos[p] = conn
		for _, n := range p.Neighbors8() {
			if _, ok := connectionsByPos[n]; !ok {
				continue
			}
			q = append(q, n)
		}
	}
}

func Part2(input string) (string, error) {
	connectionByPos := make(map[asciigrid.Pos]connections)
	for p := range fallingBytes(input, 0) {
		conn := connections{
			DownLeft: p.Col == 0 || p.Row == coordmax,
			UpRight:  p.Row == 0 || p.Col == coordmax,
		}
		for _, n := range p.Neighbors8() {
			neighborConn, ok := connectionByPos[n]
			if !ok {
				continue
			}
			conn.DownLeft = conn.DownLeft || neighborConn.DownLeft
			conn.UpRight = conn.UpRight || neighborConn.UpRight
		}
		if conn.DownLeft && conn.UpRight {
			return fmt.Sprintf("%d,%d", p.Col, p.Row), nil
		}
		for _, n := range p.Neighbors8() {
			neighborConn, ok := connectionByPos[n]
			if !ok {
				continue
			}
			if conn.DownLeft && !neighborConn.DownLeft {
				visit(connectionByPos, n, func(c *connections) { c.DownLeft = true })
			}
			if conn.UpRight && !neighborConn.UpRight {
				visit(connectionByPos, n, func(c *connections) { c.UpRight = true })
			}
		}
		connectionByPos[p] = conn
	}
	return "", errors.New("no solution found")
}
