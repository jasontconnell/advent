package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type wind rune

const (
	left  wind = '<'
	right wind = '>'
)

type xy struct {
	x, y int
}

type rock []xy

type rockpattern []rock

func print(m map[xy]bool) {
	keys := []xy{}
	for k := range m {
		keys = append(keys, k)
	}
	min, max := minmax(keys)
	fmt.Println("the grid")
	for y := max.y + 1; y >= min.y; y-- {
		for x := min.x; x <= max.x; x++ {
			if m[xy{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
	fmt.Fprintln(w, "--2022 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	w := parseInput(in)
	pattern := getRockPattern()
	return simulate(2022, 7, w, pattern)
}

func part2(in input) output {
	return 0
}

func simulate(cycles, width int, w []wind, rp rockpattern) int {
	nextwind := 0
	cycle := 0
	grid := getGrid(width)

	for cycle < cycles {
		grid, nextwind = animate(rp[cycle%len(rp)], grid, width, w, nextwind)
		cycle++
	}

	return getHeight(grid)
}

func animate(r rock, grid map[xy]bool, width int, winds []wind, widx int) (map[xy]bool, int) {
	const bottompad int = 3
	const leftpad int = 2
	falling := true

	rockpt := xy{leftpad, getHeight(grid) + bottompad + 1}
	itr := 0

	for falling {
		blow := itr%2 == 0
		if blow {
			dir := winds[widx%len(winds)]
			widx++

			switch dir {
			case left:
				if !touching(r, rockpt, grid, width, true, false, true) {
					rockpt.x--
				}
			case right:
				if !touching(r, rockpt, grid, width, false, false, true) {
					rockpt.x++
				}
			}
		} else {
			if !touching(r, rockpt, grid, width, false, true, false) {
				rockpt.y--
			} else {
				// commit to grid and do next rock
				falling = false
				grid = commitToGrid(r, rockpt, grid)
			}
		}

		itr++
	}
	return grid, widx
}

func commitToGrid(r rock, pt xy, grid map[xy]bool) map[xy]bool {
	for _, rp := range r {
		cp := xy{pt.x + rp.x, pt.y + rp.y}
		grid[cp] = true
	}
	return grid
}

func touching(r rock, pt xy, grid map[xy]bool, width int, left, ignorex, ignorey bool) bool {
	res := false
	for _, rp := range r {
		xf := xy{rp.x + pt.x, rp.y + pt.y}
		if !ignorex {
			if left {
				check := xy{xf.x - 1, xf.y}
				if _, ok := grid[check]; ok || pt.x == 0 {
					res = true
				}
			} else {
				check := xy{xf.x + 1, xf.y}
				if _, ok := grid[check]; ok || check.x == width {
					res = true
				}
			}
		}

		if !ignorey {
			check := xy{xf.x, xf.y - 1}
			if _, ok := grid[check]; ok {
				res = true
			}
		}
	}
	return res
}

func getHeight(g map[xy]bool) int {
	h := 0
	for k := range g {
		if k.y > h {
			h = k.y
		}
	}
	return h
}

func minmax(pts []xy) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for _, p := range pts {
		if p.x < min.x {
			min.x = p.x
		}
		if p.x > max.x {
			max.x = p.x
		}
		if p.y < min.y {
			min.y = p.y
		}
		if p.y > max.y {
			max.y = p.y
		}
	}
	return min, max
}

func getGrid(width int) map[xy]bool {
	g := make(map[xy]bool)
	for i := 0; i < width; i++ {
		g[xy{i, 0}] = true
	}
	return g
}

func getRockPattern() rockpattern {
	// ####
	r1 := rock{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}

	// .#.
	// ###
	// .#.
	r2 := rock{
		{1, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{1, 2},
	}

	// ..#
	// ..#
	// ###
	r3 := rock{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{2, 2},
	}

	// #
	// #
	// #
	// #
	r4 := rock{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	}

	// ##
	// ##
	r5 := rock{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	}

	return rockpattern{r1, r2, r3, r4, r5}
}

func parseInput(in input) []wind {
	w := []wind{}
	for _, c := range in {
		w = append(w, wind(c))
	}
	return w
}
