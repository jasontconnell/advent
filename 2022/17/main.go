package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
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

func (w wind) String() string {
	return string(w)
}

type xy struct {
	x, y int
}

func (p xy) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

type rock []xy

type rockpattern []rock

type cachekey struct {
	keys string
	ridx int
	widx int
}

func (k cachekey) String() string {
	return fmt.Sprintf("key: [k: %d rk: %d]", len(k.keys), k.ridx)
}

type cache struct {
	lastseen int
	height   int
}

func (c cache) String() string {
	return fmt.Sprintf("cache: [seen: %d height: %d]", c.lastseen, c.height)
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

	cycles := 1_000_000_000_000
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

	mkey := getGridKey(grid, width, 30)
	// wstr := windkey(winds, widx, 1)
	ckey := cachekey{keys: mkey, ridx: ridx, widx: widx % len(winds)}
	if state, ok := mem[ckey]; ok && state.lastseen > 0 && cycle > startcache {
		if (cycle - state.lastseen) < cycles-cycle {
			cycledelta = cycle - state.lastseen
			heightdelta = height - state.height

			// state.lastseen = cycle
			// state.height = state.height + heightdelta
			// mem[ckey] = state
			// fmt.Println(state, cycledelta, heightdelta)
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
			commitToGrid(r, rockpt, grid)
			heightdelta = getHeight(grid) - prev
			if heightdelta > 0 {
				grid, _ = trimGrid(grid, 50, width)
			}
			if heightdelta < 0 {
				heightdelta = 0
			}

			if heightdelta > 0 && cycle > startcache {
				mem[ckey] = cache{
					lastseen: cycle,
					height:   height + heightdelta,
				}
			}
		}
	}
	return grid, heightdelta, widx, cycledelta
}

func getGridKey(grid map[xy]bool, width, maxheight int) string {
	_, max := minmax(keys(grid))

	s := ""
	for y := max.y; y >= max.y-maxheight; y-- {
		if y < 0 {
			continue
		}
		for x := 0; x < width; x++ {
			p, ok := grid[xy{x, y}]

			if ok && p {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
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

func commitToGrid(r rock, pt xy, grid map[xy]bool) {
	for _, rp := range r {
		cp := xy{pt.x + rp.x, pt.y + rp.y}
		grid[cp] = true
	}
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

func sortXY(list []xy) []xy {
	sort.Slice(list, func(i, j int) bool {
		p1, p2 := list[i], list[j]
		return p1.x*100+p1.y < p2.x*100+p2.y
	})
	return list
}

func skeys(m []xy) string {
	strs := []string{}
	for _, k := range sortXY(m) {
		strs = append(strs, fmt.Sprintf("%v", k))
	}
	return strings.Join(strs, ",")
}

func windkey(m []wind, start, length int) string {
	s := ""
	for i := start; i < start+length; i++ {
		s += string(m[i%len(m)])
	}
	return s
}

func keys[K comparable, V any](m map[K]V) []K {
	ks := []K{}
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int64, integers ...int64) int64 {
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
