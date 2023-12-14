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

func (pt xy) add(p2 xy) xy {
	return xy{pt.x + p2.x, pt.y + p2.y}
}

type rock struct {
	movable bool
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	g := parseInput(in)
	tilt(g, xy{0, -1})
	return calcLoad(g)
}

func part2(in input) output {
	return 0
}

func maxes(m map[xy]rock) (int, int) {
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

func calcLoad(m map[xy]rock) int {
	_, my := maxes(m)

	sum := 0
	for k, v := range m {
		if !v.movable {
			continue
		}
		sum += (my + 1) - k.y
	}
	return sum
}

func tilt(m map[xy]rock, dir xy) {
	mx, my := maxes(m)
	checkx, checky := dir.x != 0, dir.y != 0
	done := false
	for !done {
		moved := 0
		ks := []xy{}
		for k := range m {
			if m[k].movable {
				track := false
				if checky {
					if dir.y == -1 {
						track = k.y > 0
					} else {
						track = k.y < my
					}
				} else if checkx {
					if dir.x == -1 {
						track = k.x > 0
					} else {
						track = k.x < mx
					}
				}
				if track {
					ks = append(ks, k)
				}
			}
		}

		for _, pt := range ks {
			np := pt.add(dir)
			if _, ok := m[np]; !ok {
				v := m[pt]
				delete(m, pt)
				m[np] = v
				moved++
			}
		}
		done = moved == 0
	}
}

func parseInput(in input) map[xy]rock {
	m := make(map[xy]rock)
	for y, line := range in {
		for x, c := range line {
			if c == '.' {
				continue
			}
			pt := xy{x, y}
			r := rock{movable: c == 'O'}

			m[pt] = r
		}
	}
	return m
}
