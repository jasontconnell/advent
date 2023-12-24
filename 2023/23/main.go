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

type xy struct {
	x, y int
}

func (pt xy) add(p2 xy) xy {
	return xy{pt.x + p2.x, pt.y + p2.y}
}

func (pt xy) opp() xy {
	return xy{pt.x * -1, pt.y * -1}
}

var (
	none  xy = xy{0, 0}
	north xy = xy{0, -1}
	south xy = xy{0, 1}
	west  xy = xy{-1, 0}
	east  xy = xy{1, 0}
)

type block struct {
	open     bool
	slope    bool
	slopedir xy
}

type state struct {
	path  []xy
	steps int
	point xy
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 23 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parseInput(in)
	return getLongestPath(m, true)
}

func part2(in input) output {
	m := parseInput(in)
	return getLongestPath(m, false)
}

func findStartEnd(m map[xy]block) (xy, xy) {
	var start, end xy
	_, my := maxes(m)
	for k := range m {
		open := m[k].open
		if k.y == 0 && open {
			start = k
		} else if k.y == my && open {
			end = k
		}
	}
	return start, end
}

func maxes(m map[xy]block) (int, int) {
	mx, my := 0, 0
	for k := range m {
		if k.x > mx {
			mx = k.x
		}
		if k.y > my {
			my = k.y
		}
	}
	return mx, my
}

func dist(p1, p2 xy) int {
	return int(math.Abs(float64(p2.x-p1.x)) + math.Abs(float64(p2.y-p1.y)))
}

func getLongestPath(m map[xy]block, slippery bool) int {
	max := 0
	start, end := findStartEnd(m)
	queue := common.NewPriorityQueue(func(s state) int {
		return -s.steps
	})
	// queue := common.NewQueue[state, int]()
	v := make(map[xy]bool)
	queue.Enqueue(state{path: []xy{start}, point: start, steps: 0})
	for queue.Any() {
		cur := queue.Dequeue()

		if cur.point == end {
			if cur.steps > max {
				max = cur.steps
			}
			continue
		}

		if _, ok := v[cur.point]; ok {
			continue
		}
		v[cur.point] = true

		mvs := getMoves(m, cur, slippery)
		for _, mv := range mvs {
			queue.Enqueue(state{point: mv, path: append(cur.path, mv), steps: cur.steps + 1})
		}
	}
	return max
}

func getMoves(m map[xy]block, st state, slippery bool) []xy {
	onslope := m[st.point].slope
	slopedir := m[st.point].slopedir
	mvs := []xy{}
	for _, d := range []xy{north, south, east, west} {
		dest := st.point.add(d)
		if p, ok := m[dest]; !ok || !p.open {
			continue
		}

		if slippery && onslope && d != slopedir {
			continue
		}

		// check against the grain
		s := m[dest]
		if s.slope && slippery {
			if s.slopedir.opp() == d {
				continue
			}
		}

		mvs = append(mvs, dest)
	}
	return mvs
}

func parseInput(in input) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			open := c != '#'
			slopedir := none
			switch c {
			case '>':
				slopedir = east
			case '<':
				slopedir = west
			case '^':
				slopedir = north
			case 'v':
				slopedir = south
			}
			b := block{open: open, slope: slopedir != none, slopedir: slopedir}
			m[pt] = b
		}
	}
	return m
}
