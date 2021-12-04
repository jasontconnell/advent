package main

import (
	"fmt"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []int

func main() {
	startTime := time.Now()

	vals, err := common.ReadInts(inputFilename)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	p1 := part1(vals)
	p2 := part2(vals)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}
func part1(vals input) int {
	if len(vals) < 2 {
		return -1
	}
	incs := 0
	for i := 1; i < len(vals); i++ {
		if vals[i] > vals[i-1] {
			incs++
		}
	}
	return incs
}

func part2(vals input) int {
	if len(vals) < 2 {
		return -1
	}

	incs := 0

	last := vals[2] + vals[1] + vals[0]
	for i := 3; i < len(vals); i++ {
		cur := vals[i] + vals[i-1] + vals[i-2]

		if cur > last {
			incs++
		}
		last = cur
	}
	return incs
}
