package day05

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/math/span"
	"go.saser.se/adventofgo/striter"
)

type mapEntry struct {
	Src, Dst span.Span[uint]
}

func parseMapEntry(s string) (*mapEntry, error) {
	var src, dst, n uint
	if _, err := fmt.Sscanf(s, "%d %d %d", &dst, &src, &n); err != nil {
		return nil, fmt.Errorf("parse map entry: %v", err)
	}
	e := &mapEntry{
		Src: span.New(src, src+n),
		Dst: span.New(dst, dst+n),
	}
	return e, nil
}

// Translate maps v according to this map entry.
func (e *mapEntry) Translate(v uint) uint {
	return v + (e.Dst.Start - e.Src.Start)
}

type rangeMap struct {
	Entries []*mapEntry
}

func parseRangeMap(s string) (*rangeMap, error) {
	r := &rangeMap{}
	lines := striter.OverLines(s)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		e, err := parseMapEntry(line)
		if err != nil {
			return nil, fmt.Errorf("parse range map: %v", err)
		}
		r.Entries = append(r.Entries, e)
	}
	return r, nil
}

// translate splits up a given span into a set of spans and translates them
// according to the map entries. The order of the returned spans is undefined.
func (r *rangeMap) translate(s span.Span[uint]) []span.Span[uint] {
	cmpFunc := func(m *mapEntry, u uint) int {
		// span.Span.Contains returns <0 if the target is before the span, but
		// the comparison function should return <0 if the span is before the
		// target. Similarly, Contains return >0 if the target is after the
		// span, but the comparison function should return >0 if the span is
		// after the target. Because of this, we invert the result of Contains.
		return -m.Src.Contains(u)
	}
	// We can ignore the "found" bool here; r.Entries is created so that every
	// single uint value is covered by a span, so we know we found one.
	idx, _ := slices.BinarySearchFunc(r.Entries, s.Start, cmpFunc)
	var spans []span.Span[uint]
	for idx < len(r.Entries) {
		entry := r.Entries[idx]
		i := span.Intersection(s, entry.Src)
		if i == span.Empty[uint]() {
			break
		}
		start := entry.Translate(i.Start)
		spans = append(spans, span.New(start, start+i.Len()))
		s.Start = i.End
		idx++
	}
	return spans
}

// Translate splits up the given set of spans -- which is assumed to be
// non-overlapping -- and translates them according to the entries in the map.
func (r *rangeMap) Translate(spans []span.Span[uint]) []span.Span[uint] {
	var newSpans []span.Span[uint]
	for _, s := range spans {
		newSpans = append(newSpans, r.translate(s)...)
	}
	return newSpans
}

type almanac struct {
	Seeds []uint

	SeedToSoil,
	SoilToFertilizer,
	FertilizerToWater,
	WaterToLight,
	LightToTemperature,
	TemperatureToHumidity,
	HumidityToLocation *rangeMap
}

// MinLocation finds the minimum location for any number in the given span.
func (a *almanac) MinLocation(s span.Span[uint]) uint {
	spans := []span.Span[uint]{s}
	for _, r := range []*rangeMap{
		a.SeedToSoil,
		a.SoilToFertilizer,
		a.FertilizerToWater,
		a.WaterToLight,
		a.LightToTemperature,
		a.TemperatureToHumidity,
		a.HumidityToLocation,
	} {
		spans = r.Translate(spans)
	}
	return slices.MinFunc(spans, func(a, b span.Span[uint]) int { return cmp.Compare(a.Start, b.Start) }).Start
}

