package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type xy struct {
	x, y int
}

type lowpoint struct {
	xy
	value int
}

type basin struct {
	point lowpoint
	size  int
}

type state struct {
	loc xy
}

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
	heights := parseInput(in)
	lows := getLowpoints(heights)
	s := 0
	for _, l := range lows {
		s += l.value + 1
	}
	return s
}

func part2(in input) output {
	heights := parseInput(in)
	basins := getBasins(heights)
	sort.Slice(basins, func(i, j int) bool {
		return basins[i].size < basins[j].size
	})
	s := 1
	for _, b := range basins[len(basins)-3:] {
		s *= b.size
	}
	return s
}

func getLowpoints(heights [][]int) []lowpoint {
	lowpoints := []lowpoint{}
	for y := 0; y < len(heights); y++ {
		for x := 0; x < len(heights[y]); x++ {
			adj := getAdjacent(heights, x, y)

			v := heights[y][x]
			low := true

			for _, pt := range adj {
				if heights[pt.y][pt.x] <= v {
					low = false
					break
				}
			}

			if low {
				lowpoints = append(lowpoints, lowpoint{xy: xy{x, y}, value: v})
			}
		}
	}
	return lowpoints
}

func getBasins(heights [][]int) []basin {
	lowpoints := getLowpoints(heights)
	basins := []basin{}
	for _, lp := range lowpoints {
		sz := findBasinAt(heights, lp.xy)

		b := basin{point: lp, size: sz}
		basins = append(basins, b)
	}
	return basins
}

func findBasinAt(heights [][]int, pt xy) int {
	visited := make(map[xy]bool)

	if heights[pt.y][pt.x] == 9 {
		return 0
	}

	initial := &state{loc: pt}
	queue := []*state{initial}
	size := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := visited[cur.loc]; ok {
			continue
		}
		visited[cur.loc] = true
		size++

		adj := getAdjacent(heights, cur.loc.x, cur.loc.y)

		for _, pt := range adj {
			if heights[pt.y][pt.x] != 9 {
				st := &state{loc: pt}
				queue = append(queue, st)
			}
		}
	}
	return size
}

func getAdjacent(heights [][]int, x, y int) []xy {
	xs := []xy{{x, y + 1}, {x, y - 1}, {x + 1, y}, {x - 1, y}}

	adj := []xy{}
	for _, p := range xs {
		if p.x < 0 || p.y > len(heights)-1 || p.x > len(heights[0])-1 || p.y < 0 {
			continue
		}
		adj = append(adj, p)
	}
	return adj
}

func parseInput(in input) [][]int {
	heights := make([][]int, len(in))

	for y, line := range in {

		for _, c := range line {
			h, _ := strconv.Atoi(string(c))
			heights[y] = append(heights[y], h)
		}
	}
	return heights
}
