package main

import (
	"bufio"
	"fmt"
	"github.com/jasontconnell/advent/2019/intcode"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "11.txt"

type color int

const (
	black color = 0
	white color = 1
)

type xy struct {
	x, y int
}

type point struct {
	xy
	color   color
	painted bool
}

type robot struct {
	position point
	dir      point
	visited  map[xy]point
}

func (r *robot) doPaint(c color) {
	r.position.color = c
	if p, ok := r.visited[r.position.xy]; ok {
		p.color = c
		p.painted = true
		r.visited[r.position.xy] = p
	} else {
		r.visited[r.position.xy] = point{color: c, painted: true}
	}
}

// 0 is left
func (r *robot) doMove(turn int) {
	switch turn {
	case 0:
		if r.dir.x != 0 {
			if r.dir.x == -1 {
				r.dir.y = -1
			} else {
				r.dir.y = 1
			}
			r.dir.x = 0
		} else {
			if r.dir.y == -1 {
				r.dir.x = 1
			} else {
				r.dir.x = -1
			}
			r.dir.y = 0
		}
	case 1:
		if r.dir.x != 0 {
			if r.dir.x == -1 {
				r.dir.y = 1
			} else {
				r.dir.y = -1
			}
			r.dir.x = 0
		} else {
			if r.dir.y == -1 {
				r.dir.x = -1
			} else {
				r.dir.x = 1
			}
			r.dir.y = 0
		}
	}

	r.position.x += r.dir.x
	r.position.y += r.dir.y

	if _, ok := r.visited[r.position.xy]; !ok {
		r.visited[r.position.xy] = point{xy: r.position.xy, painted: false}
	}
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	opcodes := []int{}
	if scanner.Scan() {
		var txt = scanner.Text()
		sopcodes := strings.Split(txt, ",")
		for _, s := range sopcodes {
			i, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				continue
			}

			opcodes = append(opcodes, i)
		}
	}

	prog := make([]int, len(opcodes))
	copy(prog, opcodes)

	c := intcode.NewComputer(prog)
	p1, _ := paint(c, 0)
	fmt.Println("Part 1: ", p1)

	c2 := intcode.NewComputer(opcodes)
	_, r := paint(c2, 1)
	fmt.Println("Part 2: ")
	printMessage(r)

	fmt.Println("Time", time.Since(startTime))
}

func paint(c *intcode.Computer, input int) (int, *robot) {
	outtype := 0

	r := &robot{position: point{xy{x: 0, y: 0}, black, false}, dir: point{xy: xy{y: 1, x: 0}}, visited: make(map[xy]point)}

	c.AddInput(input)
	c.OnOutput = func(val int) {
		if outtype == 0 {
			r.doPaint(color(val))
		} else {
			r.doMove(val)
			in := int(r.visited[r.position.xy].color)
			c.AddInput(in)
		}

		outtype = (outtype + 1) % 2
	}

	c.Exec()

	painted := 0
	for _, p := range r.visited {
		if p.painted {
			painted++
		}
	}
	return painted, r
}

func printMessage(r *robot) {
	minx, miny := 10000, 10000
	maxx, maxy := -10000, -10000

	for _, p := range r.visited {
		if p.color == black {
			continue
		}
		if p.x < minx {
			minx = p.x
		}
		if p.y < miny {
			miny = p.y
		}
		if p.x > maxx {
			maxx = p.x
		}
		if p.y > maxy {
			maxy = p.y
		}
	}
	yrange := int(math.Abs(float64(miny))+math.Abs(float64(maxy))) + 1
	xrange := int(math.Abs(float64(minx))+math.Abs(float64(maxx))) + 1

	grid := make([][]rune, yrange)
	for i := range grid {
		grid[i] = make([]rune, xrange)
	}

	for _, p := range r.visited {
		if p.x < minx || p.y < miny || p.x > maxx || p.y > maxy {
			continue
		}
		c := ' '
		if p.color == white {
			c = '#'
		}

		x, y := p.x+int(math.Abs(float64(minx))), p.y+int(math.Abs(float64(miny)))
		//x = int(math.Abs(float64(x - maxx)))
		y = int(math.Abs(float64(y - maxy)))
		grid[y][x] = c
	}

	for r := len(grid) - 1; r >= 0; r-- {
		fmt.Println(string(grid[r]))
	}
}
