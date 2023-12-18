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

type cachekey struct {
	ridx int
	widx int
}

type cache struct {
	lastseen int
	height   int
}

func print(m map[xy]bool) {
	keys := []xy{}
	for k := range m {
		keys = append(keys, k)
	}
	min, max := minmax(keys)
	fmt.Println("the grid - height:", getHeight(m))
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
	w := parseInput(in)
	pattern := getRockPattern()
	fmt.Println(lcm(len(w), len(pattern)))

	cycles := 1_000_000_000_000
	fmt.Println(cycles/(len(w)*5), cycles%(len(w)*5))
	return simulate(cycles, 7, w, pattern)
}

func simulate(cycles, width int, w []wind, rp rockpattern) int {
	nextwind := 0
	cycle := 0
	cycledelta := 1
	heightdelta := 0
	grid := getGrid(width)

	height := 0
	mem := make(map[cachekey]cache)

	for cycle < cycles {
		grid, heightdelta, nextwind, cycledelta = animate(cycle, cycles, rp, grid, width, height, w, nextwind, mem)
		height += heightdelta
		cycle += cycledelta
	}

	return height
}

func animate(cycle, cycles int, rp rockpattern, grid map[xy]bool, width, height int, winds []wind, widx int, mem map[cachekey]cache) (map[xy]bool, int, int, int) {
	const bottompad int = 3
	const leftpad int = 2
	fall := true
	startcache := int(lcm(int64(len(rp)), int64(len(winds))))

	ridx := cycle % len(rp)
	r := rp[ridx]
	rockpt := xy{leftpad, getHeight(grid) + bottompad + 1}
	cycledelta := 1
	heightdelta := 0

	ckey := cachekey{ridx: ridx, widx: widx % len(winds)}
	if state, ok := mem[ckey]; ok && state.lastseen > 0 && cycle >= startcache {
		fmt.Println(cycle, state.lastseen, cycles)
		if (cycle - state.lastseen) < cycles-cycle {
			cycledelta = cycle - state.lastseen
			heightdelta = height - state.height
			fall = false
		}
	}

	for fall {
		dir := winds[widx%len(winds)]
		widx++

		// wind blow first
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

		// fall down
		if !touching(r, rockpt, grid, width, false, true, false) {
			rockpt.y--
		} else {
			// commit to grid and do next rock
			// var trimmed bool
			fall = false
			prev := getHeight(grid)
			grid = commitToGrid(r, rockpt, grid)
			heightdelta = getHeight(grid) - prev
			if heightdelta > 0 {
				grid, _ = trimGrid(grid, 50, width)
			}
			if heightdelta < 0 {
				heightdelta = 0
			}

			if heightdelta > 0 && cycle >= startcache {
				mem[ckey] = cache{
					lastseen: cycle,
					height:   height + heightdelta,
				}
			}
		}
	}
	return grid, heightdelta, widx, cycledelta
}

func trimGrid(grid map[xy]bool, maxheight, width int) (map[xy]bool, bool) {
	_, max := minmax(keys(grid))
	ngrid := make(map[xy]bool)
	list := []xy{}
	for k, v := range grid {
		if !v {
			continue
		}
		if k.y < max.y-maxheight {
			continue
		}
		np := xy{k.x, k.y}
		list = append(list, np)
	}

	for _, k := range list {
		ngrid[k] = true
	}
	return ngrid, len(grid) > len(ngrid)
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
			if _, ok := grid[check]; ok || xf.y <= 0 {
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

func keys[K comparable, V any](m map[K]V) []K {
	ks := []K{}
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func gcd[N int | int64](a, b N) N {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm[N int | int64](a, b N, integers ...N) N {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func parseInput(in input) []wind {
	w := []wind{}
	for _, c := range in {
		w = append(w, wind(c))
	}
	return w
}
