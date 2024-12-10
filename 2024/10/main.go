package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 10 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parse(in)
	return findAllPaths(m, 0, 9, true)
}

func part2(in input) output {
	m := parse(in)
	return findAllPaths(m, 0, 9, false)
}

func findAllPaths(m map[xy]int, start, goal int, distinct bool) int {
	sum := 0
	starts := findStarts(m, start)
	for _, pt := range starts {
		score := countPaths(m, pt, goal, distinct)
		sum += score
	}
	return sum
}

func countPaths(m map[xy]int, start xy, goal int, distinct bool) int {
	np := 0
	v := make(map[xy]bool)
	queue := []xy{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := v[cur]; ok && distinct {
			continue
		}
		v[cur] = true

		if m[cur] == goal {
			np++
			continue
		}

		mvs := getMoves(m, cur, m[cur])

		queue = append(queue, mvs...)
	}

	return np
}

func getMoves(m map[xy]int, cur xy, v int) []xy {
	dirs := []xy{{cur.x + 1, cur.y}, {cur.x - 1, cur.y}, {cur.x, cur.y + 1}, {cur.x, cur.y - 1}}
	valid := []xy{}
	for _, dir := range dirs {
		if c, ok := m[dir]; ok && c == v+1 {
			valid = append(valid, dir)
		}
	}
	return valid
}

func findStarts(m map[xy]int, goal int) []xy {
	st := []xy{}
	for k, v := range m {
		if v == goal {
			st = append(st, k)
		}
	}
	return st
}

func parse(in []string) map[xy]int {
	res := make(map[xy]int)
	for y, line := range in {
		for x, c := range line {
			v, _ := strconv.Atoi(string(c))
			res[xy{x, y}] = v
		}
	}
	return res
}
