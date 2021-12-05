package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

var reg *regexp.Regexp = regexp.MustCompile("^([0-9]+),([0-9]+) -> ([0-9]+),([0-9]+)$")

type xy struct {
	x, y int
}

type line struct {
	left, right xy
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
	lines := getInput(in)
	horver := filterHorVer(lines)
	return findOverlaps(horver)
}

func part2(in input) output {
	lines := getInput(in)
	return findOverlaps(lines)
}

func filterHorVer(lines []line) []line {
	matching := []line{}
	for _, check := range lines {
		if isHorizontal(check) || isVertical(check) {
			matching = append(matching, check)
		}
	}
	return matching
}

func isHorizontal(check line) bool {
	return check.left.y == check.right.y
}

func isVertical(check line) bool {
	return check.left.x == check.right.x
}

func findOverlaps(lines []line) output {
	m := make(map[xy]int)
	for _, l := range lines {
		pts := points(l)
		for _, p := range pts {
			m[p]++
		}
	}
	c := 0
	for _, v := range m {
		if v > 1 {
			c++
		}
	}
	return c
}

func slope(l, r int) int {
	if l == r {
		return 0
	}

	if l > r {
		return -1
	}
	return 1
}

func dist(l, r xy) int {
	diff := l.y - r.y
	if diff == 0 {
		diff = l.x - r.x
	}

	return int(math.Abs(float64(diff)))
}

func points(l line) []xy {
	pts := []xy{}

	xdelta := slope(l.left.x, l.right.x)
	ydelta := slope(l.left.y, l.right.y)

	itr := dist(l.left, l.right)

	for i := 0; i <= itr; i++ {
		x := l.left.x + i*xdelta
		y := l.left.y + i*ydelta

		p := xy{x, y}
		pts = append(pts, p)
	}
	return pts
}

func minmax(l, r int) (int, int) {
	min := int(math.Min(float64(l), float64(r)))
	max := l
	if min == l {
		max = r
	}

	return min, max
}

func getInput(in input) []line {
	lines := []line{}
	for _, row := range in {
		groups := reg.FindStringSubmatch(row)

		if len(groups) == 5 {
			sx, _ := strconv.Atoi(groups[1])
			sy, _ := strconv.Atoi(groups[2])
			ex, _ := strconv.Atoi(groups[3])
			ey, _ := strconv.Atoi(groups[4])

			lines = append(lines, line{left: xy{sx, sy}, right: xy{ex, ey}})
		}
	}

	return lines
}
