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

type pair struct {
	n1, n2 int
}

func main() {
	in, err := common.ReadStringsCsv(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 02 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	pairs := parseInput(in)
	return countValid(pairs, false)
}

func part2(in input) output {
	pairs := parseInput(in)
	return countValid(pairs, true)
}

func countValid(pairs []pair, p2 bool) int {
	c := 0
	for _, p := range pairs {
		for i := p.n1; i <= p.n2; i++ {
			s := strconv.Itoa(i)
			if !p2 && !isValid(s) {
				c += i
			} else if p2 && !isValidPart2(s) {
				c += i
			}
		}
	}
	return c
}

func parseInput(in input) []pair {
	pairs := []pair{}
	for _, p := range in {
		sp := strings.Split(p, "-")
		n1, _ := strconv.Atoi(sp[0])
		n2, _ := strconv.Atoi(sp[1])
		pairs = append(pairs, pair{n1, n2})
	}
	return pairs
}

func isValidPart2(s string) bool {
	if !isValid(s) {
		return false
	}

	invalid := false
	for j := 1; j < len(s)/2+1; j++ {
		ss := s[:j]
		if s == strings.Repeat(ss, len(s)/len(ss)) {
			invalid = true
			break
		}
	}
	return !invalid
}

func isValid(s string) bool {
	if len(s)%2 == 1 {
		return true
	}

	h := len(s) / 2
	invalid := true
	for i := 0; i < h; i++ {
		if s[i] != s[i+h] {
			invalid = false
			break
		}
	}
	return !invalid
}
