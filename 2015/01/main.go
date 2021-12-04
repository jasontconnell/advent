package main

import (
	"fmt"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) int {
	return findFloor(in[0])
}

func part2(in input) int {
	return findBasement(in[0])
}

func findFloor(input string) int {
	floor := 0
	for _, chr := range input {
		switch chr {
		case '(':
			floor += 1
			break
		case ')':
			floor -= 1
			break
		}
	}
	return floor
}

func findBasement(input string) int {
	floor := 0
	index := 0
	for i, chr := range input {
		switch chr {
		case '(':
			floor += 1
			break
		case ')':
			floor -= 1
			break
		}

		if floor < 0 {
			index = i
			break
		}
	}
	return index
}
