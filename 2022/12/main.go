package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

const (
	start int = 96
	end   int = 123
)

type xy struct {
	x, y int
}

func (pt xy) dist(from xy) int {
	dx := from.x - pt.x
	dy := from.y - pt.y

	return int(math.Abs(float64(dx))) + int(math.Abs(float64(dy)))
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 12 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	grid := parseInput(in)
	return traverse(grid, false, xy{})
}

func part2(in input) output {
	grid := parseInput(in)
	starts := []xy{}
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == start+1 {
				starts = append(starts, xy{c, r})
			}
		}
	}

	vals := []int{}
	for _, sxy := range starts {
		v := traverse(grid, true, sxy)
		if v > 0 {
			vals = append(vals, v)
		}
	}

	sort.Ints(vals)
	return vals[0]
}

func getStartEnd(grid [][]int) (xy, xy) {
	sxy, exy := xy{}, xy{}
	for y, r := range grid {
		for x, c := range r {
			if c == start {
				sxy = xy{x, y}
			} else if c == end {
				exy = xy{x, y}
			}
		}
	}
	return sxy, exy
}

func traverse(grid [][]int, override bool, startoverride xy) int {
	sxy, exy := getStartEnd(grid)
	visit := make(map[xy]int)
	queue := common.NewQueue[xy, int]()

	if override {
		sxy = startoverride
	}
	queue.Enqueue(sxy)

	for queue.Any() {
		cur := queue.Dequeue()

		if cur == exy {
			break
		}

		mvs := getMoves(grid, cur)
		for _, mv := range mvs {
			if _, ok := visit[mv]; !ok {
				visit[mv] = visit[cur] + 1
				queue.Enqueue(mv)
			}
		}
	}
	return visit[exy]
}

func getMoves(grid [][]int, pt xy) []xy {
	mvs := []xy{}

	for _, mv := range []xy{
		{x: 1, y: 0},
		{x: -1, y: 0},
		{x: 0, y: -1},
		{x: 0, y: 1},
	} {
		to := xy{x: pt.x + mv.x, y: pt.y + mv.y}
		if to.x == -1 || to.x >= len(grid[0]) || to.y == -1 || to.y >= len(grid) {
			continue
		}

		cur := grid[pt.y][pt.x]
		dest := grid[to.y][to.x]

		if cur == start {
			cur = start + 1
		}

		if dest == end {
			dest = end - 1
		}

		diff := dest - cur
		if diff <= 1 {
			mvs = append(mvs, to)
		}
	}
	return mvs
}

func parseInput(in input) [][]int {
	grid := [][]int{}
	for _, line := range in {
		r := []int{}
		for _, c := range line {
			cval := int(c)
			if c == 'E' {
				cval = end
			} else if c == 'S' {
				cval = start
			}
			r = append(r, cval)
		}
		grid = append(grid, r)
	}
	return grid
}
