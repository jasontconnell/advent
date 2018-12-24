package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var input = "17_test.txt"

const (
	XAxis int = iota
	YAxis
)

const (
	Sand int = iota
	Clay
	Water
	Spring
)

type spec struct {
	axis       int
	origin     int
	start, end int
}

type xy struct {
	x, y int
}

type block struct {
	xy
	contents int
	wet bool
}

func (b block) canHold() bool {
	return b.contents == Clay || b.contents == Water
}

func (b block) isWater() bool {
	return b.contents == Water || (b.contents == Sand && b.wet)
}

func (b block) String() string {
	str := ""
	switch b.contents {
	case Clay:
		str = "#"
	case Sand:
		str = "."
		if b.wet {
			str = "|"
		}
	case Water:
		str = "~"
	case Spring:
		str = "+"
	}
	return str
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	specs := []spec{}
	for scanner.Scan() {
		var txt = scanner.Text()
		spec := getSpec(txt)
		if spec != nil {
			specs = append(specs, *spec)
		}
	}

	blocks, xshift := makeGrid(specs, 500, 0)
	water := solve(blocks, 500-xshift, 0)

	fmt.Println("Part 1:", water)

	fmt.Println("Time", time.Since(startTime))
}

func print(blocks [][]block) {
	for y := 0; y < len(blocks); y++ {
		line := ""
		for x := 0; x < len(blocks[y]); x++ {
			s := blocks[y][x].String()
			line += s
		}
		fmt.Println(line)
	}
}

func solve(blocks [][]block, startx, starty int) int {
	blocks = dropWater(blocks, startx, starty)
	return countWater(blocks)
}

func countWater(blocks [][]block) int {
	count := 0
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			if blocks[y][x].isWater(){
				count++
			}
		}
	}
	return count
}

func dropWater(blocks [][]block, x, starty int) [][]block {
	filled := false
	y := starty
	for !filled {
		c := blocks[y][x].contents
		switch c {
		case Sand:
			blocks[y][x].wet = true
		case Water, Clay:
			blocks, filled = fillArea(blocks, x, y-1)
			y-=2
		}
		y++
		if y == len(blocks) {
			break
		}
	}
	return blocks
}

func fillArea(blocks [][]block, x, y int) ([][]block, bool) {
	res, left, right := reservoirCheck(blocks, x, y)

	for fx := left; fx < right; fx++ {
		blocks[y][fx].contents = Water
		if !res {
			blocks[y][fx].contents = Sand
			blocks[y][fx].wet = true
		}
	}

	if !res && left > 0 && right < len(blocks[y]) {
		// drop water from overflow areas
		xs := getOverflowX(blocks[y], left, right)
		for _, ovx := range xs {
			blocks = dropWater(blocks, ovx, y)
		}
	}
	return blocks, !res
}

func getOverflowX(row []block, left, right int) []int {
	ret := []int{}
	if !row[left-1].canHold() {
		ret = append(ret, left)
	}

	if !row[right].canHold() {
		ret = append(ret, right)
	}
	return ret
}

func reachedBottom(blocks [][]block, y int) bool {
	return y == len(blocks) - 1
}

func reservoirCheck(blocks [][]block, x, y int) (bool, int, int) {
	reservoir := false
	leftdone := false
	left, right := 0, 0
	for mx := x-1; !leftdone; mx-- {
		if blocks[y][mx].canHold() && blocks[y+1][mx].canHold() {
			leftdone = true
			reservoir = true
			left = mx+1
		} else if mx == 0 || (!blocks[y][mx].canHold() && !blocks[y][mx-1].canHold() && !blocks[y+1][mx-1].canHold()) {
			leftdone = true
			left = mx
		}
	}

	rightdone := false
	for mx := x; !rightdone; mx++ {
		if blocks[y][mx].canHold() && blocks[y+1][mx].canHold() {
			rightdone = true
			right = mx
		} else if y == len(blocks)-1 || !blocks[y][mx].canHold() && !blocks[y][mx+1].canHold() && !blocks[y+1][mx+1].canHold() {
			rightdone = true
			right = mx + 1
			reservoir = false
		}
	}
	
	return reservoir, left, right // right is non-inclusive upper bound
}

func waterAtBottom(blocks [][]block) bool {
	row := len(blocks) - 1
	val := false
	for x := 0; x < len(blocks[row]); x++ {
		switch blocks[row][x].contents {
		case Water:
			val = true
			break
		case Sand:
			val = blocks[row][x].wet
		}
	}
	return val
}

func makeGrid(specs []spec, springx, springy int) ([][]block, int) {
	minx := 10000
	maxx, maxy := 0, 0

	for _, s := range specs {
		switch s.axis {
		case XAxis:
			if s.origin > maxx {
				maxx = s.origin
			}
			if s.origin < minx {
				minx = s.origin
			}
			if s.end > maxy {
				maxy = s.end
			}
		case YAxis:
			if s.origin > maxy {
				maxy = s.origin
			}
			if s.end > maxx {
				maxx = s.end
			}
			if s.start < minx {
				minx = s.start
			}
		}
	}

	ybound := maxy+1
	xbound := maxx - minx + 3
	blocks := make([][]block, ybound)
	for y := 0; y < ybound; y++ {
		blocks[y] = make([]block, xbound)
		for x := 0; x < xbound; x++ {
			blocks[y][x].contents = Sand
			blocks[y][x].x = x
			blocks[y][x].y = y
		}
	}

	blocks[springy][springx-minx+1].contents = Spring
	blocks[springy][springx-minx+1].x = springx
	blocks[springy][springx-minx+1].y = springy

	for _, s := range specs {
		if s.axis == XAxis {
			for y := s.start; y < s.end+1; y++ {
				blocks[y][s.origin-minx+1].contents = Clay
			}
		} else if s.axis == YAxis {
			for x := s.start - minx + 1; x < s.end-minx+2; x++ {
				blocks[s.origin][x].contents = Clay
			}
		}
	}

	return blocks, minx-1
}

func getSpec(line string) *spec {
	reg := regexp.MustCompile(`^(y|x)=(\d+), (y|x)=(\d+)\.\.(\d+)$`)
	m := reg.FindAllStringSubmatch(line, -1)
	var s *spec
	if m != nil && len(m) > 0 {
		g := m[0]
		c := g[1]
		f, _ := strconv.Atoi(g[2])
		o1, _ := strconv.Atoi(g[4])
		o2, _ := strconv.Atoi(g[5])
		t := YAxis
		if c == "x" {
			t = XAxis
		}

		s = &spec{axis: t, origin: f, start: o1, end: o2}

	} else {
		fmt.Errorf("couldn't regex parse %s", line)
	}
	return s
}
