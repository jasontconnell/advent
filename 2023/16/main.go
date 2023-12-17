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
	mir       *mirror
	spl       *splitter
	energized bool
}

type mirror struct {
	slant slant
}

type splitter struct {
	dir splitdir
}

type lightbeam struct {
	pos xy
	dir dir
}

type dir xy
type slant rune
type splitdir rune

const (
	front      slant    = '/'
	back       slant    = '\\'
	eastwest   splitdir = '-'
	northsouth splitdir = '|'
)

var (
	north dir = dir{0, -1}
	south dir = dir{0, 1}
	west  dir = dir{-1, 0}
	east  dir = dir{1, 0}
)

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
	m := parseInput(in)
	return findMaxEnergized(m)
}

func countEnergized(m map[xy]block) int {
	total := 0
	mx, my := maxes(m)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			if m[xy{x, y}].energized {
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

func findMaxEnergized(m map[xy]block) int {
	mx, my := maxes(m)
	vals := []int{}
	for y := 0; y <= my; y++ {
		left := calculateEnergized(m, lightbeam{pos: xy{0, y}, dir: east})
		right := calculateEnergized(m, lightbeam{pos: xy{mx, y}, dir: west})
		vals = append(vals, left, right)
	}

	for x := 0; x <= mx; x++ {
		top := calculateEnergized(m, lightbeam{pos: xy{x, 0}, dir: south})
		bottom := calculateEnergized(m, lightbeam{pos: xy{x, my}, dir: north})
		vals = append(vals, top, bottom)
	}

	return max(vals)
}

func max(list []int) int {
	m := 0
	for _, v := range list {
		if v > m {
			m = v
		}
	}
	return m
}

func reset(m map[xy]block) {
	for k, b := range m {
		b.energized = false
		m[k] = b
	}
}

func calculateEnergized(m map[xy]block, beam lightbeam) int {
	reset(m)
	trackLight(m, beam, make(map[lightbeam]bool))
	return countEnergized(m)
}

func deflect(s slant, d dir) dir {
	var nd dir = dir{x: d.y, y: d.x}
	if s == front {
		nd.x = -nd.x
		nd.y = -nd.y
	}
	return nd
}

func trackLight(m map[xy]block, beam lightbeam, v map[lightbeam]bool) {
	done := false
	for !done {
		b, ok := m[beam.pos]
		if !ok {
			break
		}

		if _, ok := v[beam]; ok {
			break
		}

		v[beam] = true
		b.energized = true

		if b.mir != nil {
			beam.dir = deflect(b.mir.slant, beam.dir)
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
	}
}

func parseInput(in input) map[xy]block {
	m := make(map[xy]block)
	for y, line := range in {
		for x, c := range line {
			var b block
			pt := xy{x, y}

			if c == '-' || c == '|' {
				b.spl = &splitter{}
				if c == '-' {
					b.spl.dir = eastwest
				} else {
					b.spl.dir = northsouth
				}
			} else if c == '/' || c == '\\' {
				b.mir = &mirror{}
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
