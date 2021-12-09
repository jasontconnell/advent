package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	points := parseInput(in)
	return countTriangles(points)
}

func part2(in input) output {
	points := parseInput(in)
	return countVerticalTriangles(points)
}

func countTriangles(points [][]int) output {
	count := 0
	for _, trip := range points {
		if isTriangle(trip...) {
			count++
		}
	}
	return count
}

func countVerticalTriangles(points [][]int) output {
	all := []int{}
	for y := 0; y < len(points); y += 3 {
		for col := 0; col < len(points[y]); col++ {
			trip := []int{points[y][col], points[y+1][col], points[y+2][col]}
			all = append(all, trip...)
		}
	}

	count := 0
	for i := 0; i < len(all); i += 3 {
		if isTriangle(all[i : i+3]...) {
			count++
		}
	}
	return count
}

func parseInput(in input) [][]int {
	ret := make([][]int, len(in))
	for y, line := range in {
		flds := strings.Fields(line)
		li := []int{}
		for _, s := range flds {
			i, _ := strconv.Atoi(s)
			li = append(li, i)
		}
		ret[y] = li
	}

	return ret
}

func isTriangle(x ...int) bool {
	a, b, c := x[0], x[1], x[2]
	return a+b > c && a+c > b && b+c > a
}
