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
type output = int

type xy struct {
	x, y int
}

type grid [][]int
type state struct {
	pt xy
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("--2021 day 11 solution--")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	g := parseInput(in)
	flashes := run(g, 100)
	return flashes
}

func part2(in input) output {
	g := parseInput(in)
	return runToSimult(g, 100)
}

func runToSimult(g grid, flashes int) int {
	step := 1
	for runOne(g) != flashes {
		step++
	}
	return step
}

func run(g grid, steps int) int {
	flashes := 0
	for i := 0; i < steps; i++ {
		flashes += runOne(g)
	}
	return flashes
}

func runOne(g grid) int {
	inc(g)
	total := flash(g)
	clear(g)
	return total
}

func clear(g grid) {
	iter(g, func(x, y int) {
		if g[y][x] > 9 {
			g[y][x] = 0
		}
	})
}

func flash(g grid) int {
	total := 0
	f := getFlashes(g)
	queue := []*state{}
	for _, pt := range f {
		queue = append(queue, &state{pt: pt})
	}
	flashed := make(map[xy]bool)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := flashed[cur.pt]; ok {
			continue
		}

		g[cur.pt.y][cur.pt.x]++

		if g[cur.pt.y][cur.pt.x] > 9 {
			total++
			flashed[cur.pt] = true
			sur := getSurrounding(g, cur.pt)
			for _, pt := range sur {
				queue = append(queue, &state{pt: pt})
			}
			g[cur.pt.y][cur.pt.x] = 0
		}
	}

	return total
}

func getSurrounding(g grid, pt xy) []xy {
	list := []xy{}
	for _, p := range []xy{{pt.x, pt.y + 1}, {pt.x, pt.y - 1}, {pt.x + 1, pt.y + 1}, {pt.x + 1, pt.y - 1}, {pt.x - 1, pt.y + 1}, {pt.x - 1, pt.y - 1}, {pt.x + 1, pt.y}, {pt.x - 1, pt.y}} {
		if p.x < 0 || p.x == len(g[0]) || p.y < 0 || p.y == len(g) {
			continue
		}
		list = append(list, p)
	}
	return list
}

func getFlashes(g grid) []xy {
	list := []xy{}
	iter(g, func(x, y int) {
		if g[y][x] == 10 {
			list = append(list, xy{x, y})
		}
	})
	return list
}

func inc(g grid) {
	iter(g, func(x, y int) {
		g[y][x]++
	})
}

func iter(g grid, m func(x, y int)) {
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			m(x, y)
		}
	}
}

func parseInput(in input) grid {
	g := grid{}
	for _, line := range in {
		row := []int{}
		for _, c := range line {
			i, _ := strconv.Atoi(string(c))
			row = append(row, i)
		}
		g = append(g, row)
	}
	return g
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
