package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

var dirs = []xy{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{1, 0},
	{1, -1},
	{1, 1},
	{0, 1},
	{0, -1},
}

type xy struct {
	x, y int
}

func (pt xy) add(pt2 xy) xy {
	return xy{x: pt.x + pt2.x, y: pt.y + pt2.y}
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 04 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parseInput(in)
	total := 0
	for k := range m {
		count := countAdjacent(m, k)
		if count < 4 {
			total++
		}
	}
	return total
}

func part2(in input) output {
	m := parseInput(in)
	total := 0
	removed := -1
	for removed != 0 {
		m, removed = removeValid(m)
		total += removed
	}
	return total
}

func removeValid(m map[xy]bool) (map[xy]bool, int) {
	m2 := make(map[xy]bool)
	count := 0
	for k := range m {
		adj := countAdjacent(m, k)
		if adj >= 4 {
			m2[k] = m[k]
		} else {
			count++
		}
	}
	return m2, count
}

func countAdjacent(m map[xy]bool, pt xy) int {
	c := 0
	for _, d := range dirs {
		pt2 := pt.add(d)
		if isAt, ok := m[pt2]; ok && isAt {
			c++
		}
	}
	return c
}

func parseInput(in input) map[xy]bool {
	m := make(map[xy]bool)
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			if c == '@' {
				m[pt] = true
			}
		}
	}
	return m
}
