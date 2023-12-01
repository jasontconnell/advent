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

type content rune

const (
	space content = '.'
	elf   content = '#'
)

type xy struct {
	x, y int
}

type dir string

const (
	STAY dir = ""
	N    dir = "N"
	S    dir = "S"
	W    dir = "W"
	E    dir = "E"
	NE   dir = "NE"
	NW   dir = "NW"
	SE   dir = "SE"
	SW   dir = "SW"
)

type side string

const (
	left  side = "left"
	right side = "right"
	up    side = "up"
	down  side = "down"
	all   side = "all"
)

var dirdelta map[dir]xy = map[dir]xy{
	N: {0, -1}, S: {0, 1}, E: {1, 0}, W: {-1, 0},
	NE: {1, -1}, NW: {-1, -1}, SE: {1, 1}, SW: {-1, 1},
}

var checkpoints map[side][]xy = map[side][]xy{
	left:  {dirdelta[NW], dirdelta[W], dirdelta[SW]},
	up:    {dirdelta[NW], dirdelta[N], dirdelta[NE]},
	down:  {dirdelta[SW], dirdelta[S], dirdelta[SE]},
	right: {dirdelta[NE], dirdelta[E], dirdelta[SE]},
	all:   {dirdelta[N], dirdelta[S], dirdelta[E], dirdelta[W], dirdelta[NW], dirdelta[NE], dirdelta[SW], dirdelta[SE]},
}

var dirside map[dir]side = map[dir]side{
	N: up, S: down, E: right, W: left,
}

var defaultProposeOrder []dir = []dir{N, S, W, E}

type block struct {
	contents     content
	proposedMove dir
}

func print(grid map[xy]block) {
	min, max := minmax(grid)
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			b, ok := grid[xy{x, y}]
			if !ok {
				fmt.Print(".")
			} else {
				fmt.Print(string(b.contents))
			}
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
	fmt.Fprintln(w, "--2022 day 23 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	grid := parseInput(in)
	prop := make([]dir, len(defaultProposeOrder))
	copy(prop, defaultProposeOrder)
	return simulate(grid, prop, 10)
}

func part2(in input) output {
	grid := parseInput(in)
	prop := make([]dir, len(defaultProposeOrder))
	copy(prop, defaultProposeOrder)
	return findNoMoves(grid, prop)
}

func countSpace(grid map[xy]block) int {
	min, max := minmax(grid)
	count := 0
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			b, ok := grid[xy{x, y}]
			if !ok || b.contents == space {
				count++
			}
		}
	}
	return count
}

func findNoMoves(grid map[xy]block, proposed []dir) int {
	var moved bool
	var round int = 1
	for {
		grid, proposed, moved = simulateOne(grid, proposed)
		if !moved {
			break
		}
		round++
	}
	return round
}

func simulate(grid map[xy]block, proposed []dir, rounds int) int {
	for i := 0; i < rounds; i++ {
		grid, proposed, _ = simulateOne(grid, proposed)
	}
	return countSpace(grid)
}

func simulateOne(grid map[xy]block, proposed []dir) (map[xy]block, []dir, bool) {
	for pt, b := range grid {
		if b.contents != elf {
			continue
		}

		anyElves := !isFree(grid, pt, checkpoints[all])
		b.proposedMove = STAY
		if anyElves {
			for _, d := range proposed {
				if isFree(grid, pt, checkpoints[dirside[d]]) {
					b.proposedMove = d
					break
				}
			}
		}

		grid[pt] = b // update the map
	}

	mvmap := make(map[xy][]xy)
	for pt, b := range grid {
		if b.proposedMove != STAY {
			to := xy{pt.x + dirdelta[b.proposedMove].x, pt.y + dirdelta[b.proposedMove].y}
			mvmap[to] = append(mvmap[to], pt)
		}
	}

	moved := false
	for k, v := range mvmap {
		if len(v) == 1 {
			from := v[0]
			elf := grid[from]
			delete(grid, from)
			grid[k] = elf
			moved = true
		}
	}

	proposed = moveToBack(proposed, proposed[0])

	return grid, proposed, moved
}

func isFree(grid map[xy]block, pt xy, deltas []xy) bool {
	free := true
	for _, delta := range deltas {
		check := xy{pt.x + delta.x, pt.y + delta.y}
		if b, ok := grid[check]; ok && b.contents == elf {
			free = false
			break
		}
	}
	return free
}

func moveToBack[T comparable](list []T, item T) []T {
	idx := -1
	for i := 0; i < len(list); i++ {
		if item == list[i] {
			idx = i
			continue
		}
	}
	if idx == -1 {
		return list
	}
	list = append(list[:idx], list[idx+1:]...)
	list = append(list, item)
	return list
}

func minmax(grid map[xy]block) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for k := range grid {
		if grid[k].contents == space {
			continue
		}
		if k.x < min.x {
			min.x = k.x
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.y > max.y {
			max.y = k.y
		}
	}
	return min, max
}

func parseInput(in input) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			c := content(c)

			b := block{contents: c}
			m[pt] = b
		}
	}
	return m
}
