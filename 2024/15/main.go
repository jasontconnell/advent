package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

type block struct {
	wall, box, open bool
}

type grid map[xy]block

func print(g grid, r xy) {
	for y := 0; y < 11; y++ {
		for x := 0; x < 11; x++ {
			pt := xy{x, y}
			b, ok := g[pt]
			if !ok {
				continue
			}
			if r == pt {
				fmt.Print("@")
			} else if b.wall {
				fmt.Print("#")
			} else if b.box {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type dir int

const (
	Up dir = iota
	Down
	Left
	Right
)

func (d dir) xy() xy {
	var p xy
	switch d {
	case Up:
		p.y = -1
	case Down:
		p.y = 1
	case Left:
		p.x = -1
	case Right:
		p.x = 1
	}
	return p
}

func getDir(r rune) dir {
	var d dir
	switch r {
	case '^':
		d = Up
	case 'v':
		d = Down
	case '<':
		d = Left
	case '>':
		d = Right
	}
	return d
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m, start, dirs := parse(in)
	m, _ = navigate(m, start, dirs)
	return calcGPS(m)
}

func part2(in input) output {
	for i := 0; i < len(in); i++ {
		in[i] = double(in[i])
	}
	// log.Println(in)
	return 0
}

func navigate(m grid, robot xy, dirs []dir) (grid, xy) {
	cur := robot

	for _, d := range dirs {
		np := cur.add(d.xy())
		b, ok := m[np]
		if !ok {
			continue
		}
		if b.open {
			cur = np
		}
		if b.wall || b.open {
			continue
		}
		if b.box {
			var moved bool
			m, moved = moveBox(m, np, d)
			if moved {
				cur = np
			}
		}
	}
	return m, cur
}

func calcGPS(m grid) int {
	total := 0
	for pt, b := range m {
		if b.box {
			total += pt.x + pt.y*100
		}
	}
	return total
}

func moveBox(m grid, start xy, d dir) (grid, bool) {
	op := nextOpen(m, start, d)
	if op != nil {
		f, t := m[start], m[*op]
		m[start] = t
		m[*op] = f
	}
	return m, op != nil
}

func nextOpen(m grid, start xy, d dir) *xy {
	cur := start
	found := false
	for {
		cur = cur.add(d.xy())
		if b, ok := m[cur]; ok {
			if b.open {
				found = true
				break
			}
			if b.wall {
				break
			}
		} else {
			break
		}
	}
	var openspot *xy
	if found {
		openspot = &cur
	}
	return openspot
}

func double(in string) string {
	s := strings.ReplaceAll(in, "#", "##")
	s = strings.ReplaceAll(s, "O", "[]")
	s = strings.ReplaceAll(s, ".", "..")
	s = strings.ReplaceAll(s, "@", "@.")
	return s
}

func parse(in []string) (grid, xy, []dir) {
	m := make(grid)
	dirs := []dir{}
	var robot xy
	for y, line := range in {
		for x, c := range line {
			pt := xy{x, y}
			var b block
			isb := true
			switch c {
			case '#':
				b.wall = true
			case 'O':
				b.box = true
			case '.':
				b.open = true
			case '@':
				robot = pt
				b.open = true
			case '^', 'v', '<', '>':
				d := getDir(c)
				dirs = append(dirs, d)
				isb = false
			default:
				continue
			}
			if isb {
				m[pt] = b
			}
		}
	}

	return m, robot, dirs
}
