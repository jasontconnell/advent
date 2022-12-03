package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type elf struct {
	meals []int
	sum   int
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
	fmt.Fprintln(w, "--2022 day 01 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	elves := parseInput(in)
	return max(elves, math.MaxInt32)
}

func part2(in input) output {
	elves := parseInput(in)
	var top3 int
	less := math.MaxInt32
	for i := 0; i < 3; i++ {
		m := max(elves, less)
		top3 += m
		less = m
	}

	return top3
}

func max(elves []elf, less int) int {
	m := 0
	for _, elf := range elves {
		if elf.sum > m && elf.sum < less {
			m = elf.sum
		}
	}
	return m
}

func parseInput(lines []string) []elf {
	elves := []elf{{}}
	for _, line := range lines {
		if line == "" {
			elves = append(elves, elf{})
			continue
		}
		idx := len(elves) - 1
		i, _ := strconv.Atoi(line)
		elves[idx].meals = append(elves[idx].meals, i)
		elves[idx].sum += i
	}

	return elves
}

func newElf() *elf {
	e := &elf{}
	e.meals = []int{}
	return e
}
