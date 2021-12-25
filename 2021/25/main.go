package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int
type xy struct {
	x, y int
}
type dir string

const (
	east  dir = ">"
	south dir = "v"
)

type move struct {
	from, to xy
	cuc      cucumber
}
type cucumber struct {
	facing dir
}
type grid struct {
	m   map[xy]cucumber
	max xy
}

func (g grid) String() string {
	s := ""
	for y := 0; y < g.max.y; y++ {
		for x := 0; x < g.max.x; x++ {
			pt := xy{x, y}
			if cuc, ok := g.m[pt]; ok {
				s += string(cuc.facing)
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 25 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	g := parseInput(in)
	return findStop(g)
}

func findStop(g grid) int {
	s := 1
	done := false
	for !done {
		var emvs, smvs int

		g, emvs = processOne(g, east)
		g, smvs = processOne(g, south)

		done = emvs+smvs == 0

		if !done {
			s++
		}
	}
	return s
}

func processOne(g grid, facing dir) (grid, int) {
	mvs := []move{}
	for y := 0; y < g.max.y; y++ {
		for x := 0; x < g.max.x; x++ {
			pt := xy{x, y}
			cuc, ok := g.m[pt]
			if !ok {
				continue
			}
			if cuc.facing != facing {
				continue
			}
			apt := pt
			if cuc.facing == south {
				apt.y = (pt.y + 1) % g.max.y
			} else {
				apt.x = (pt.x + 1) % g.max.x
			}

			if _, ok := g.m[apt]; !ok {
				mvs = append(mvs, move{from: pt, to: apt, cuc: cuc})
			}
		}
	}

	for _, mv := range mvs {
		g.m[mv.to] = mv.cuc
		delete(g.m, mv.from)
	}

	gp := grid{max: g.max, m: map[xy]cucumber{}}
	for k, v := range g.m {
		gp.m[k] = v
	}
	return gp, len(mvs)
}

func parseInput(in input) grid {
	g := grid{}
	g.m = map[xy]cucumber{}

	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			if c != '.' {
				face := dir(string(c))
				cuc := cucumber{facing: face}
				g.m[pt] = cuc
			}
		}
	}
	g.max = xy{len(in[0]), len(in)}
	return g
}
