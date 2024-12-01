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
	point xy
	steps int
}

func print(m map[xy]block, path []common.Edge[xy, int]) {
	mx, my := maxes(m)
	cm := make(map[xy]rune)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			pt := xy{x, y}
			c := '#'
			b := m[pt]
			if b.open {
				c = '.'
			}
			if b.slope {
				c = 'S'
			}
			cm[pt] = c
		}
	}

	for _, e := range path {
		p := e.GetLeft()
		cm[p] = 'O'
	}

	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			pt := xy{x, y}
			c := cm[pt]
			fmt.Print(string(c))
		}
		fmt.Println()
	}
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
	start, end := findStartEnd(m)
	g := getGraph(m, start, end, slippery)
	edges := g.GetEdges()
	for _, e := range edges {
		fmt.Println(e.GetLeft(), e.GetRight(), e.GetWeight())
	}
	path := g.DFS(start, end)
	sum := 0
	for _, edge := range path {
		sum += edge.GetWeight()
	}
	return sum
}

func getGraph(m map[xy]block, start, end xy, slippery bool) common.Graph[xy, int] {
	g := common.NewGraph[xy]()

	intersections := getIntersections(m, slippery)
	intersections[start] = true
	intersections[end] = true

	list := []xy{start, end}
	for k := range intersections {
		g.AddVertex(k)
		if k != start && k != end {
			list = append(list, k)
		}
	}

	for k := range intersections {
		queue := []state{{point: k, steps: 0}}
		v := make(map[xy]bool)

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

			if _, ok := intersections[cur.point]; ok && k != cur.point {
				g.AddWeightedEdge(k, cur.point, cur.steps)
				continue
			}

			mvs := getMoves(m, cur.point, slippery)
			for _, mv := range mvs {
				if _, ok := v[mv]; ok {
					continue
				}
				v[mv] = true

				st := state{point: mv, steps: cur.steps + 1}
				queue = append(queue, st)
			}
		}
	}
	return g
}

func getIntersections(m map[xy]block, slippery bool) map[xy]bool {
	imap := make(map[xy]bool)
	for k, b := range m {
		if !b.open {
			continue
		}
		opencount := 0
		for _, d := range []xy{north, south, east, west} {
			dest := k.add(d)
			db, ok := m[dest]
			if !ok || !db.open {
				continue
			}
			if db.open {
				opencount++
			}
		}
		if opencount > 2 {
			imap[k] = true
		}
	}
	return imap
}

func getMoves(m map[xy]block, pt xy, slippery bool) []xy {
	onslope := m[pt].slope
	slopedir := m[pt].slopedir
	mvs := []xy{}
	for _, d := range []xy{north, south, east, west} {
		dest := pt.add(d)
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
