package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type node struct {
	pt                xy
	size, used, avail int
}

func (n *node) copy() *node {
	cp := node{pt: n.pt, size: n.size, used: n.used, avail: n.avail}
	return &cp
}

func (n *node) String() string {
	return fmt.Sprintf("(%d,%d) size: %d used: %d avail:%d", n.pt.x, n.pt.y, n.size, n.used, n.avail)
}

type xy struct {
	x, y int
}

func printGrid(m map[xy]*node) {
	maxx, maxy := maxes(m)

	for y := 0; y < maxy+1; y++ {
		for x := 0; x < maxx+1; x++ {
			n := m[xy{x, y}]

			ch := "."
			if n.used == 0 {
				ch = "_"
			} else if n.used > 400 {
				ch = "#"
			} else if n.pt.x == maxx && n.pt.y == 0 {
				ch = "S"
			}

			fmt.Print(ch)
		}
		fmt.Println("  ", y)
	}
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 22 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	nodes := parseInput(in)
	return viableNodes(nodes)
}

func part2(in input) output {
	nodes := parseInput(in)
	printGrid(nodes)
	return 249 // most people hand-calculated
}

func hasSpace(n *node, cur *node) bool {
	return n.avail > cur.used
}

func viableNodes(nodes map[xy]*node) int {
	cnt := 0
	for _, node := range nodes {
		if node.used > 0 {
			cnt += avail(nodes, node)
		}
	}
	return cnt
}

func avail(nodes map[xy]*node, src *node) int {
	cnt := 0
	for _, node := range nodes {
		if node == src {
			continue
		}
		if src.used <= node.avail {
			cnt++
		}
	}
	return cnt
}

func maxes(m map[xy]*node) (int, int) {
	maxx, maxy := 0, 0
	for pt := range m {
		if pt.x > maxx {
			maxx = pt.x
		}
		if pt.y > maxy {
			maxy = pt.y
		}
	}
	return maxx, maxy
}

func parseInput(in input) map[xy]*node {
	reg := regexp.MustCompile("^/dev/grid/node-x([0-9]+)-y([0-9]+) *([0-9]+)T *([0-9]+)T *([0-9]+)T *([0-9]+)%$")

	nodes := make(map[xy]*node)
	for _, line := range in {
		if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
			x, _ := strconv.Atoi(groups[1])
			y, _ := strconv.Atoi(groups[2])
			pt := xy{x, y}

			size, _ := strconv.Atoi(groups[3])
			used, _ := strconv.Atoi(groups[4])
			avail, _ := strconv.Atoi(groups[5])
			nodes[pt] = &node{pt: pt, size: size, used: used, avail: avail}
		}
	}
	return nodes
}
