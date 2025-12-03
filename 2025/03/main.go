package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type bank struct {
	val, pos int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 03 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	batteries := parseInput(in)
	sum := 0
	for _, b := range batteries {
		sum += findLargest(b)
	}
	return sum
}

func part2(in input) output {
	return 0
}

func findLargest(batteries []int) int {
	fm := 0
	fp := -1

	for i, x := range batteries[:len(batteries)-1] {
		if x > fm {
			fm = x
			fp = i
		}
	}

	lm := 0
	for _, x := range batteries[fp+1:] {
		if x > lm {
			lm = x
		}
	}

	return fm*10 + lm
}

func parseInput(in input) [][]int {
	list := [][]int{}
	for _, line := range in {
		ls := []int{}
		for _, c := range line {
			x, _ := strconv.Atoi(string(c))
			ls = append(ls, x)
		}
		list = append(list, ls)
	}
	return list
}
