package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type content int

const (
	air  content = 0
	rock content = 1
	sand content = 2
)

type xy struct {
	x, y int
}

type rockpath struct {
	list []xy
}

type block struct {
	contents content
	coord    xy
}

func getMinMax(grid map[xy]block) (xy, xy) {
	minx, miny := math.MaxInt32, math.MaxInt32
	maxx, maxy := math.MinInt32, math.MinInt32

	for k := range grid {
		if k.x < minx {
			minx = k.x
		}
		if k.y < miny {
			miny = k.y
		}
		if k.x > maxx {
			maxx = k.x
		}
		if k.y > maxy {
			maxy = k.y
		}
	}

	if miny > 0 {
		miny = 0
	}

	return xy{minx, miny}, xy{maxx, maxy}
}

func print(grid map[xy]block, spout xy) {
	min, max := getMinMax(grid)
	if min.y > 0 {
		min.y = 0
	}
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			cur := xy{x, y}
			if b, ok := grid[cur]; ok {
				s := "#"
				if b.contents == air {
					s = "."
				}
				fmt.Print(s)
			} else if cur == spout {
				fmt.Print("+")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
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
	fmt.Fprintln(w, "--2022 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	paths := parseInput(in)
	grid := getRockGrid(paths)
	spout := xy{500, 0}
	// print(grid, spout)
	return simulate(grid, spout, false)
}

func part2(in input) output {
	paths := parseInput(in)
	grid := getRockGrid(paths)
	spout := xy{500, 0}
	// print(grid, spout)
	return simulate(grid, spout, true)
}

func simulate(grid map[xy]block, spout xy, floor bool) int {
	_, max := getMinMax(grid)

	stopy := max.y + 1
	var ground *int
	if floor {
		floory := max.y + 2
		ground = &floory
	}
	grain := spout
	grains := 0
	for {
		if grain.y >= stopy && !floor {
			break
		}

		if floor && grid[spout].contents == sand {
			break
		}

		if check(grid, xy{grain.x, grain.y + 1}, ground) {
			grain.y++
			continue
		}

		if check(grid, xy{grain.x - 1, grain.y + 1}, ground) {
			grain.x--
			grain.y++
			continue
		}

		if check(grid, xy{grain.x + 1, grain.y + 1}, ground) {
			grain.x++
			grain.y++
			continue
		}

		b := block{contents: sand, coord: grain}
		grid[grain] = b
		grain = spout
		grains++
	}
	return grains
}

func check(grid map[xy]block, pt xy, floor *int) bool {
	c := grid[pt]
	if floor != nil && c.contents == air && pt.y == *floor {
		c.contents = rock
	}
	return c.contents == air
}

func getRockGrid(rocks []rockpath) map[xy]block {
	grid := make(map[xy]block)
	for _, path := range rocks {
		for i := 1; i < len(path.list); i++ {
			start := path.list[i-1]
			grid[start] = block{contents: rock, coord: start}
			cur := path.list[i]
			delta := getDelta(start, cur)

			ptr := start
			for ptr != cur {
				ptr.x += delta.x
				ptr.y += delta.y

				grid[ptr] = block{contents: rock, coord: ptr}
			}
		}
	}
	return grid
}

func getDelta(start, end xy) xy {
	delta := xy{0, 0}
	if start.x > end.x {
		delta.x = -1
	} else if start.x < end.x {
		delta.x = 1
	}

	if start.y > end.y {
		delta.y = -1
	} else if start.y < end.y {
		delta.y = 1
	}
	return delta
}

func parseInput(in input) []rockpath {
	p := []rockpath{}
	for _, line := range in {
		coords := strings.Split(line, " -> ")
		list := []xy{}
		for _, coord := range coords {
			ns := strings.Split(coord, ",")
			x, _ := strconv.Atoi(ns[0])
			y, _ := strconv.Atoi(ns[1])

			pt := xy{x, y}
			list = append(list, pt)
		}
		path := rockpath{list: list}
		p = append(p, path)
	}
	return p
}
