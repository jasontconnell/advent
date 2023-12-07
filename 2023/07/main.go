package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type deal struct {
	hand     hand
	bid      int
	strength int
}

type hand struct {
	cards []card
}

type card struct {
	val  int
	valc rune
}

type strength int

const (
	HighCard  strength = 0
	OnePair            = 1000
	TwoPair            = 2000
	ThreeKind          = 3000
	FullHouse          = 4000
	FourKind           = 5000
	FiveKind           = 6000
)

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 07 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	deals := parseInput(in)
	ordered := sortWinners(deals)
	return totalWinnings(ordered)
}

func part2(in input) output {
	return 0
}

func totalWinnings(deals []deal) int {
	w := 0
	for i := 0; i < len(deals); i++ {
		w += (i + 1) * deals[i].bid
	}
	return w
}

func sortWinners(deals []deal) []deal {
	res := []deal{}
	for _, d := range deals {
		s := handStrength(d.hand)
		res = append(res, deal{hand: d.hand, bid: d.bid, strength: int(s)})
	}

	sort.Slice(res, func(i, j int) bool {
		less := res[i].strength < res[j].strength
		if res[i].strength == res[j].strength {
			for x, c := range res[i].hand.cards {
				less = c.val < res[j].hand.cards[x].val
				if c.val != res[j].hand.cards[x].val {
					break
				}
			}
		}
		return less
	})
	return res
}

func handStrength(h hand) strength {
	m := map[rune]int{}
	for _, c := range h.cards {
		m[c.valc]++
	}

	if len(m) == len(h.cards) {
		return HighCard
	}

	st := HighCard
	for _, v := range m {
		if v == 5 {
			st = FiveKind
		} else if v == 4 {
			st = FourKind
		}

		if v == 3 {
			if st == OnePair {
				st = FullHouse
			} else {
				st = ThreeKind
			}
		}

		if v == 2 {
			if st == ThreeKind {
				st = FullHouse
			} else if st == OnePair {
				st = TwoPair
			} else {
				st = OnePair
			}
		}
	}
	return st
}

func parseInput(in input) []deal {
	deals := []deal{}
	for _, line := range in {
		s := strings.Fields(line)
		h := getHand(s[0])
		b, _ := strconv.Atoi(s[1])

		d := deal{hand: h, bid: b}
		deals = append(deals, d)
	}
	return deals
}

func getHand(str string) hand {
	h := hand{}
	for _, c := range str {
		card := card{valc: c}
		switch c {
		case '2', '3', '4', '5', '6', '7', '8', '9':
			v, _ := strconv.Atoi(string(c))
			card.val = v
		case 'T':
			card.val = 10
		case 'J':
			card.val = 11
		case 'Q':
			card.val = 12
		case 'K':
			card.val = 13
		case 'A':
			card.val = 14
		}
		h.cards = append(h.cards, card)
	}
	return h
}
