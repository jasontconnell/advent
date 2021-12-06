package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

const (
	deckSize   int64 = 10007
	deckSize2  int64 = 119315717514047
	iterations int64 = 101741582076661
)

type input = []string
type output = int64

type action int

const (
	deal action = iota
	dealnew
	cut
)

type cardAction struct {
	text   string
	action action
	num    int64
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	actions := getInput(in)
	deck := getDeck(deckSize)

	deck = shuffle(deck, actions)
	return position(deck, 2019)
}

func part2(in input) *big.Int {
	actions := getInput(in)
	cards, shuff := big.NewInt(119315717514047), big.NewInt(101741582076661)

	return bigShuffle(cards, shuff, actions)
}

func position(deck []int64, num int64) int64 {
	var pos int64
	for i, n := range deck {
		if n == num {
			pos = int64(i)
			break
		}
	}
	return pos
}

func bigShuffle(cards, shuff *big.Int, actions []cardAction) *big.Int {
	a, b := big.NewInt(1), big.NewInt(0)
	for _, act := range actions {
		switch act.action {
		case deal:
			inc := big.NewInt(act.num)
			a.Mul(a, new(big.Int).ModInverse(inc, cards))
		case dealnew:
			a.Neg(a)
			b.Add(b, a)
		case cut:
			c := big.NewInt(act.num)
			b.Add(b, new(big.Int).Mul(a, c))
		}
	}

	b.Mul(b, new(big.Int).Sub(big.NewInt(1), new(big.Int).Exp(a, shuff, cards)))
	b.Mul(b, new(big.Int).ModInverse(new(big.Int).Sub(big.NewInt(1), a), cards))
	a.Exp(a, shuff, cards)

	return new(big.Int).Mod(new(big.Int).Add(new(big.Int).Mul(a, big.NewInt(2020)), b), cards)
}

func shuffle(deck []int64, actions []cardAction) []int64 {
	cp := getCopy(deck)
	for _, act := range actions {
		switch act.action {
		case deal:
			cp = dealIncr(cp, act.num)
		case dealnew:
			cp = dealNewStack(cp)
		case cut:
			cp = cutDeck(cp, act.num)
		}
	}
	return cp
}

func getCopy(deck []int64) []int64 {
	cp := make([]int64, len(deck))
	copy(cp, deck)
	return cp
}

func dealIncr(deck []int64, incr int64) []int64 {
	cp := getCopy(deck)
	var i int64
	for i < int64(len(deck)) {
		pos := (i * incr) % int64(len(deck))
		cp[pos] = deck[i]
		i++
	}
	return cp
}

func dealNewStack(deck []int64) []int64 {
	cp := getCopy(deck)
	for i := 0; i < len(deck)/2; i++ {
		cp[i], cp[len(deck)-i-1] = cp[len(deck)-i-1], cp[i]
	}
	return cp
}

func cutDeck(deck []int64, size int64) []int64 {
	cp := make([]int64, len(deck))

	start := size
	if size < 0 {
		start = int64(len(deck)) + size
	}

	var i, pos int64
	for i = start; i < int64(len(deck)); i++ {
		cp[pos] = deck[i]
		pos++
	}
	for i = 0; i < start; i++ {
		cp[pos] = deck[i]
		pos++
	}
	return cp
}

func getDeck(n int64) []int64 {
	deck := []int64{}
	var i int64
	for i = 0; i < n; i++ {
		deck = append(deck, i)
	}
	return deck
}

func getInput(in input) []cardAction {
	actions := []cardAction{}
	for _, line := range in {
		var n int64
		act := cut
		if strings.HasPrefix(line, "deal") {
			act = deal
		}
		if strings.HasSuffix(line, "new stack") {
			act = dealnew
		} else {
			sp := strings.Fields(line)
			n, _ = strconv.ParseInt(sp[len(sp)-1], 10, 64)
		}

		ca := cardAction{action: act, num: n, text: line}
		actions = append(actions, ca)
	}
	return actions
}
