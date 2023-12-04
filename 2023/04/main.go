package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type card struct {
	id      int
	winners []int
	mine    []int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 04 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	list := parseInput(in)
	return getWorth(list)
}

func part2(in input) output {
	list := parseInput(in)
	return playCards(list)
}

func getWorth(cards []card) int {
	totalWorth := 0
	for _, c := range cards {
		worth := 0
		matches := getMatches(c)
		for i := 0; i < matches; i++ {
			if worth == 0 {
				worth = 1
				continue
			}
			worth *= 2
		}
		totalWorth += worth
	}
	return totalWorth
}

func playCards(cards []card) int {
	counts := make(map[int]int)
	for _, c := range cards {
		counts[c.id] = 1
	}

	for _, card := range cards {
		m := getMatches(card)
		v := counts[card.id]

		for midx := 1; midx <= m; midx++ {
			counts[card.id+midx] += v
		}
	}

	total := 0
	for _, v := range counts {
		total += v
	}
	return total
}

func getMatches(c card) int {
	matches := 0
	wm := make(map[int]int)
	for _, w := range c.winners {
		wm[w] = w
	}
	for _, m := range c.mine {
		if _, ok := wm[m]; ok {
			matches++
		}
	}
	return matches
}

func parseInput(in input) []card {
	cards := []card{}
	for _, line := range in {
		parts := strings.Split(line, ":")
		idstr := strings.Fields(parts[0])[1]
		id, _ := strconv.Atoi(idstr)

		sp := strings.Split(parts[1], "|")
		w, m := []int{}, []int{}
		for _, wn := range strings.Fields(sp[0]) {
			n, _ := strconv.Atoi(wn)
			w = append(w, n)
		}

		for _, mn := range strings.Fields(sp[1]) {
			n, _ := strconv.Atoi(mn)
			m = append(m, n)
		}
		cards = append(cards, card{id: id, winners: w, mine: m})
	}
	return cards
}
