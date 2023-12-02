package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

const (
	right int = 1
	left  int = -1
)

const (
	clean int = iota
	weakened
	infected
	flagged
	maxstate
)

type xy struct {
	x, y int
}
type block struct {
	xy
	val   bool
	state int
}

type carrier struct {
	xy
	dir xy
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 22 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	grid := parseGrid(in)
	return solve(grid, 10000)
}

func part2(in input) output {
	grid := parseGrid(in)
	return solveStates(grid, 10000000)
}

func minmax(grid []block) (minx, maxx, miny, maxy int) {
	minx, miny = math.MaxInt32, math.MaxInt32
	maxx, maxy = math.MinInt32, math.MinInt32

	for _, b := range grid {
		if b.x < minx {
			minx = b.x
		}
		if b.x > maxx {
			maxx = b.x
		}
		if b.y < miny {
			miny = b.y
		}
		if b.y > maxy {
			maxy = b.y
		}
	}

	return minx, maxx, miny, maxy
}

func solve(grid []block, moves int) int {
	var mid xy
	_, maxx, _, maxy := minmax(grid)
	mid.x = maxx / 2
	mid.y = maxy / 2
	switched := 0

	m := make(map[xy]block)
	for _, b := range grid {
		m[b.xy] = b
	}

	c := carrier{xy: mid, dir: xy{0, -1}} // up is negative
	for i := 0; i < moves; i++ {
		cur := m[c.xy]
		if cur.val {
			c.dir = turn(c.dir, right)
		} else {
			c.dir = turn(c.dir, left)
			switched++
		}

		cur.val = !cur.val
		m[c.xy] = cur

		c.x += c.dir.x
		c.y += c.dir.y

		if _, ok := m[c.xy]; !ok {
			m[c.xy] = block{xy: c.xy, val: false}
		}

	}
	return switched
}

func solveStates(grid []block, moves int) int {
	var mid xy
	_, maxx, _, maxy := minmax(grid)
	mid.x = maxx / 2
	mid.y = maxy / 2
	switched := 0

	m := make(map[xy]block)
	for _, b := range grid {
		m[b.xy] = b
	}

	c := carrier{xy: mid, dir: xy{0, -1}} // up is negative
	for i := 0; i < moves; i++ {
		cur := m[c.xy]

		switch cur.state {
		case clean:
			c.dir = turn(c.dir, left)
		case infected:
			c.dir = turn(c.dir, right)
		case weakened:
			switched++
			break
		case flagged:
			c.dir = xy{x: -c.dir.x, y: -c.dir.y}
		}

		cur.state = (cur.state + 1) % maxstate
		m[c.xy] = cur

		c.x += c.dir.x
		c.y += c.dir.y

		if _, ok := m[c.xy]; !ok {
			m[c.xy] = block{xy: c.xy, val: false, state: 0}
		}

	}
	return switched
}

func turn(facing xy, dir int) xy {
	turned := xy{}
	if facing.y != 0 {
		if facing.y == dir {
			turned.x = -1
		} else {
			turned.x = 1
		}
		turned.y = 0
	} else if facing.x != 0 {
		if facing.x == dir {
			turned.y = 1
		} else {
			turned.y = -1
		}
		turned.x = 0
	}

	return turned
}

func parseGrid(lines []string) []block {
	grid := []block{}
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			pos := xy{x, y}
			val := lines[y][x] == '#'
			state := clean
			if val {
				state = infected
			}
			b := block{xy: pos, val: val, state: state}
			grid = append(grid, b)
		}
	}
	return grid
}
