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

func (pt xy) add(pt2 xy) xy {
	return xy{pt.x + pt2.x, pt.y + pt2.y}
}

func maxPoint(m map[xy]bool) xy {
	max := xy{0, 0}
	for k := range m {
		if k.x >= max.x && k.y >= max.y {
			max = k
		}
	}
	return max
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 07 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m, start := parseInput(in)
	return simulate(m, start)
}

func part2(in input) output {
	m, start := parseInput(in)
	return countTimelines(m, start)
}

func simulate(m map[xy]bool, start xy) int {
	down, left, right := xy{0, 1}, xy{-1, 0}, xy{1, 0}
	queue := []xy{start.add(down)}
	visit := make(map[xy]bool)
	count := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if _, ok := visit[cur]; ok {
			// doubled up beams don't count
			continue
		}
		visit[cur] = true

		if _, ok := m[cur]; !ok {
			// reached bottom
			continue
		}

		if split, ok := m[cur]; ok && split {
			lbeam, rbeam := cur.add(left), cur.add(right)
			queue = append(queue, lbeam, rbeam)
			count++
		} else {
			dbeam := cur.add(down)
			queue = append(queue, dbeam)
		}
	}
	return count
}

func countTimelines(m map[xy]bool, start xy) int {
	max := maxPoint(m)
	paths := make(map[xy]int)
	for y := max.y; y >= 0; y-- {
		for x := 0; x <= max.x; x++ {
			pt := xy{x, y}
			if pt == start {
				continue
			}
			if split, ok := m[pt]; ok && split {
				paths[pt] = calcPaths(m, paths, pt)
			}
		}
	}
	log.Println(paths)
	return calcPaths(m, paths, start)
}

func calcPaths(m map[xy]bool, paths map[xy]int, start xy) int {
	down, left, right := xy{0, 1}, xy{-1, 0}, xy{1, 0}
	queue := []xy{start.add(down)}
	count := 1

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if fromHere, ok := paths[cur]; ok {
			log.Println("found path from", cur, fromHere)
			count += fromHere
			break
		}

		if _, ok := m[cur]; !ok {
			// reached bottom
			continue
		}

		if split, ok := m[cur]; ok && split {
			lbeam, rbeam := cur.add(left), cur.add(right)
			queue = append([]xy{lbeam, rbeam}, queue...)
			count += 2
		} else {
			dbeam := cur.add(down)
			queue = append([]xy{dbeam}, queue...)
		}
	}
	return count
}

func parseInput(in input) (map[xy]bool, xy) {
	var start xy
	m := make(map[xy]bool)

	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			if c == 'S' {
				start = pt
			}
			m[pt] = c == '^'
		}
	}

	return m, start
}
