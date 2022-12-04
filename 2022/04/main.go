package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type assignment struct {
	start, end int
}

type pair struct {
	elves []assignment
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 04 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	pairs := parseInput(in)
	return contained(pairs)
}

func part2(in input) output {
	return 0
}

func contained(pairs []pair) int {
	total := 0
	for _, p := range pairs {
		if isContained(p) {
			total++
		}
	}
	return total
}

func isContained(p pair) bool {
	left, right := p.elves[0], p.elves[1]
	outer := left
	inner := right
	if right.start <= left.start {
		outer = right
		inner = left
	}

	if outer.end < inner.end && outer.start == inner.start {
		outer, inner = inner, outer
	}

	return outer.start <= inner.start && outer.end >= inner.end
}

func parseInput(in input) []pair {
	pairs := []pair{}
	for _, line := range in {
		sp := strings.Split(line, ",")

		p := pair{}
		for _, s := range sp {
			ls := strings.Split(s, "-")
			left, _ := strconv.Atoi(ls[0])
			right, _ := strconv.Atoi(ls[1])

			p.elves = append(p.elves, assignment{start: left, end: right})
		}
		pairs = append(pairs, p)
	}
	return pairs
}
