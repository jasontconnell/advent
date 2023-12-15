package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type result struct {
	orig string
	val  int
}

func main() {
	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	steps := parseInput(in)
	return sum(compute(steps))
}

func part2(in input) output {
	return 0
}

func sum(results []result) int {
	s := 0
	for _, r := range results {
		s += r.val
	}
	return s
}

func compute(steps []string) []result {
	results := []result{}
	for _, step := range steps {
		r := computeOne(step)
		results = append(results, r)
	}
	return results
}

func computeOne(step string) result {
	r := result{orig: step}
	for _, c := range step {
		r.val += int(c)
		r.val *= 17
		r.val = r.val % 256
	}
	return r
}

func parseInput(in input) []string {
	sp := strings.Split(in, ",")
	return sp
}
