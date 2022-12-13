package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/2017/knot"
	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type Point struct {
	X, Y int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	grid := getGrid(in)
	return count1s(grid)

}

func part2(in input) output {
	grid := getGrid(in)
	return countGroups(grid)
}

func getGrid(in input) []string {
	grid := []string{}
	for i := 0; i < 128; i++ {
		h := hash(in, i)
		grid = append(grid, h)
	}
	return grid
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
	check := []Point{{X: p.X, Y: p.Y - 1}, {X: p.X, Y: p.Y + 1}, {X: p.X + 1, Y: p.Y}, {X: p.X - 1, Y: p.Y}}
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
