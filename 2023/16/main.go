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

type block struct {
	pt        xy
	mir       *mirror
	spl       *splitter
	energized bool
}

type mirror struct {
	pt    xy
	slant slant
}

type splitter struct {
	pt  xy
	dir splitdir
}

type lightbeam struct {
	pos xy
	dir dir
}

type dir xy
type slant rune
type splitdir rune

var (
	front      slant    = '/'
	back       slant    = '\\'
	north      dir      = dir{0, -1}
	south      dir      = dir{0, 1}
	west       dir      = dir{-1, 0}
	east       dir      = dir{1, 0}
	eastwest   splitdir = '-'
	northsouth splitdir = '|'
)

func print(g map[xy]block) {
	mx, my := maxes(g)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			pt := xy{x, y}
			if b, ok := g[pt]; ok {
				c := '.'
				if b.energized {
					c = '#'
				}
				fmt.Print(string(c))
			}
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
	fmt.Fprintln(w, "--2023 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parseInput(in)
	v := make(map[lightbeam]bool)
	init := lightbeam{pos: xy{0, 0}, dir: east}
	trackLight(m, init, v)
	return countEnergized(m)
}

func part2(in input) output {
	return 0
}

func countEnergized(m map[xy]block) int {
	total := 0
	mx, my := maxes(m)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			b, ok := m[xy{x, y}]
			if !ok {
				continue
			}

			if b.energized {
				total++
			}
		}
	}
	return total
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

func trackLight(m map[xy]block, beam lightbeam, v map[lightbeam]bool) {
	mx, my := maxes(m)
	done := false
	for !done {
		b, ok := m[beam.pos]
		if !ok {
			done = true
			break
		}

		if _, ok := v[beam]; ok {
			// we've seen this beam going in this direction already
			// nothing changed
			beam.pos = beam.pos.add(xy(beam.dir))
			done = true
			continue
		}
		v[beam] = true

		b.energized = true

		if b.mir != nil {
			if b.mir.slant == front {
				if beam.dir == east {
					beam.dir = north
				} else if beam.dir == west {
					beam.dir = south
				} else if beam.dir == north {
					beam.dir = east
				} else if beam.dir == south {
					beam.dir = west
				}
			} else if b.mir.slant == back {
				if beam.dir == east {
					beam.dir = south
				} else if beam.dir == west {
					beam.dir = north
				} else if beam.dir == north {
					beam.dir = west
				} else if beam.dir == south {
					beam.dir = east
				}
			}
		} else if b.spl != nil {
			// continue tracking this beam in this line of logic
			// but track a new light beam on a recursive call
			if b.spl.dir == northsouth && beam.dir.x != 0 {
				beam.dir = north
				splitbeam := lightbeam{pos: beam.pos, dir: south}
				trackLight(m, splitbeam, v)
			} else if b.spl.dir == eastwest && beam.dir.y != 0 {
				beam.dir = west
				splitbeam := lightbeam{pos: beam.pos, dir: east}
				trackLight(m, splitbeam, v)
			}
		}

		m[beam.pos] = b
		beam.pos = beam.pos.add(xy(beam.dir))
		done = beam.pos.x < 0 || beam.pos.y < 0 || beam.pos.x > mx || beam.pos.y > my
	}
}

func parseInput(in input) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			var b block
			pt := xy{x, y}

			if c == '-' || c == '|' {
				b.spl = &splitter{pt: pt}
				if c == '-' {
					b.spl.dir = eastwest
				} else {
					b.spl.dir = northsouth
				}
			} else if c == '/' || c == '\\' {
				b.mir = &mirror{pt: pt}
				if c == '/' {
					b.mir.slant = front
				} else {
					b.mir.slant = back
				}
			}
			m[pt] = b
		}
	}
	return m
}
