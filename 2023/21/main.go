package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int64

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
	steps int64
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
	start, m := parseInput(in)
	return calculateLarge(start, m, 26501365)
}

func calculateLarge(start xy, m map[xy]block, steps int64) int64 {
	size := int64(math.Sqrt(float64(len(m))))

	gridw := steps/size - 1
	oddg := int64(math.Pow(float64(gridw/2*2+1), 2))
	eveng := int64(math.Pow(float64((gridw+1)/2*2), 2))

	odd := pointsFromStart(start, m, size*2+1)
	even := pointsFromStart(start, m, size*2)

	var corners int64
	for _, p := range []xy{
		{start.x, 0},
		{start.x, int(size - 1)},
		{0, start.y},
		{int(size - 1), start.y},
	} {
		corners += pointsFromStart(p, m, size-1)
	}

	var smallcorners int64
	for _, p := range []xy{
		{0, 0},
		{0, int(size - 1)},
		{int(size - 1), 0},
		{int(size - 1), int(size - 1)},
	} {
		smallcorners += pointsFromStart(p, m, size/2-1)
	}

	var largecorners int64
	for _, p := range []xy{
		{0, 0},
		{0, int(size - 1)},
		{int(size - 1), 0},
		{int(size - 1), int(size - 1)},
	} {
		largecorners += pointsFromStart(p, m, 3*size/2-1)
	}

	return oddg*odd + eveng*even + corners + (gridw+1)*smallcorners + gridw*largecorners
}

func pointsFromStart(start xy, m map[xy]block, steps int64) int64 {
	queue := common.NewQueue[state, int]()
	queue.Enqueue(state{pt: start, steps: steps})

	var goals int64

	v := make(map[xy]bool)
	for queue.Any() {
		cur := queue.Dequeue()

		if _, ok := v[cur.pt]; ok {
			continue
		}
		v[cur.pt] = true

		if cur.steps%2 == 0 {
			goals++
		}

		if cur.steps == 0 {
			continue
		}

		mvs := getMoves(m, cur.pt, cur.steps)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	return goals
}

func getMoves(m map[xy]block, pt xy, steps int64) []state {
	dirs := []xy{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	mvs := []state{}
	for _, d := range dirs {
		np := pt.add(d)

		if b, ok := m[np]; !ok || b.rock {
			continue
		}
		st := state{pt: np, steps: steps - 1}
		mvs = append(mvs, st)
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
