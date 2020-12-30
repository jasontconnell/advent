package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "24.txt"

type xy struct {
	x, y float32
}

type direction string

const (
	ne direction = "ne"
	nw direction = "nw"
	se direction = "se"
	sw direction = "sw"
	e  direction = "e"
	w  direction = "w"
)

type path struct {
	moves []xy
}

var deltas []xy = []xy{
	{-1, 0},    // west
	{1, 0},     // east
	{-.5, -.5}, // sw
	{.5, .5},   // ne
	{-.5, .5},  // nw
	{.5, -.5},  // se
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	directions := readDirections(lines)
	paths := getPaths(directions)
	p1 := getBlack(paths)

	fmt.Println("Part 1:", p1)

	p2 := simulate(paths, 100)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getBlack(paths []path) int {
	// false is white
	m := make(map[xy]bool)
	for _, path := range paths {
		last := path.moves[len(path.moves)-1]
		m[last] = !m[last]
	}

	count := 0
	for _, v := range m {
		if v {
			count++
		}
	}
	return count
}

func countBlack(m map[xy]bool) int {
	count := 0
	for _, v := range m {
		if v {
			count++
		}
	}
	return count
}

func simulate(paths []path, count int) int {
	m := make(map[xy]bool)
	for _, path := range paths {
		last := path.moves[len(path.moves)-1]
		m[last] = !m[last]
	}

	fmt.Println("black to start", countBlack(m))
	for i := 0; i < count; i++ {
		m = fillSurrounding(m)
		m = simulateOne(m)
		if i < 11 || (i+1)%10 == 0 {
			fmt.Println("day", i+1, countBlack(m), "map size", len(m))
		}
	}

	return countBlack(m)
}

func fillSurrounding(m map[xy]bool) map[xy]bool {
	bps := []xy{}
	for k, v := range m {
		if v {
			bps = append(bps, k)
		}
	}

	for _, pt := range bps {
		for _, d := range deltas {
			dv := xy{pt.x + d.x, pt.y + d.y}
			if _, ok := m[dv]; !ok {
				m[dv] = false
			}
		}
	}
	return m
}

func simulateOne(m map[xy]bool) map[xy]bool {
	cm := make(map[xy]bool, len(m))

	for k, v := range m {
		blkcnt := 0

		for _, d := range deltas {
			pt := xy{k.x + d.x, k.y + d.y}
			if dv, ok := m[pt]; ok && dv {
				blkcnt++
			}
		}

		// true is black
		wv := v
		if v && (blkcnt == 0 || blkcnt > 2) {
			wv = false
		} else if !v && blkcnt == 2 {
			wv = true
		}

		cm[k] = wv
	}
	return cm
}

func getPaths(directions [][]direction) []path {
	start := xy{0, 0}
	cur := start
	paths := []path{}
	for _, list := range directions {

		pth := path{}
		for _, d := range list {
			pt := xy{cur.x, cur.y}
			switch d {
			case e:
				pt.x = pt.x + 1
			case w:
				pt.x = pt.x - 1
			case nw:
				pt.x = pt.x - .5
				pt.y = pt.y + .5
			case ne:
				pt.x = pt.x + .5
				pt.y = pt.y + .5
			case sw:
				pt.x = pt.x - .5
				pt.y = pt.y - .5
			case se:
				pt.x = pt.x + .5
				pt.y = pt.y - .5
			}
			cur = pt
			pth.moves = append(pth.moves, pt)
		}
		paths = append(paths, pth)
		cur = xy{0, 0}
	}
	return paths
}

func readDirections(lines []string) [][]direction {
	dirs := make([][]direction, len(lines))
	for j, line := range lines {
		for i := 0; i < len(line); i++ {
			p := ""
			ch := line[i]
			if ch == 'n' || ch == 's' {
				p = string(ch) + string(line[i+1])
				i++
			} else {
				p = string(ch)
			}

			dirs[j] = append(dirs[j], direction(p))
		}
	}
	return dirs
}
