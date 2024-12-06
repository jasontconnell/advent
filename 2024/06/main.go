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

type dir int

const (
	E dir = iota
	W
	N
	S
)

var turndir map[dir]xy = map[dir]xy{
	E: {1, 0}, W: {-1, 0}, S: {0, 1}, N: {0, -1},
}

type xy struct {
	x, y int
}

type visit struct {
	pt xy
	d  dir
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

type block struct {
	wall bool
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 06 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	grid, start, d := parse(in)
	visited, _ := patrol(grid, start, d)
	return visited
}

func part2(in input) output {
	grid, start, d := parse(in)
	obstructions := placeObstructions(grid, start, d)
	return obstructions
}

func minmax(g map[xy]block) (xy, xy) {
	min := xy{math.MaxInt32, math.MaxInt32}
	max := xy{math.MinInt32, math.MinInt32}

	for k := range g {
		if k.x < min.x {
			min.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y > max.y {
			max.y = k.y
		}
	}
	return min, max
}

func placeObstructions(g map[xy]block, start xy, startdir dir) int {
	count := 0
	cp := make(map[xy]block)
	for k, v := range g {
		cp[k] = v
	}

	for k := range cp {
		if cp[k].wall {
			continue
		}

		cp[k] = block{wall: true}
		if _, looped := patrol(cp, start, startdir); looped {
			count++
		}
		cp[k] = block{wall: false}
	}
	return count
}

func patrol(g map[xy]block, start xy, startdir dir) (int, bool) {
	visited := make(map[visit]bool)
	min, max := minmax(g)
	cur := start
	d := startdir
	loops := false
	for {
		if cur.x < min.x || cur.y < min.y || cur.x > max.x || cur.y > max.y {
			break
		}

		vkey := visit{pt: cur, d: d}
		if _, ok := visited[vkey]; ok {
			loops = true
			break
		}

		visited[vkey] = true
		test := cur.add(turndir[d])
		if b, ok := g[test]; (ok && !b.wall) || !ok {
			cur = test
			continue
		} else if b.wall {
			d = turnRight(d)
		}
	}
	return len(visited), loops
}

func turnRight(d dir) dir {
	var nd dir
	switch d {
	case E:
		nd = S
	case W:
		nd = N
	case S:
		nd = W
	case N:
		nd = E
	}
	return nd
}

func parse(in []string) (map[xy]block, xy, dir) {
	m := make(map[xy]block)
	var d dir
	var start xy
	for y := 0; y < len(in); y++ {
		for x := 0; x < len(in[y]); x++ {
			c := in[y][x]
			m[xy{x, y}] = block{wall: c == '#'}

			if c != '.' && c != '#' {
				start = xy{x, y}
				switch c {
				case '<':
					d = W
				case '>':
					d = E
				case '^':
					d = N
				case 'v':
					d = S
				}
			}
		}
	}
	return m, start, d
}
