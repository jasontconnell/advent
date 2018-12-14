package main

import (
	"fmt"
	"github.com/jasontconnell/advent/2017/knot"
	"strconv"
	"time"
)

var input = "oundnydw"

var testInput = []string{
	"1111111011111110",
	"1001000101111110",
	"1001011101111110",
	"1111111101111110",
}

type Point struct {
	X, Y int
}

func main() {
	startTime := time.Now()

	grid := []string{}
	for i := 0; i < 128; i++ {
		h := hash(input, i)
		grid = append(grid, h)
	}

	n := count1s(grid)
	g := countGroups(grid)

	fmt.Println("Number of 1s:           :", n)
	fmt.Println("Number of groups is     :", g)

	fmt.Println("Time", time.Since(startTime))
}

func hash(key string, rnum int) string {
	k := key + "-" + strconv.Itoa(rnum)
	h := knot.KnotHashBinary(k)

	return h
}

func count1s(s []string) int {
	ones := 0
	for _, line := range s {
		for _, c := range line {
			if c == '1' {
				ones++
			}
		}
	}
	return ones
}

func countGroups(s []string) int {
	visited := make(map[Point]bool)
	g := 0

	for y := 0; y < len(s); y++ {
		for x := 0; x < len(s[y]); x++ {
			p := Point{X: x, Y: y}
			isOne := checkPoint(s, p)
			_, ok := visited[p]
			visited[p] = true

			if !ok && isOne {
				markGroup(visited, s, p)
				g++
			}
		}
	}

	return g
}

func markGroup(visited map[Point]bool, s []string, p Point) {
	check := []Point{Point{X: p.X, Y: p.Y - 1}, Point{X: p.X, Y: p.Y + 1}, Point{X: p.X + 1, Y: p.Y}, Point{X: p.X - 1, Y: p.Y}}
	for _, chk := range check {
		isOne := checkPoint(s, chk)
		_, ok := visited[chk]
		visited[chk] = true

		if !ok && isOne {
			markGroup(visited, s, chk)
		}
	}
}

func checkPoint(s []string, p Point) bool {
	b := false
	if p.X != -1 && p.Y != -1 && p.Y < len(s) && p.X < len(s[p.Y]) {
		b = s[p.Y][p.X] == '1'
	}

	return b
}
