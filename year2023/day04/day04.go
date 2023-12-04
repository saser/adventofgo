package day04

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"

	"go.saser.se/adventofgo/striter"
)

type uint128 struct {
	hi, lo uint64
}

func (u *uint128) SetBit(i int) {
	if i >= 64 {
		u.hi |= 1 << (i - 64)
	} else {
		u.lo |= 1 << i
	}
}

func (u *uint128) Bit(i int) uint64 {
	if i >= 64 {
		return (u.hi >> (i - 64)) & 1
	} else {
		return (u.lo >> i) & 1
	}
}

func (u *uint128) And(other *uint128) *uint128 {
	return &uint128{
		hi: u.hi & other.hi,
		lo: u.lo & other.lo,
	}
}

func (u *uint128) OnesCount() int {
	return bits.OnesCount64(u.hi) + bits.OnesCount64(u.lo)
}

type scratchcard struct {
	ID int

	winning *uint128
	card    *uint128
}

func (s scratchcard) CountMatchingNumbers() int {
	return s.winning.And(s.card).OnesCount()
}

func parseLine(line string) (scratchcard, error) {
	var s scratchcard
	id, rest, ok := strings.Cut(line, ": ")
	if !ok {
		return scratchcard{}, fmt.Errorf("invalid line %q: couldn't find card ID", line)
	}
	var err error
	idx := strings.IndexFunc(id, unicode.IsDigit)
	if idx == -1 {
		return scratchcard{}, fmt.Errorf("invalid line %q: parse ID from %q: no digits found", line, id)
	}
	s.ID, err = strconv.Atoi(id[idx:])
	if err != nil {
		return scratchcard{}, fmt.Errorf("invalid line %q: parse ID from %q: %v", line, id, err)
	}
	parts := striter.OverSplit(rest, " ")
	inWinning := true
	s.winning = new(uint128)
	s.card = new(uint128)
	for part, ok := parts.Next(); ok; part, ok = parts.Next() {
		if part == "" {
			continue
		}
		if part == "|" {
			inWinning = false
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return scratchcard{}, fmt.Errorf("invalid line %q: %v", line, err)
		}
		if inWinning {
			s.winning.SetBit(n)
		} else {
			s.card.SetBit(n)
		}
	}
	return s, nil
}

func Part1(input string) (string, error) {
	sum := 0
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		card, err := parseLine(line)
		if err != nil {
			return "", fmt.Errorf("parse input: %v", err)
		}
		matches := card.CountMatchingNumbers()
		if matches != 0 {
			sum += (1 << (matches - 1))
		}
	}
	return fmt.Sprint(sum), nil
}

func Part2(input string) (string, error) {
	cardCount := make(map[int]int) // card ID -> how many cards of it we have
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		card, err := parseLine(line)
		if err != nil {
			return "", fmt.Errorf("parse input: %v", err)
		}
		cardCount[card.ID]++
		for i := 1; i <= card.CountMatchingNumbers(); i++ {
			// For each instance of this card we have, add a corresponding
			// number of instances of the next card we won.
			cardCount[card.ID+i] += cardCount[card.ID]
		}
	}
	sum := 0
	for _, n := range cardCount {
		sum += n
	}
	return fmt.Sprint(sum), nil
}
