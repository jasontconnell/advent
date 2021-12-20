package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type grid map[xy]bool
type xy struct {
	x, y int
}

type imagealgo []bool

var extend int = 15

func printGrid(g grid, min, max xy) {
	ly := -1000
	traverse(g, min, max, func(gg grid, pt xy, val bool) {
		if pt.y != ly {
			fmt.Println()
			ly = pt.y
		}
		c := "."
		if val {
			c = "#"
		}
		fmt.Print(c)
	})
	fmt.Println()
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
	fmt.Fprintln(w, "--2021 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	algo, g := parseInput(in)
	g = loopProcess(g, algo, 2)
	return countOn(g)
}

func part2(in input) output {
	algo, g := parseInput(in)
	g = loopProcess(g, algo, 50)
	return countOn(g)
}

func loopProcess(g grid, algo imagealgo, num int) grid {

	ng := grid{}

	for k, v := range g {
		ng[k] = v
	}
	omin, omax := minmax(g)
	for i := 0; i < num; i++ {
		ng = processImage(ng, algo, omin, omax, i+1)
		if (i+1)%2 == 0 {
			omin, omax = minmax(ng)
		}
	}
	return ng
}

func processImage(g grid, algo imagealgo, min, max xy, step int) grid {
	ng := grid{}

	lex := extend
	traverse(g, min, max, func(gg grid, pt xy, val bool) {
		// ignore right and bottom edge
		if step%2 == 0 && (pt.x == max.x+lex || pt.y == max.y+lex || pt.x == min.x-lex || pt.y == min.y-lex) {
			// ng[pt] = false
			return
		}
		index := getIndex(gg, pt)
		newbit := algo[index]
		ng[pt] = newbit
	})
	return ng
}

func countOn(g grid) int {
	count := 0
	min, max := minmax(g)
	traverse(g, min, max, func(gg grid, pt xy, val bool) {
		if val {
			count++
		}
	})
	return count
}

func getIndex(g grid, p xy) int {
	pts := []xy{
		{p.x - 1, p.y - 1}, {p.x, p.y - 1}, {p.x + 1, p.y - 1},
		{p.x - 1, p.y}, {p.x, p.y}, {p.x + 1, p.y},
		{p.x - 1, p.y + 1}, {p.x, p.y + 1}, {p.x + 1, p.y + 1},
	}

	ret := 0b000000000
	for i, pt := range pts {
		if v, ok := g[pt]; ok && v {
			ret = ret | 0b1<<(8-i)
		}
	}
	return ret
}

func traverse(g grid, min, max xy, f func(gg grid, pt xy, val bool)) {
	for y := min.y - extend; y <= max.y+extend; y++ {
		for x := min.x - extend; x <= max.x+extend; x++ {
			p := xy{x, y}
			v := false

			if a, ok := g[p]; ok {
				v = a
			}

			f(g, p, v)
		}
	}
}

func minmax(g grid) (xy, xy) {
	minx, miny := math.MaxInt64, math.MaxInt64
	maxx, maxy := math.MinInt64, math.MinInt64

	for k := range g {
		if k.x < minx {
			minx = k.x
		}

		if k.x > maxx {
			maxx = k.x
		}

		if k.y < miny {
			miny = k.y
		}

		if k.y > maxy {
			maxy = k.y
		}
	}

	return xy{minx, miny}, xy{maxx, maxy}
}

func parseInput(in input) (imagealgo, grid) {
	algoline := in[0]

	enhalgo := imagealgo{}
	for _, c := range algoline {
		v := false
		if c == '#' {
			v = true
		}
		enhalgo = append(enhalgo, v)
	}

	g := grid{}
	for y, line := range in[2:] {
		for x, c := range line {
			v := false
			if c == '#' {
				v = true
			}
			g[xy{x, y}] = v
		}
	}
	return enhalgo, g
}
