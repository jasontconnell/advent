package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

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
	return getArea(in)
}

func part2(in input) int {
	return getRibbon(in)
}

func getArea(in input) int {
	totalArea := 0
	for _, line := range in {
		sides := strings.Split(line, "x")
		l, _ := strconv.Atoi(sides[0])
		w, _ := strconv.Atoi(sides[1])
		h, _ := strconv.Atoi(sides[2])

		a1 := l * w
		a2 := w * h
		a3 := h * l

		area := 2*a1 + 2*a2 + 2*a3
		smallest := int(math.Min(float64(a1), math.Min(float64(a2), float64(a3))))
		area += smallest

		totalArea += area
	}
	return totalArea
}

func getRibbon(in input) int {
	totalRibbon := 0
	for _, line := range in {
		sides := strings.Split(line, "x")
		l, _ := strconv.Atoi(sides[0])
		w, _ := strconv.Atoi(sides[1])
		h, _ := strconv.Atoi(sides[2])

		list := []int{l, w, h}
		sort.Ints(list)

		ribbon := 2*list[0] + 2*list[1]
		totalRibbon += ribbon + (l * w * h)
	}
	return totalRibbon
}
