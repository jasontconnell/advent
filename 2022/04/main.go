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
	sections []assignment
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
	pairs := parseInput(in)
	return overlaps(pairs)
}

func contained(pairs []pair) int {
	total := 0
	for _, p := range pairs {
		if isContained(p.sections[0], p.sections[1]) {
			total++
		}
	}
	return total
}

func overlaps(pairs []pair) int {
	total := 0
	for _, p := range pairs {
		if isOverlap(p.sections[0], p.sections[1]) {
			total++
		}
	}
	return total
}

func isOverlap(left, right assignment) bool {
	left, right = getLeftRight(left, right)
	ovr := left.start <= right.start && left.end >= right.start && right.end >= left.end
	return ovr || isContained(left, right)
}

func isContained(left, right assignment) bool {
	left, right = getLeftRight(left, right)
	return left.start <= right.start && left.end >= right.end
}

func getLeftRight(left, right assignment) (assignment, assignment) {
	outer := left
	inner := right
	if right.start <= left.start {
		outer = right
		inner = left
	}

	if outer.end < inner.end && outer.start == inner.start {
		outer, inner = inner, outer
	}

	return outer, inner
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

			p.sections = append(p.sections, assignment{start: left, end: right})
		}
		pairs = append(pairs, p)
	}
	return pairs
}
