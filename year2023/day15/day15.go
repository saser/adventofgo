package day15

import (
	"container/list"
	"fmt"
	"strings"

	"go.saser.se/adventofgo/striter"
)

func Part1(input string) (string, error) {
	sum := 0
	hash := 0
	for _, c := range input {
		if c == ',' {
			sum += hash
			hash = 0
			continue
		}
		hash += int(c)
		hash *= 17
		hash %= 256
	}
	sum += hash // Add the hash of the last instruction!
	return fmt.Sprint(sum), nil
}

type lens struct {
	Label       string
	FocalLength uint8
}

func hash(s string) int {
	h := 0
	for _, c := range s {
		h += int(c)
		h *= 17
		h %= 256
	}
	return h
}

func Part2(input string) (string, error) {
	boxes := make([]*list.List, 256)
	steps := striter.OverSplit(input, ",")
stepLoop:
	for step, ok := steps.Next(); ok; step, ok = steps.Next() {
		opIndex := strings.IndexAny(step, "-=")
		label := step[:opIndex]
		op := step[opIndex]

		// Ensure that there is a linked list for this box, so we avoid nil
		// checks later on -- makes for more readable code.
		h := hash(label)
		box := boxes[h]
		if box == nil {
			box = list.New()
			boxes[h] = box
		}

		if op == '-' {
			for e := box.Front(); e != nil; e = e.Next() {
				if e.Value.(*lens).Label == label {
					box.Remove(e)
					break
				}
			}
			continue stepLoop
		}
		// op == '='
		focal := step[opIndex+1] - '0'
		for e := box.Front(); e != nil; e = e.Next() {
			if l := e.Value.(*lens); l.Label == label {
				l.FocalLength = focal
				continue stepLoop
			}
		}
		// At this point we know that we should "set" a lens, but we haven't
		// found one to replace, so we insert it at the back.
		box.PushBack(&lens{Label: label, FocalLength: focal})
	}

	// We're done going through the steps, so power up the focusing power.
	power := 0
	for boxNr, box := range boxes {
		if box == nil {
			continue
		}
		slot := 1
		for e := box.Front(); e != nil; e = e.Next() {
			power += (boxNr + 1) * slot * int(e.Value.(*lens).FocalLength)
			slot++
		}
	}
	return fmt.Sprint(power), nil
}
