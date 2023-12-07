package day07

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.saser.se/adventofgo/striter"
	"golang.org/x/exp/maps"
)

type card int

const (
	joker card = iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
	ace
)

func parseCard(r rune, withJokers bool) (card, error) {
	cards := map[rune]card{
		'2': two,
		'3': three,
		'4': four,
		'5': five,
		'6': six,
		'7': seven,
		'8': eight,
		'9': nine,
		'T': ten,
		'J': jack,
		'Q': queen,
		'K': king,
		'A': ace,
	}
	if withJokers {
		cards['J'] = joker
	}
	c, ok := cards[r]
	if !ok {
		return 0, fmt.Errorf("invalid card %q", r)
	}
	return c, nil
}

func (c card) String() string {
	return map[card]string{
		joker: "J",
		two:   "2",
		three: "3",
		four:  "4",
		five:  "5",
		six:   "6",
		seven: "7",
		eight: "8",
		nine:  "9",
		ten:   "T",
		jack:  "J",
		queen: "Q",
		king:  "K",
		ace:   "A",
	}[c]
}

type hand [5]card

func parseHand(s string, withJokers bool) (hand, error) {
	var h hand
	for i, r := range s {
		c, err := parseCard(r, withJokers)
		if err != nil {
			return hand{}, fmt.Errorf("parse hand %q: %v", s, err)
		}
		h[i] = c
	}
	return h, nil
}

func (h hand) String() string {
	var sb strings.Builder
	for _, c := range h {
		sb.WriteString(c.String())
	}
	return sb.String()
}

// compareHands returns a negative number if a is weaker than b; zero if they
// are equal; and a positive number if a is stronger than b.
func compareHands(a, b hand) int {
	return slices.Compare(a[:], b[:])
}

// signature represents the "type" of a hand, as described in the problem
// description. A signature is a string whose characters are ordered in
// ascending order, and where each character is a "class" of card. The first
// class is 'a', the second is 'b', etc. Two different classes always represent
// two different cards.
//
// For example: a signature like "aabbc" means that:
//
//   - There are two cards of class 'a'
//   - There are two cards of class 'b'
//   - There is one card of class 'c'
//
// An example of a hand matching this would be "KJKAJ" -- there are two kings,
// two jacks, and one ace.
//
// Signatures have the interesting property that in this particular game, the
// ascending alphabetical order of signatures is the same as the descending
// order of hand "types" (from strongest to weakest):
type signature string

const (
	fiveOfAKind  signature = "aaaaa"
	fourOfAKind  signature = "aaaab"
	fullHouse    signature = "aaabb"
	threeOfAKind signature = "aaabc"
	twoPair      signature = "aabbc"
	onePair      signature = "aabcd"
	highCard     signature = "abcde"
)

func handSignature(h hand, withJokers bool) signature {
	// The best possible strategy when playing with jokers is to let each joker
	// count as the most frequent non-joker card. If we are playing with jokers,
	// we "take out" the jokers, remember their count, and then add their count
	// to the most frequent card when we've found it.
	jokerCount := 0
	freq := make(map[card]int)
	for _, c := range h {
		freq[c]++
	}
	if withJokers {
		jokerCount = freq[joker]
		delete(freq, joker)
		// Special case: a hand of only jokers is a "five of a kind" hand.
		if len(freq) == 0 {
			return fiveOfAKind
		}
	}
	// The signature of a hand can be determined by the combination of how many
	// there are of the most frequent card, and how many different cards there
	// are.
	//
	// This idea is not mine; it was taken from another AoC solution:
	// https://github.com/ViddeM/advent-of-code-2023/blob/4c536cdd0c46f480f5dd3e829d2014f3fcccf451/day07/src/solution.rs#L41-L61
	max := slices.Max(maps.Values(freq)) + jokerCount
	labels := len(freq)
	switch {
	case max == 5:
		return fiveOfAKind
	case max == 4:
		return fourOfAKind
	case max == 3 && labels == 2:
		return fullHouse
	case max == 3 && labels == 3:
		return threeOfAKind
	case max == 2 && labels == 3:
		return twoPair
	case max == 2 && labels == 4:
		return onePair
	case max == 1:
		return highCard
	default:
		panic(fmt.Sprintf("invalid signature for %s (with jokers = %v)", h, withJokers))
	}
}

// compareSignatures returns a negative number if a is weaker than b; zero if
// they are equal; and a positive number if a is stronger than b.
func compareSignatures(a, b signature) int {
	// Ascending alphabetical order of signatures results in a descending order
	// of strength. By flipping the arguments to cmp.Compare, we get descending
	// alphabetical order => ascending order of strength.
	return cmp.Compare(b, a)
}

type bid struct {
	Hand hand
	Bid  int
}

func parseBid(s string, withJokers bool) (bid, error) {
	handString, bidString, ok := strings.Cut(s, " ")
	if !ok {
		return bid{}, fmt.Errorf("invalid bid %q", s)
	}
	var b bid
	var err error
	b.Hand, err = parseHand(handString, withJokers)
	if err != nil {
		return bid{}, fmt.Errorf("parse bid %q: %v", s, err)
	}
	b.Bid, err = strconv.Atoi(bidString)
	if err != nil {
		return bid{}, fmt.Errorf("parse bid %q: parse bid value %q: %v", s, bidString, err)
	}
	return b, nil
}

func (b bid) String() string {
	return fmt.Sprintf("%s %d", b.Hand, b.Bid)
}

func parse(input string, withJokers bool) ([]bid, error) {
	var bids []bid
	lines := striter.OverLines(input)
	for line, ok := lines.Next(); ok; line, ok = lines.Next() {
		b, err := parseBid(line, withJokers)
		if err != nil {
			return nil, fmt.Errorf("parse input: parse line %q: %v", line, err)
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func solve(input string, part int) (string, error) {
	withJokers := part == 2
	bids, err := parse(input, withJokers)
	if err != nil {
		return "", err
	}
	// Sort the bids in _ascending_ order. This means that the weakest bid will
	// be the first element. This corresponds well to the ranks in the problem
	// description: the weakest bid will have the lowest rank, i.e. rank 1.
	slices.SortFunc(bids, func(a, b bid) int {
		if n := compareSignatures(handSignature(a.Hand, withJokers), handSignature(b.Hand, withJokers)); n != 0 {
			return n
		}
		return compareHands(a.Hand, b.Hand)
	})
	sum := 0
	for i, b := range bids {
		rank := i + 1
		sum += b.Bid * rank
	}
	return fmt.Sprint(sum), nil
}

func Part1(input string) (string, error) {
	return solve(input, 1)
}

func Part2(input string) (string, error) {
	return solve(input, 2)
}
