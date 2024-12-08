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

type block struct {
	c       rune
	antenna bool
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 08 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parse(in)
	return findAntinodes(m)
}

func part2(in input) output {
	return 0
}

func minmax(g map[xy]block) (xy, xy) {
	min := xy{math.MaxInt32, math.MaxInt32}
	max := xy{math.MinInt32, math.MinInt32}

	for k := range g {
		if k.x < min.x {
			min.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y > max.y {
			max.y = k.y
		}
	}
	return min, max
}

func findAntinodes(m map[xy]block) int {
	lookup := make(map[rune][]xy)
	for k, v := range m {
		if v.antenna {
			lookup[v.c] = append(lookup[v.c], k)
		}
	}

	min, max := minmax(m)

	antinodes := make(map[xy]bool)
	for _, v := range lookup {
		for i, p1 := range v {
			for j, p2 := range v {
				if i == j {
					continue
				}

				an := plotAntinodes(p1, p2)
				for _, p := range an {
					if p.x < min.x || p.y < min.y || p.x > max.x || p.y > max.y {
						continue
					}
					antinodes[p] = true
				}
			}
		}
	}
	return len(antinodes)
}

func distance(p1, p2 xy) int {
	return int(math.Abs(math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y))))
}

func plotAntinodes(p1, p2 xy) []xy {
	dx := p1.x - p2.x
	dy := p1.y - p2.y

	test := []xy{
		{p1.x - dx, p1.y - dy},
		{p1.x + dx, p1.y + dy},
		{p1.x - dx, p1.y + dy},
		{p1.x + dx, p1.y - dy},
		{p2.x - dx, p2.y - dy},
		{p2.x + dx, p2.y + dy},
		{p2.x - dx, p2.y + dy},
		{p2.x + dx, p2.y - dy},
	}

	chkdist := distance(p1, p2)

	pts := []xy{}
	for _, p := range test {
		dp1 := distance(p, p1)
		dp2 := distance(p, p2)
		if (dp1 == chkdist && dp2 == chkdist*2) || (dp1 == chkdist*2 && dp2 == chkdist) {
			pts = append(pts, p)
		}
	}

	return pts
}

func parse(in []string) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			b := block{c: c}
			if c != '.' {
				b.antenna = true
			}
			m[pt] = b
		}
	}
	return m
}
