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

type pipe struct {
	pt xy
	ch rune
}

type state struct {
	pt       xy
	steps    int
	distance int
}

func (p pipe) String() string {
	return fmt.Sprintf("(%d,%d) %c", p.pt.x, p.pt.y, p.ch)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 10 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	start, pipes := parseInput(in)
	return traverse(start, pipes)
}

func part2(in input) output {
	return 0
}

func traverse(start pipe, pipes map[xy]pipe) int {
	maxstate := state{pt: xy{}, steps: 0, distance: 0}

	v := make(map[xy]bool)
	v[start.pt] = true

	var cur state
	queue := []state{{pt: start.pt, steps: 0}}
	for len(queue) > 0 {
		cur = queue[0]
		queue = queue[1:]

		v[cur.pt] = true

		if cur.steps > maxstate.steps {
			maxstate.distance = cur.distance
			maxstate.steps = cur.steps
			maxstate.pt = cur.pt
		}

		mvs := getMoves(cur.pt, pipes)
		for _, mv := range mvs {
			if _, ok := v[mv]; ok {
				continue
			}
			d := distance(mv, start.pt)

			queue = append(queue, state{pt: mv, steps: cur.steps + 1, distance: d})
		}
	}
	return maxstate.steps
}

func distance(p1, p2 xy) int {
	d := math.Abs(float64(p2.y-p1.y)) + math.Abs(float64(p2.x-p1.x))
	return int(d)
}

func getMoves(pt xy, pipes map[xy]pipe) []xy {
	p := pipes[pt]
	mvs := []xy{}

	if p.ch == 'S' {
		for y := pt.y - 1; y <= pt.y+1; y++ {
			for x := pt.x - 1; x <= pt.x+1; x++ {
				p2 := xy{x, y}
				if p2 == pt {
					continue
				}

				smvs := getMoves(p2, pipes)
				for _, smv := range smvs {
					if smv == pt {
						mvs = append(mvs, p2)
					}
				}
			}
		}
		return mvs
	}

	switch p.ch {
	case '|':
		mvs = append(mvs, xy{pt.x, pt.y - 1}, xy{pt.x, pt.y + 1})
	case '-':
		mvs = append(mvs, xy{pt.x - 1, pt.y}, xy{pt.x + 1, pt.y})
	case 'L':
		mvs = append(mvs, xy{pt.x + 1, pt.y}, xy{pt.x, pt.y - 1})
	case 'J':
		mvs = append(mvs, xy{pt.x, pt.y - 1}, xy{pt.x - 1, pt.y})
	case '7':
		mvs = append(mvs, xy{pt.x - 1, pt.y}, xy{pt.x, pt.y + 1})
	case 'F':
		mvs = append(mvs, xy{pt.x, pt.y + 1}, xy{pt.x + 1, pt.y})
	}
	return mvs
}

func parseInput(in input) (pipe, map[xy]pipe) {
	pipes := make(map[xy]pipe)
	var start pipe
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			p := pipe{pt: pt, ch: c}
			pipes[pt] = p
			if c == 'S' {
				start = p
			}
		}
	}
	return start, pipes
}
