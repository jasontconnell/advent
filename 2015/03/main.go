package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

type xy struct {
	x, y int
}

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) int {
	return getVisited(in[0])
}

func part2(in input) int {
	return getVisitedWithRobot(in[0])
}

func getVisited(dirs string) int {
	houses := make(map[xy]int)
	x, y := 0, 0

	for _, ch := range dirs {
		dir := string(ch)

		switch dir {
		case "<":
			x--
			break
		case ">":
			x++
			break
		case "v":
			y--
			break
		case "^":
			y++
			break
		}

		houses[xy{x, y}]++
	}

	return len(houses)
}

func getVisitedWithRobot(dirs string) int {
	santa := xy{0, 0}
	robo := xy{0, 0}

	houses := make(map[xy]int)

	for i, ch := range dirs {
		dir := string(ch)

		var current *xy

		if i%2 == 0 {
			current = &santa
		} else {
			current = &robo
		}

		switch dir {
		case "<":
			current.x--
			break
		case ">":
			current.x++
			break
		case "v":
			current.y--
			break
		case "^":
			current.y++
			break
		}

		houses[*current]++
	}

	return len(houses)
}
