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
	pt, dir := traverse(grid, instrs)
	return (pt.y+1)*1000 + (pt.x+1)*4 + dirval[dir]
}

func part2(in input) output {
	grid, instrs := parseInput(in)
	pt, dir := traverseCube(grid, instrs)
	return pt.y*1000 + pt.x*4 + dirval[dir]
}

func traverseCube(grid map[xy]block, instrs []instruction) (xy, direction) {
	facing := right
	start := xy{0, 0}
	for grid[start].contents != space && grid[start].contents != wall {
		start.x++
	}
	cur := start

	getCube(grid)

	return cur, facing
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

func getCube(grid map[xy]block) []map[xy]block {
	sides := make([]map[xy]block, 6)
	min, max := minmax(grid)
	sectors := make(map[xy]int)
	occupied := make(map[int]xy)

	notblank := 0
	for _, b := range grid {
		if b.contents != blank {
			notblank++
		}
	}

	for i := 0; i < 6; i++ {
		sides[i] = make(map[xy]block)
	}

	size := int(math.Sqrt(float64(notblank / 6)))
	perrow := max.x/size + 1
	ys := 0
	for y := min.y; y <= max.y; y++ {
		ys = y / size
		for x := min.x; x <= max.x; x++ {
			pt := xy{x, y}
			sector := (ys*perrow + x/size)
			sxy := xy{x / size, y / size}
			if b, ok := grid[pt]; ok && b.contents != blank {
				occupied[sector] = sxy
			}
			sectors[pt] = sector
		}
	}

	// normalize to top bottom east west front and back

	keys := sortXYVals(occupied)
	result := make(map[xy]side)
	for _, pt := range occupied {
		result[pt] = tbd
	}
	start := keys[0]
	result[start] = top
	visit := make(map[xy]bool)
	queue := []cubestate{
		{pt: xy{start.x - 1, start.y}, prev: start, prevside: top, dir: left},
		{pt: xy{start.x + 1, start.y}, prev: start, prevside: top, dir: right},
		{pt: xy{start.x, start.y + 1}, prev: start, prevside: top, dir: down},
		{pt: xy{start.x, start.y - 1}, prev: start, prevside: top, dir: up},
	}
	fmt.Println(keys)

	avail := make(map[side]bool)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// fmt.Println(start, result[start], xy{start.x + 1, start.y}, result[xy{start.x + 1, start.y}], cur.pt)

		// fmt.Println(cur, result[cur.pt])

		if _, ok := visit[cur.pt]; ok {
			continue
		}
		visit[cur.pt] = true

		if s, ok := result[cur.pt]; !ok || s != tbd {
			continue
		}

		var newside side

		switch cur.dir {
		case left:
			switch cur.prevside {
			case top, bottom, front:
				newside = west
			case west:
				newside = back
			case east:
				newside = front
			case back:
				newside = east
			}
		case right:
			switch cur.prevside {
			case top, bottom, front:
				newside = east
			case west:
				newside = front
			case east:
				newside = back
			case back:
				newside = west
			}
		case down:
			switch cur.prevside {
			case top:
				newside = front
			case east, west, front:
				if _, ok := avail[bottom]; !ok {
					newside = bottom
				} else {
					newside = back
				}
			case bottom:
				newside = back
			case back:
				newside = top
			}
		case up:
			switch cur.prevside {
			case top:
				newside = back
			case east, west, front:
				newside = top
			case bottom:
				newside = front
			case back:
				newside = bottom
			}
		}

		if _, ok := avail[newside]; ok {
			fmt.Println("going", cur.dir, "to", cur.pt, "from", cur.prevside, cur.prev, "result", newside)
			fmt.Println(result)
			fmt.Println(avail)
			fmt.Println(cur.dir, "from", cur.pt, "prevside", cur.prevside, "from", cur.prev)
			panic("already got side")
		}
		if cs, ok := result[cur.pt]; ok && cs == tbd {
			fmt.Println("going", cur.dir, "to", cur.pt, "from", cur.prevside, cur.prev, "result", newside)
			result[cur.pt] = newside
			avail[newside] = true
		}

		mvs := []cubestate{
			{prev: cur.pt, prevside: newside, pt: xy{cur.pt.x + 1, cur.pt.y}, dir: right},
			{prev: cur.pt, prevside: newside, pt: xy{cur.pt.x - 1, cur.pt.y}, dir: left},
			{prev: cur.pt, prevside: newside, pt: xy{cur.pt.x, cur.pt.y - 1}, dir: up}, // should never be up
			{prev: cur.pt, prevside: newside, pt: xy{cur.pt.x, cur.pt.y + 1}, dir: down},
		}
		queue = append(queue, mvs...)
	}

	fmt.Println(result)

	return sides
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

func minmax(grid map[xy]block) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for k := range grid {
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

func sortXYVals[K comparable](m map[K]xy) []xy {
	list := []xy{}
	for _, v := range m {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		lval := list[i].y*10 + list[i].x
		rval := list[j].y*10 + list[j].x
		return lval < rval
	})
	return list
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
			m[xy{x, y}] = b
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
