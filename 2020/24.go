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
