package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type sack struct {
	left  string
	right string
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
	fmt.Fprintln(w, "--2022 day 03 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	sacks := parseInput(in)
	return sumDuped(sacks)
}

func part2(in input) output {
	sacks := parseInput(in)
	return grouped(sacks)
}

func parseInput(in input) []sack {
	sacks := []sack{}
	for _, line := range in {
		ln := len(line)
		s := sack{left: line[:ln/2], right: line[ln/2:]}
		sacks = append(sacks, s)
	}
	return sacks
}

func grouped(sacks []sack) int {
	total := 0
	for i := 0; i < len(sacks); i += 3 {
		first := duplicates(sacks[i].left+sacks[i].right, sacks[i+1].left+sacks[i+1].right)
		final := duplicates(sacks[i+2].left+sacks[i+2].right, string(first))
		total += priority(final[0])
	}
	return total
}

func sumDuped(sacks []sack) int {
	total := 0
	for _, s := range sacks {
		rs := duplicates(s.left, s.right)
		for _, ch := range rs {
			total += priority(ch)
		}
	}
	return total
}

func priority(r rune) int {
	upper := unicode.IsUpper(r)
	sub := 96
	if upper {
		sub = 38
	}
	return int(r) - sub
}

func duplicates(left, right string) []rune {
	seen := make(map[rune]bool)
	dups := []rune{}
	for _, ch := range left {
		if _, ok := seen[ch]; ok {
			continue
		}
		for _, chr := range right {
			if ch == chr {
				dups = append(dups, ch)
				break
			}
		}
		seen[ch] = true
	}
	return dups
}
