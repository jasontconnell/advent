package main

import (
	"fmt"
	"log"
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

type block struct {
	ch     rune
	rock   bool
	garden bool
	start  bool
}

type state struct {
	pt    xy
	prev  xy
	steps int
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

func print(m map[xy]block, v map[xy]bool) {
	mx, my := maxes(m)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			pt := xy{x, y}
			b := m[pt]
			c := b.ch
			if _, ok := v[pt]; ok {
				c = 'O'
			}
			fmt.Print(string(c))
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 21 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	start, m := parseInput(in)
	return pointsFromStart(start, m, 64)
}

func part2(in input) output {
	return 0
}

func pointsFromStart(start xy, m map[xy]block, steps int) int {
	queue := common.NewQueue[state, int]()
	queue.Enqueue(state{pt: start, steps: 0})

	goals := make(map[xy]bool)

	vs := make(map[int]map[xy]bool)

	for queue.Any() {
		cur := queue.Dequeue()
		if vs[cur.steps] == nil {
			vs[cur.steps] = make(map[xy]bool)
		} else if _, ok := vs[cur.steps][cur.pt]; ok {
			continue
		}
		vs[cur.steps][cur.pt] = true

		if _, ok := goals[cur.pt]; ok {
			continue
		}

		if cur.steps == steps {
			goals[cur.pt] = true
			continue
		}

		mvs := getMoves(m, cur.pt, cur.prev)
		for _, mv := range mvs {
			st := state{pt: mv, prev: cur.pt, steps: cur.steps + 1}
			queue.Enqueue(st)
		}
	}
	return len(goals)
}

func getMoves(m map[xy]block, pt, prev xy) []xy {
	dirs := []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	mvs := []xy{}
	for _, d := range dirs {
		np := pt.add(d)
		if b, ok := m[np]; !ok || b.rock || pt == prev {
			continue
		}
		mvs = append(mvs, np)
	}
	return mvs
}

func parseInput(in input) (xy, map[xy]block) {
	m := make(map[xy]block)
	var startpt xy
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			rock := c == '#'
			garden := c == '.' || c == 'S'
			start := c == 'S'
			b := block{ch: c, rock: rock, garden: garden, start: start}
			if start {
				startpt = pt
			}
			m[pt] = b
		}
	}
	return startpt, m
}
