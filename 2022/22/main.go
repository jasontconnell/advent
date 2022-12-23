package main

import (
	"fmt"
	"log"
	"math"
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
	perimeter   bool
	opendir     direction
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
	grid, instrs := parseInput(in)
	graphCube(grid)
	return 0
	block, dir := traverse(grid, instrs)
	return (block.pt.y+1)*1000 + (block.pt.x+1)*4 + dirval[dir]
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
	min, max := minmax(grid)
	getGraph(grid)

	notblank := 0
	for _, b := range grid {
		if b.contents != blank {
			notblank++
		}
	}

	size := int(math.Sqrt(float64(notblank / 6)))
	perimeter := getPerimeter(grid)
	last := len(perimeter)

	for len(perimeter) > 0 {
		for _, pr := range perimeter {
			// fmt.Println("mapping", pr, pr.up, pr.down)
			if allJoined(pr) {
				continue
			}
			m1, m2 := mapInternalPoints(pr, perimeter, size, min, max)

			if m1 != nil && m2 != nil {
				if allJoined(m1) {
					delete(perimeter, m1.pt)
				}
				if allJoined(m2) {
					delete(perimeter, m2.pt)
				}
			}
		}
		if last == len(perimeter) {
			break
		}
		last = len(perimeter)
	}

	fmt.Println("len", len(perimeter), "mapping external")
	for len(perimeter) > 0 {
		for _, pr := range perimeter {
			// fmt.Println("mapping", pr, pr.up, pr.down)
			if allJoined(pr) {
				continue
			}
			m1, m2 := mapExternalPoints(pr, perimeter, size, min, max)
			if m1 != nil && m2 != nil {
				if allJoined(m1) {
					delete(perimeter, m1.pt)
				}
				if allJoined(m2) {
					delete(perimeter, m2.pt)
				}
				// fmt.Println("len", len(perimeter), "mapped external")
			}
		}

		if len(perimeter) == last {
			break
		}
		last = len(perimeter)
	}
	fmt.Println("perimeter points remaining", len(perimeter))
	if len(perimeter) > 0 {
		// // fmt.Println("can't find matches for", len(perimeter))
		// for _, pr := range perimeter {
		// 	fmt.Println(pr.pt, "left", pr.left, "right", pr.right, "up", pr.up, "down", pr.down)
		// }
		// panic("can't find matches for all points")
	}
	// fmt.Println(perimeter)
}

func allJoined(b *block) bool {
	return countJoined(b) == 4
}

func countJoined(b *block) int {
	count := 0
	if b.left != nil {
		count++
	}
	if b.right != nil {
		count++
	}
	if b.up != nil {
		count++
	}
	if b.down != nil {
		count++
	}
	return count
}

func mapExternalPoints(b *block, m map[xy]*block, cubesize int, min, max xy) (*block, *block) {
	var minblock *block
	var mindist int = math.MaxInt32
	for _, check := range m {
		if b.pt == check.pt || b.pt.x == check.pt.x || b.pt.y == check.pt.y {
			continue
		}

		fmt.Println(b.pt, "left", b.left, "right", b.right, "up", b.up, "down", b.down, check.pt)

		d := dist(b.pt, check.pt)
		sq := abs(b.pt.x-check.pt.y) == abs(b.pt.y-check.pt.x)
		fmt.Println(b.pt, check.pt, d, sq, d%2)
		if d < mindist && d%2 == 1 {
			mindist = d
			minblock = check
			break
		}
	}

	if minblock == nil {
		// panic("external minblock is nil " + b.String())
		return nil, nil
	}

	mapToSide(b, minblock)
	mapToSide(minblock, b)

	fmt.Println("closest point to", b, "is", minblock)
	return b, minblock
}

func mapInternalPoints(b *block, m map[xy]*block, cubesize int, min, max xy) (*block, *block) {
	if b.up == nil || b.down == nil {
		return nil, nil
	}
	var minblock *block
	var mindist int = math.MaxInt32
	for _, check := range m {
		if b.pt == check.pt {
			continue
		}
		if b.left == nil {
			if check.pt.x >= b.pt.x && (check.up != nil && check.down != nil) {
				continue
			}
		}
		if b.right == nil {
			if check.pt.x <= b.pt.x && (check.up != nil && check.down != nil) {
				continue
			}
		}

		d := dist(b.pt, check.pt)
		sq := abs(b.pt.x-check.pt.y) == abs(b.pt.y-check.pt.x)

		if d < mindist && sq {
			mindist = d
			minblock = check
			break
		}
	}

	if minblock == nil {
		// panic("minblock is nil " + b.String())
		return nil, nil
	}

	mapToSide(b, minblock)
	mapToSide(minblock, b)

	return b, minblock
}

func mapToSide(b, tomap *block) {
	if b.down == nil {
		b.down = tomap
	} else if b.up == nil {
		b.up = tomap
	} else if b.left == nil {
		b.left = tomap
	} else if b.right == nil {
		b.right = tomap
	}
}

func abs(n int) int {
	return int(math.Abs(float64(n)))
}

func dist(p1, p2 xy) int {
	dx := p1.x - p2.x
	dy := p1.y - p2.y

	return abs(dx) + abs(dy)
}

func getPerimeter(grid map[xy]*block) map[xy]*block {
	getGraph(grid) // this method could be called before it, let's make sure to graph points first
	p := make(map[xy]*block)
	for pt, b := range grid {
		if b.perimeter {
			p[pt] = b
		}
	}
	return p
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

		b.perimeter = !allJoined(b)
		if b.perimeter {
			if b.left == nil {
				b.opendir = left
			}
			if b.right == nil {
				b.opendir = right
			}
			if b.up == nil {
				b.opendir = up
			}
			if b.down == nil {
				b.opendir = down
			}
		}
	}
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
