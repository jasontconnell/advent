package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

type direction string

const (
	forward direction = "forward"
	up      direction = "up"
	down    direction = "down"
)

type instruction struct {
	dir    direction
	amount int
}

type xy struct {
	x, y int
}

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
	insts := getInstructions(in)
	x, y := 0, 0
	for _, inst := range insts {
		switch inst.dir {
		case up:
			y = y - inst.amount
		case down:
			y = y + inst.amount
		case forward:
			x = x + inst.amount
		}
	}

	return x * y
}

func part2(in input) int {
	insts := getInstructions(in)
	x, aim, depth := 0, 0, 0
	for _, inst := range insts {
		switch inst.dir {
		case up:
			aim = aim - inst.amount
		case down:
			aim = aim + inst.amount
		case forward:
			x = x + inst.amount
			depth = depth + (aim * inst.amount)
		}
	}
	return x * depth
}

func getInstructions(lines input) []instruction {
	insts := []instruction{}
	reg := regexp.MustCompile("^(forward|up|down) ([0-9]+)$")
	for _, txt := range lines {

		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			dir := groups[1]
			amount, _ := strconv.Atoi(groups[2])

			inst := instruction{dir: direction(dir), amount: amount}
			insts = append(insts, inst)
		}
	}
	return insts
}
