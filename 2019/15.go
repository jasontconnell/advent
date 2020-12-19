package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/2019/intcode"
)

var input = "15.txt"

const (
	north int = 1
	south int = 2
	west  int = 3
	east  int = 4
)

var dirs []string = []string{
	"", "north", "south", "west", "east",
}

func dirn(dir int) string {
	return dirs[dir]
}

type mode int

const (
	forward mode = iota
	backtrack
)

type blocktype string

const (
	wall blocktype = "#"
	path           = "."
	goal           = "!"
)

type block struct {
	contents blocktype
}

type xy struct {
	x, y int
}

type neighbor struct {
	xy
	dir int
}

var matrix []neighbor = []neighbor{
	{xy{0, -1}, north}, {xy{0, 1}, south}, {xy{-1, 0}, west}, {xy{1, 0}, east},
}

var dirxy []xy = []xy{
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
}

type robot struct {
	pos     xy
	path    []xy
	moves   []int
	mode    mode
	dir     int
	visited map[xy]block
}

func (p xy) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
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

	r := solve(prog, false)

	generatorPosition := r.pos
	p1 := len(r.moves)

	oxyprog := make([]int, len(opcodes))
	copy(oxyprog, opcodes)

	rexp := solve(oxyprog, true)
	p2 := fillOxygen(rexp.visited, generatorPosition)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func fillOxygen(m map[xy]block, start xy) int {
	openCount := 0
	for _, b := range m {
		if b.contents == wall {
			continue
		}
		openCount++
	}

	omap := make(map[xy]bool)
	omap[start] = true
	minutes := 0

	q := [][]xy{[]xy{start}}

	for len(q) > 0 {
		pts := q[0]
		q = q[1:]
		opts := []xy{}
		oreached := false

		if len(pts) == 0 {
			break
		}

		for _, pt := range pts {
			n := getNeighbors(pt)

			for _, np := range n {
				b, inMap := m[np.xy]
				_, ox := omap[np.xy]
				if ox {
					continue
				}

				if !inMap {
					continue
				}

				if b.contents == wall {
					continue
				}

				if _, ok := omap[np.xy]; !ok {
					opts = append(opts, np.xy)
					omap[np.xy] = true
					oreached = true
				}
			}
		}
		if oreached {
			minutes++
		}
		q = append(q, opts)
	}

	return minutes
}

func getDirXy(dir int) xy {
	return dirxy[dir-1]
}

func getBacktrack(dir int) int {
	bt := 0
	switch dir {
	case north:
		bt = south
	case south:
		bt = north
	case east:
		bt = west
	case west:
		bt = east
	}
	return bt
}

func solve(prog []int, exploreMode bool) robot {
	c := intcode.NewComputer(prog)
	m := make(map[xy]block)
	r := robot{dir: north, pos: xy{0, 0}, path: []xy{}, moves: []int{}, mode: forward, visited: m}

	homeCount := 1

	c.OnOutput = func(out int) {
		switch out {
		case 1, 2:
			mxy := getDirXy(r.dir)
			r.pos = xy{r.pos.x + mxy.x, r.pos.y + mxy.y}

			var bt blocktype = path
			if out == 2 && !exploreMode {
				bt = goal
			}

			b := block{contents: bt}
			r.visited[r.pos] = b

			if r.mode == forward {
				r.path = append(r.path, r.pos)
				r.moves = append(r.moves, r.dir)
			} else {
				r.path = r.path[:len(r.path)-1]
				r.moves = r.moves[:len(r.moves)-1]
			}

			if out == 2 && !exploreMode {
				c.Complete = true
				break
			}

			if exploreMode && (r.pos.x == 0 && r.pos.y == 0) {
				homeCount++
				c.Complete = homeCount == 5 // start at 0,0, plus come back to it 4 times maybe.
				if c.Complete {
					break
				}
			}

			d, m := bestDir(r)
			r.dir = d
			r.mode = m
		case 0:
			mxy := getDirXy(r.dir)
			wpos := xy{r.pos.x + mxy.x, r.pos.y + mxy.y}
			b := block{contents: wall}
			r.visited[wpos] = b

			d, m := bestDir(r)
			r.dir = d
			r.mode = m
		}

		if !c.Complete {
			c.AddInput(r.dir)
		}
	}

	c.AddInput(r.dir)
	c.Exec()

	return r
}

func bestDir(r robot) (int, mode) {
	n := getNeighbors(r.pos)
	m := r.mode
	best := 0
	unvisited := false

	for _, np := range n {
		if _, ok := r.visited[np.xy]; !ok {
			best = np.dir
			m = forward
			unvisited = true
			break
		}
	}

	if !unvisited && len(r.moves) > 0 {
		best = getBacktrack(r.moves[len(r.moves)-1])
		m = backtrack
	}

	return best, m
}

func getNeighbors(p xy) []neighbor {
	n := []neighbor{}
	for _, mp := range matrix {
		p2 := xy{p.x + mp.x, p.y + mp.y}
		n = append(n, neighbor{xy: p2, dir: mp.dir})
	}
	return n
}
