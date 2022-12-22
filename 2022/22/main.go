package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type content rune
type direction rune

const (
	blank content = ' '
	space content = '.'
	wall  content = '#'
)

func (c content) String() string {
	return string(c)
}

const (
	none  direction = 0
	left  direction = 'L'
	right direction = 'R'
	up    direction = 'U'
	down  direction = 'D'
)

func (d direction) String() string {
	return string(d)
}

type instruction struct {
	count int
	dir   direction
}

type xy struct {
	x, y int
}

type block struct {
	contents content
}

func (b block) String() string {
	return fmt.Sprintf("block: %v", b.contents)
}

var delta map[direction]xy = map[direction]xy{
	up:    {0, -1},
	down:  {0, 1},
	left:  {-1, 0},
	right: {1, 0},
}

type turnkey struct {
	facing, turn direction
}

var turn map[turnkey]direction = map[turnkey]direction{
	{up, right}:    right,
	{up, left}:     left,
	{down, right}:  left,
	{down, left}:   right,
	{right, right}: down,
	{right, left}:  up,
	{left, right}:  up,
	{left, left}:   down,
}

var dirval map[direction]int = map[direction]int{
	right: 0,
	down:  1,
	left:  2,
	up:    3,
}

var opposite map[direction]direction = map[direction]direction{
	left:  right,
	right: left,
	up:    down,
	down:  up,
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	grid, instrs := parseInput(in)
	pt, dir := traverse(grid, instrs)
	return pt.y*1000 + pt.x*4 + dirval[dir]
}

func part2(in input) output {
	return 0
}

func traverse(grid map[xy]block, instrs []instruction) (xy, direction) {

	facing := right
	start := xy{1, 1}
	for grid[start].contents != space && grid[start].contents != wall {
		start.x++
	}

	cur := start
	for _, instr := range instrs {
		d := delta[facing]

		for i := 0; i < instr.count; i++ {
			curtmp := xy{cur.x + d.x, cur.y + d.y}
			if b, ok := grid[curtmp]; !ok || b.contents == blank {
				curtmp = getOpposite(grid, cur, facing)
			}
			b := grid[curtmp]
			if b.contents == space {
				cur = curtmp
			} else {
				break
			}
		}

		facing = doTurn(facing, instr.dir)
	}
	return cur, facing
}

func doTurn(facing direction, t direction) direction {
	if t == none {
		return facing
	}
	return turn[turnkey{facing, t}]
}

func getOpposite(grid map[xy]block, pos xy, facing direction) xy {
	opp := opposite[facing]
	d := delta[opp]

	cur := pos
	tmp := cur
	done := false
	for !done {
		tmp.x += d.x
		tmp.y += d.y

		if _, ok := grid[tmp]; ok {
			b := grid[tmp]
			if b.contents != blank {
				cur = tmp
			}
		} else {
			done = true
		}
	}
	return cur
}

func parseInput(in input) (map[xy]block, []instruction) {
	m := make(map[xy]block)
	instrs := []instruction{}
	instridx := 0

	for y, line := range in {
		if line == "" {
			instridx = y
			break
		}

		for x, c := range line {
			b := block{contents: content(c)}
			m[xy{x + 1, y + 1}] = b
		}
	}

	for i := instridx; i < len(in); i++ {
		line := in[i]
		if line == "" {
			continue
		}

		rep := line
		rep = strings.Replace(rep, "L", " L ", -1)
		rep = strings.Replace(rep, "R", " R ", -1)
		sp := strings.Fields(rep)

		for _, s := range sp {
			n, err := strconv.Atoi(s)
			var instr instruction
			if err != nil {
				instr = instruction{count: 0, dir: direction(s[0])}
			} else {
				instr = instruction{count: n, dir: none}
			}
			instrs = append(instrs, instr)
		}
	}

	return m, instrs
}