func parse(input string) (*almanac, error) {
	a := &almanac{}
	fragments := striter.OverSplit(input, "\n\n")

	seedString, ok := fragments.Next()
	if !ok {
		return nil, fmt.Errorf("parse input: no seeds fragment")
	}
	const prefix = "seeds: "
	if !strings.HasPrefix(seedString, prefix) {
		return nil, fmt.Errorf("parse input: parse seeds from %q: missing prefix %q", seedString, prefix)
	}
	seeds := striter.OverSplit(strings.TrimPrefix(seedString, prefix), " ")
	for seed, ok := seeds.Next(); ok; seed, ok = seeds.Next() {
		i, err := strconv.Atoi(seed)
		if err != nil {
			return nil, fmt.Errorf("parse input: parse seeds from %q: parse %q: %v", seedString, seed, err)
		}
		a.Seeds = append(a.Seeds, uint(i))
	}

	maps := []**rangeMap{
		&a.SeedToSoil,
		&a.SoilToFertilizer,
		&a.FertilizerToWater,
		&a.WaterToLight,
		&a.LightToTemperature,
		&a.TemperatureToHumidity,
		&a.HumidityToLocation,
	}
	mapIndex := 0
	for fragment, ok := fragments.Next(); ok; fragment, ok = fragments.Next() {
		newline := strings.IndexRune(fragment, '\n')
		if newline == -1 {
			return nil, fmt.Errorf("parse input: parse fragment %q: first line not found", fragment)
		}
		m, err := parseRangeMap(fragment[newline+1:])
		if err != nil {
			return nil, fmt.Errorf("parse input: parse fragment %q: %v", fragment, err)
		}
		// Sort entries in ascending order based on the start of the source ranges.
		slices.SortFunc(m.Entries, func(a, b *mapEntry) int { return cmp.Compare(a.Src.Start, b.Src.Start) })

		// Now that the entries are sorted, we can figure out where the "gaps"
		// between the entries are. We will fill those gaps with new entries that
		// have Source and Destination spans equal, which is equivalent to the
		// values for which no mapping exists. In doing so we get complete coverage
		// of all possible values from the lowest source to the highest.

		i := 0
		for i < len(m.Entries)-1 {
			lo, hi := m.Entries[i], m.Entries[i+1]
			//     vvv  = gap.Len() = 3
			// ####...####
			//     ^  ^ = gap.End   (Exclusive)
			//     |    = gap.Start (inclusive)
			gap := span.New(lo.Src.End, hi.Src.Start)
			if gap.Len() > 0 {
				entry := &mapEntry{Src: gap, Dst: gap}
				m.Entries = slices.Insert(m.Entries, i+1, entry)
				i += 2 // Skip over inserted entry.
			} else {
				i++ // No gap, so check again
			}
		}

		// Finally, we ensure that every single possible value is covered by a
		// span. This means that if the lowest span starts at x>0, we insert a
		// span [0, x) in the beginning. Similarly, if the highest span ends
		// at x<MaxUint, we insert a span [x, MaxUint) at the end.
		if x := m.Entries[0].Src.Start; x > 0 {
			s := span.New(0, x)
			e := &mapEntry{Src: s, Dst: s}
			m.Entries = append([]*mapEntry{e}, m.Entries...)
		}
		if x := m.Entries[len(m.Entries)-1].Src.End; x < math.MaxUint {
			s := span.New(x, math.MaxUint)
			e := &mapEntry{Src: s, Dst: s}
			m.Entries = append(m.Entries, e)
		}
		*maps[mapIndex] = m
		mapIndex++
	}
	return a, nil
}

func solve(input string, part int) (string, error) {
	a, err := parse(input)
	if err != nil {
		return "", err
	}
	var spans []span.Span[uint]
	if part == 1 {
		spans = make([]span.Span[uint], len(a.Seeds))
		for i, s := range a.Seeds {
			spans[i] = span.New(s, s+1)
		}
	} else {
		spans = make([]span.Span[uint], 0, len(a.Seeds)/2)
		for i := 0; i < len(a.Seeds)-1; i += 2 {
			start := a.Seeds[i]
			end := start + a.Seeds[i+1]
			spans = append(spans, span.New(start, end))
		}
	}
	best, set := uint(0), false
	for _, s := range spans {
		loc := a.MinLocation(s)
		if !set {
			best = loc
			set = true
		} else {
			best = min(best, loc)
		}
	}
	return fmt.Sprint(best), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
