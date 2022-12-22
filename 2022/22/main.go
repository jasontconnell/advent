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
	pt          xy
	contents    content
	left, right *block
	up, down    *block
}

func (b block) String() string {
	return fmt.Sprintf("block: %v %v", b.pt, b.contents)
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

type cubestate struct {
	pt       xy
	prev     xy
	prevside side
	dir      direction
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

type side int

const (
	tbd    side = -1
	top    side = 0
	bottom side = 1
	west   side = 2
	east   side = 3
	front  side = 4
	back   side = 5
)

var sides []string = []string{"top", "bottom", "west", "east", "front", "back"}

func (s side) String() string {
	if s == tbd {
		return "tbd"
	}
	return sides[s]
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
	graphSides(grid)
	block, dir := traverse(grid, instrs)
	return (block.pt.y+1)*1000 + (block.pt.x+1)*4 + dirval[dir]
}

func part2(in input) output {
	return 0
}

func traverse(grid map[xy]*block, instrs []instruction) (*block, direction) {
	facing := right
	start := xy{0, 0}
	for grid[start].contents != space && grid[start].contents != wall {
		start.x++
	}

	cur := grid[start]
	for _, instr := range instrs {
		for i := 0; i < instr.count; i++ {
			var curtmp *block
			switch facing {
			case left:
				curtmp = cur.left
			case right:
				curtmp = cur.right
			case up:
				curtmp = cur.up
			case down:
				curtmp = cur.down
			}
			if curtmp.contents == wall {
				break
			}
			cur = curtmp
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

func graphSides(grid map[xy]*block) {
	getGraph(grid)
	for _, b := range grid {
		if b.contents == blank {
			continue
		}
		if b.left == nil {
			b.left = getOppositePoint(b, right)
			b.left.right = b
		}
		if b.right == nil {
			b.right = getOppositePoint(b, left)
			b.right.left = b
		}
		if b.up == nil {
			b.up = getOppositePoint(b, down)
			b.up.down = b
		}
		if b.down == nil {
			b.down = getOppositePoint(b, up)
			b.down.up = b
		}
	}
}

func getOppositePoint(b *block, dir direction) *block {
	opp := b
	switch dir {
	case left:
		for opp.left != nil {
			opp = opp.left
		}
	case right:
		for opp.right != nil {
			opp = opp.right
		}

	case down:
		for opp.down != nil {
			opp = opp.down
		}
	case up:
		for opp.up != nil {
			opp = opp.up
		}
	}
	return opp
}

func getGraph(grid map[xy]*block) {
	for pt, b := range grid {
		if b.contents == blank {
			continue
		}
		lpt := xy{pt.x - 1, pt.y}
		rpt := xy{pt.x + 1, pt.y}
		dpt := xy{pt.x, pt.y + 1}
		upt := xy{pt.x, pt.y - 1}

		if n, ok := grid[lpt]; ok && b.left == nil && n.contents != blank {
			b.left = n
			n.right = b
		}

		if n, ok := grid[rpt]; ok && b.right == nil && n.contents != blank {
			b.right = n
			n.left = b
		}

		if n, ok := grid[upt]; ok && b.up == nil && n.contents != blank {
			b.up = n
			n.down = b
		}

		if n, ok := grid[dpt]; ok && b.down == nil && n.contents != blank {
			b.down = n
			n.up = b
		}
	}
}

func parseInput(in input) (map[xy]*block, []instruction) {
	m := make(map[xy]*block)
	instrs := []instruction{}
	instridx := 0

	for y, line := range in {
		if line == "" {
			instridx = y
			break
		}

		for x, c := range line {
			b := &block{pt: xy{x, y}, contents: content(c)}
			m[b.pt] = b
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
