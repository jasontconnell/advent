package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

func (p xy) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type block struct {
	pt          xy
	contents    content
	left, right *block
	up, down    *block
	sector      xy
	side        side
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

var dirval map[direction]int = map[direction]int{
	right: 0,
	down:  1,
	left:  2,
	up:    3,
}

type side int

const (
	north side = iota
	west
	east
	south
	front
	back
)

func (s side) String() string {
	return []string{"north", "west", "east", "south", "front", "back"}[s]
}

type sidefacing struct {
	side   side
	facing direction
}

var sidefacingmap map[sidefacing]direction = map[sidefacing]direction{
	{north, up}: right, {north, right}: right, {north, left}: right, {north, down}: down,
	{front, up}: up, {front, right}: up, {front, left}: down, {front, down}: down,
	{west, up}: right, {west, right}: right, {west, left}: right, {west, down}: down,
	{east, up}: up, {east, right}: left, {east, left}: left, {east, down}: left,
	{south, up}: up, {south, right}: left, {south, left}: left, {south, down}: left,
	{back, up}: up, {back, right}: up, {back, left}: down, {back, down}: down,
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 22 solution--")
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
	grid, instrs := parseInput(in)
	graphCube(grid)
	block, dir := traverseCube(grid, instrs)
	return (block.pt.y+1)*1000 + (block.pt.x+1)*4 + dirval[dir]
}

func traverseCube(grid map[xy]*block, instrs []instruction) (*block, direction) {
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
			if curtmp.sector != cur.sector {
				// switch facing

				sfkey := sidefacing{cur.side, facing}
				if dir, ok := sidefacingmap[sfkey]; ok {
					facing = dir
				}
			}
			cur = curtmp
		}

		facing = doTurn(facing, instr.dir)
	}
	return cur, facing
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

func graphCube(grid map[xy]*block) {
	size := getSize(grid)

	getGraph(grid)
	setSectors(grid, size)

	graphSide(grid, size, xy{1, 0}, xy{0, 3}, north, west, false)
	graphSide(grid, size, xy{1, 0}, xy{0, 2}, west, west, true)

	graphSide(grid, size, xy{2, 0}, xy{0, 3}, north, south, false)
	graphSide(grid, size, xy{2, 0}, xy{1, 2}, east, east, true)
	graphSide(grid, size, xy{2, 0}, xy{1, 1}, south, east, false)

	graphSide(grid, size, xy{0, 2}, xy{1, 1}, north, west, false)
	graphSide(grid, size, xy{0, 3}, xy{1, 2}, east, south, false)

	setCubeSide(grid, xy{1, 0}, north)
	setCubeSide(grid, xy{2, 0}, east)
	setCubeSide(grid, xy{1, 1}, front)
	setCubeSide(grid, xy{0, 2}, west)
	setCubeSide(grid, xy{1, 2}, south)
	setCubeSide(grid, xy{0, 3}, back)
}

func setCubeSide(grid map[xy]*block, sector xy, cubeside side) {
	for _, b := range grid {
		if b.sector == sector {
			b.side = cubeside
		}
	}
}

func graphSide(grid map[xy]*block, size int, s1, s2 xy, s1side, s2side side, invert bool) {
	s1blocks := []*block{}
	s2blocks := []*block{}

	for pt, b := range grid {
		if b.contents == blank || allJoined(b) {
			continue
		}
		if b.sector == s1 {
			if s1side == west && s1.x*size == pt.x {
				s1blocks = append(s1blocks, b)
			}
			if s1side == east && s1.x*size+size-1 == pt.x {
				s1blocks = append(s1blocks, b)
			}
			if s1side == north && s1.y*size == pt.y {
				s1blocks = append(s1blocks, b)
			}
			if s1side == south && s1.y*size+size-1 == pt.y {
				s1blocks = append(s1blocks, b)
			}
		}
		if b.sector == s2 {
			if s2side == west && s2.x*size == pt.x {
				s2blocks = append(s2blocks, b)
			}
			if s2side == east && s2.x*size+size-1 == pt.x {
				s2blocks = append(s2blocks, b)
			}
			if s2side == north && s2.y*size == pt.y {
				s2blocks = append(s2blocks, b)
			}
			if s2side == south && s2.y*size+size-1 == pt.y {
				s2blocks = append(s2blocks, b)
			}
		}
	}

	if len(s1blocks) != len(s2blocks) {
		fmt.Println("lists are uneven", len(s1blocks), len(s2blocks), s1, s2, s1side, s2side)
		panic("not like this")
	}

	sortBlocks(s1blocks)
	sortBlocks(s2blocks)

	for i := 0; i < len(s1blocks); i++ {
		start := s1blocks[i]
		dest := s2blocks[i]

		if invert {
			dest = s2blocks[len(s2blocks)-1-i]
		}

		switch s1side {
		case north:
			start.up = dest
		case south:
			start.down = dest
		case west:
			start.left = dest
		case east:
			start.right = dest
		}

		switch s2side {
		case north:
			dest.up = start
		case south:
			dest.down = start
		case west:
			dest.left = start
		case east:
			dest.right = start
		}
	}
}

func sortBlocks(list []*block) {
	sort.Slice(list, func(i, j int) bool {
		lval := list[i].pt.y*10 + list[i].pt.x
		rval := list[j].pt.y*10 + list[j].pt.x
		return lval < rval
	})
}

func setSectors(grid map[xy]*block, size int) {
	for pt := range grid {
		grid[pt].sector = xy{pt.x / size, pt.y / size}
	}
}

func getSize(grid map[xy]*block) int {
	notblank := 0
	for _, b := range grid {
		if b.contents != blank {
			notblank++
		}
	}

	return int(math.Sqrt(float64(notblank / 6)))
}

func abs(n int) int {
	return int(math.Abs(float64(n)))
}

func dist(p1, p2 xy) int {
	dx := p1.x - p2.x
	dy := p1.y - p2.y

	return abs(dx) + abs(dy)
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

func allJoined(b *block) bool {
	return b.contents != blank && b.left != nil && b.right != nil && b.up != nil && b.down != nil
}

func minmax(grid map[xy]*block) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for k, b := range grid {
		if b.contents == blank {
			continue
		}
		if k.x < min.x {
			min.x = k.x
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.y > max.y {
			max.y = k.y
		}
	}
	return min, max
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
