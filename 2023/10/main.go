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

func (pt xy) add(pt2 xy) xy {
	return xy{pt.x + pt2.x, pt.y + pt2.y}
}

func (pt xy) sub(pt2 xy) xy {
	return xy{pt2.x - pt.x, pt2.y - pt.y}
}

func (pt xy) opposite() xy {
	p2 := xy{-pt.x, -pt.y}
	return p2
}

func (pt xy) eq(p2 xy) bool {
	return pt.x == p2.x && pt.y == p2.y
}

type pipe struct {
	pt    xy
	ch    rune
	dir   xy
	indir xy
}

func (p pipe) String() string {
	return fmt.Sprintf("(%d,%d) %c [dir(%d,%d) out(%d,%d)]", p.pt.x, p.pt.y, p.ch, p.dir.x, p.dir.y, p.indir.x, p.indir.y)
}

func maxes(m map[xy]pipe) (int, int) {
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

func print(m map[xy]pipe, plookup map[xy]pipe, internal map[xy]bool) {
	mx, my := maxes(m)
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			pt := xy{x, y}

			c := ' '
			if _, ok := internal[pt]; ok {
				c = '.'
			}
			if p, inloop := plookup[pt]; inloop {
				c = p.ch
			}
			fmt.Print(string(c))
		}
		fmt.Print("\n")
	}
}

var dirs []xy = []xy{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

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
	start, pipes := parseInput(in)
	list := getLoopPoints(start, pipes)
	return getEnclosedArea(pipes, list)
}

func traverse(start pipe, pipes map[xy]pipe) int {
	list := getLoopPoints(start, pipes)
	return len(list) / 2
}

func getEnclosedArea(pipes map[xy]pipe, loop []pipe) int {
	area := 0
	for i := 0; i < len(loop); i++ {
		cur := loop[i]
		next := loop[(i+1)%len(loop)]

		area += cur.pt.y*cur.pt.x - cur.pt.x*next.pt.y
	}
	return area - len(loop)/2 + 1
}

func getLoopPoints(start pipe, pipes map[xy]pipe) []pipe {
	list := []pipe{}
	cur := start.pt
	dir := xy{0, 0}
	done := false
	for !done {
		p, ok := pipes[cur]
		if !ok {
			break
		}

		var np xy
		switch p.ch {
		case '|':
			if dir.eq(xy{0, 1}) || dir.eq(xy{0, -1}) {
				np = cur.add(dir) // keep going in same dir
			}
		case '-':
			if dir.eq(xy{-1, 0}) || dir.eq(xy{1, 0}) {
				np = cur.add(dir)
			}
		case 'L':
			if dir.eq(xy{0, 1}) {
				dir = xy{1, 0}
			} else if dir.eq(xy{-1, 0}) {
				dir = xy{0, -1}
			}
			np = cur.add(dir)
		case 'J':
			if dir.eq(xy{0, 1}) {
				dir = xy{-1, 0}
			} else if dir.eq(xy{1, 0}) {
				dir = xy{0, -1}
			}
			np = cur.add(dir)
		case '7':
			if dir.eq(xy{0, -1}) {
				dir = xy{-1, 0}
			} else if dir.eq(xy{1, 0}) {
				dir = xy{0, 1}
			}
			np = cur.add(dir)
		case 'F':
			if dir.eq(xy{0, -1}) {
				dir = xy{1, 0}
			} else if dir.eq(xy{-1, 0}) {
				dir = xy{0, 1}
			}
			np = cur.add(dir)
		case 'S':
			dir = xy{0, 0}
			for _, d := range dirs {
				cp := p.pt.add(d)
				dest := pipes[cp]
				// first one wins
				if d.eq(xy{0, 1}) && (dest.ch == '|' || dest.ch == 'J' || dest.ch == 'L') {
					dir = d
					break
				} else if d.eq(xy{0, -1}) && (dest.ch == '|' || dest.ch == 'F' || dest.ch == '7') {
					dir = d
					break
				} else if d.eq(xy{-1, 0}) && (dest.ch == '-' || dest.ch == 'F' || dest.ch == 'L') {
					dir = d
					break
				} else if d.eq(xy{1, 0}) && (dest.ch == '-' || dest.ch == 'J' || dest.ch == '7') {
					dir = d
					break
				}
			}
			np = cur.add(dir)
		}

		p.dir = dir
		list = append(list, p)
		cur = np

		done = cur == start.pt
	}
	return list // list will have start at the front and back
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
