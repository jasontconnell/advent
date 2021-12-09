package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = string

var numpad [][]int = [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
var numpad2 [][]int = [][]int{{0, 0, 1, 0, 0}, {0, 2, 3, 4, 0}, {5, 6, 7, 8, 9}, {0, 10, 11, 12, 0}, {0, 0, 13, 0, 0}}

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
	vs := []int{}
	x, y := 1, 1
	numpad := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	for _, line := range in {
		x, y = processLine(line, x, y, numpad)
		vs = append(vs, numpad[y][x])
	}

	val := ""
	for _, v := range vs {
		val += strconv.Itoa(v)
	}
	return val
}

func part2(in input) output {
	vs := []int{}
	x, y := 1, 1
	numpad := [][]int{{0, 0, 1, 0, 0}, {0, 2, 3, 4, 0}, {5, 6, 7, 8, 9}, {0, 10, 11, 12, 0}, {0, 0, 13, 0, 0}}
	for _, line := range in {
		x, y = processLine(line, x, y, numpad)
		vs = append(vs, numpad[y][x])
	}
	s := ""
	for _, v := range vs {
		r := strconv.Itoa(v)
		switch v {
		case 10:
			r = "A"
		case 11:
			r = "B"
		case 12:
			r = "C"
		case 13:
			r = "D"
		}
		s += r
	}
	return s
}

func processLine(line string, startx, starty int, numpad [][]int) (endx, endy int) {
	endx, endy = startx, starty
	for _, k := range line {
		switch k {
		case 'U':
			if endy > 0 && numpad[endy-1][endx] != 0 {
				endy--
			}
		case 'D':
			if endy < len(numpad)-1 && numpad[endy+1][endx] != 0 {
				endy++
			}
		case 'L':
			if endx > 0 && numpad[endy][endx-1] != 0 {
				endx--
			}
		case 'R':
			if endx < len(numpad[0])-1 && numpad[endy][endx+1] != 0 {
				endx++
			}
		}
	}

	return endx, endy
}
