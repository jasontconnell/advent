package main

import (
	"fmt"
	"container/heap"
	"math"
	"time"
)

var goal xytool = xytool{xy: xy{x: 12, y: 763}, tool: Torch}
var depth int = 7740

var debug bool = true
var dbgoal xytool = xytool{xy: xy{x: 10, y: 10}, tool: Torch}
var dbdepth int = 510

const (
	Rocky int = iota
	Wet
	Narrow
)

type xy struct {
	x, y int
}

func (pt xy) String() string {
	return fmt.Sprintf("(%d, %d)", pt.x, pt.y)
}

type block struct {
	xy
	gindex       int
	erosion      int
	erosioncheck bool
	gindexcheck  bool
	terrain      int
}

const (
	ClimbingGear int = 1 // 01
	Torch        int = 2 // 10
	None         int = 0 // 00
)

type xytool struct {
	xy
	tool int
}

func (xyt xytool) String() string {
	return fmt.Sprintf("(%d, %d, %s)", xyt.x, xyt.y, toolstr(xyt.tool))
}

type rescuer struct {
	xy
	tool    int
	minutes int
}

func (r rescuer) String() string {
	return fmt.Sprintf("rescuer at (%d, %d) with %s after %d minutes", r.x, r.y, toolstr(r.tool), r.minutes)
}

type state struct {
	current xytool
	moves   []xytool
	// current xy
	// moves   []xy
	minutes int
	steps   int
	dist    int
	// tool    int
}

func main() {
	startTime := time.Now()
	if debug {
		goal = dbgoal
		depth = dbdepth
	}

	pad := 40
	grid := makeGrid(goal.x+pad, goal.y+pad)
	calcTerrain(grid, goal, depth, pad)

	p1 := sumTerrain(grid)
	draw(grid)

	res := rescuer{xy: xy{0, 0}, tool: Torch, minutes: 0} // You start at 0,0 (the mouth of the cave) with the torch equipped

	fmt.Println("Part 1:", p1)

	p2 := solve(res, grid, goal)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func draw(grid [][]*block){
	for y := 0; y < goal.y+20; y++ {
		line := ""
		for x := 0; x < goal.x+20; x++ {
			var c rune
			switch grid[y][x].terrain {
			case Rocky: c = '.'
			case Wet: c = '='
			case Narrow: c = '|'
			}
			line += string(c)
		}
		fmt.Println(line)
	}
}

func toolstr(tool int) string {
	s := ""
	switch tool {
	case ClimbingGear | Torch:
		s = "climbing gear or torch"
	case Torch:
		s = "torch"
	case ClimbingGear:
		s = "climbing gear"
	case None:
		s = "nothing"
	}

	return s
}

func terrstr(terrain int) string {
	s := ""
	switch terrain {
	case Rocky:
		s = "rocky"
	case Wet:
		s = "Wet"
	case Narrow:
		s = "Narrow"
	}
	return s
}

func solve(res rescuer, grid [][]*block, goal xytool) int {
	return getPath(res, goal, grid, 1, 7)
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func distance(p1, p2 xy) int {
	dx := abs(p1.x - p2.x)
	dy := abs(p1.y - p2.y)
	return dx + dy
}

type item struct {
	pos xy
	tool int
	minutes int
	index int
	moves []xytool
}

type pqueue []*item

func (pq pqueue) Len() int { return len(pq) }

func (pq pqueue) Less(i, j int) bool {
	return pq[i].minutes < pq[j].minutes
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	itm := old[n-1]
	itm.index = -1
	*pq = old[0 : n-1]
	return itm
}

func (pq *pqueue) Push(x interface{}) {
	n := len(*pq)
	itm := x.(*item)
	itm.index = n
	*pq = append(*pq, itm)
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func getPath(res rescuer, goal xytool, grid [][]*block, travelTime, toolChangeTime int) int {
	start := xytool{xy: res.xy, tool: res.tool}
	bail := 8
	queue := pqueue{
		&item{pos: start.xy, tool: start.tool, minutes: 0, index: 0, moves: []xytool{start}},
	}
	heap.Init(&queue)

	distances := make(map[xytool]int)
	distances[start] = 0

	for queue.Len() > 0 {
		pt := (heap.Pop(&queue)).(*item)
		ptx := xytool{xy: pt.pos, tool: pt.tool}

		if pt.pos == goal.xy && pt.tool == goal.tool {
			fmt.Println("solved!", pt.moves)
			return pt.minutes
		}

		if pt.pos.x > bail * goal.x || pt.pos.y > bail * goal.y {
			continue
		}

		if minutes, ok := distances[ptx]; ok && minutes < pt.minutes {
			continue
		}

		mvs := getMoves(grid, ptx.xy, ptx.tool)
		fmt.Println("from", ptx, "got moves", mvs)
		for _, mv := range mvs {
			m := 1
			if mv.tool != ptx.tool {
				m = 8 // move plus swap
			}
			if minutes, ok := distances[mv]; !ok || pt.minutes + m < minutes {
				distances[mv] = pt.minutes + m
				cp := make([]xytool, len(pt.moves))
				copy(cp, pt.moves)
				cp = append(cp, mv)
				fmt.Println("trail from", ptx, cp)
				heap.Push(&queue, &item{pos: mv.xy, minutes: pt.minutes+m, tool: mv.tool, moves: cp})
			}
		}
	}

	return 0
}

func getMoves(grid [][]*block, start xy, equipped int) []xytool {
	maxy := len(grid)
	maxx := len(grid[0])
	mvs := []xytool{}
	for _, mv := range []xy{
		xy{x: 1, y: 0},
		xy{x: 0, y: 1},
		xy{x: -1, y: 0},
		xy{x: 0, y: -1},
	} {
		pt := xy{x: start.x + mv.x, y: start.y + mv.y}
		if pt.x == -1 || pt.x >= maxx || pt.y == -1 || pt.y >= maxy {
			continue
		}
		tool := validTools(grid, pt)

		if tool == equipped || (tool&equipped) == equipped {
			xyt := xytool{xy: pt, tool: equipped}
			mvs = append(mvs, xyt)
		}
	}

	v := validTools(grid, start)
	swtool := xytool{xy: start, tool: v ^ equipped}
	// if v == ClimbingGear|Torch {
	// 	swtool.tool = ClimbingGear
	// 	if equipped == ClimbingGear {
	// 		swtool.tool = Torch
	// 	}
	// } else if equipped != None {
	// 	swtool.tool = None
	// }

	mvs = append(mvs, swtool)

	return mvs
}

func sumTerrain(grid [][]*block) int {
	sum := 0
	for y := 0; y < goal.y+1; y++ {
		for x := 0; x < goal.x+1; x++ {
			sum += grid[y][x].terrain
		}
	}
	return sum
}

func validTools(grid [][]*block, pt xy) int {
	terrain := grid[pt.y][pt.x].terrain
	tools := 0

	switch terrain {
	case Rocky:
		tools = ClimbingGear | Torch
	case Wet:
		tools = ClimbingGear | None
	case Narrow:
		tools = Torch | None
	}

	return tools
}

func calcTerrain(grid [][]*block, goal xytool, depth, pad int) {
	for y := 0; y < goal.y+pad; y++ {
		for x := 0; x < goal.x+pad; x++ {
			geologicIndex(grid, grid[y][x], goal, depth)
			erosionLevel(grid, grid[y][x], goal, depth)
		}
	}
}

func makeGrid(x, y int) [][]*block {
	grid := make([][]*block, y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			grid[i] = append(grid[i], &block{xy: xy{x: j, y: i}})
		}
	}
	return grid
}

func geologicIndex(grid [][]*block, pt *block, goal xytool, depth int) int {
	if pt.gindex != 0 {
		return pt.gindex
	}
	pt.gindexcheck = true

	if pt.x == 0 && pt.y == 0 || pt.x == goal.x && pt.y == goal.y {
		pt.gindex = 0
		return pt.gindex
	}

	gindex := 0
	if pt.y == 0 {
		gindex = pt.x * 16807
	} else if pt.x == 0 {
		gindex = pt.y * 48271
	} else {
		left := erosionLevel(grid, grid[pt.y][pt.x-1], goal, depth)
		down := erosionLevel(grid, grid[pt.y-1][pt.x], goal, depth)

		gindex = left * down
	}
	pt.gindex = gindex

	return pt.gindex
}

func erosionLevel(grid [][]*block, pt *block, goal xytool, depth int) int {
	if pt.erosioncheck {
		return pt.erosion
	}
	gindex := geologicIndex(grid, pt, goal, depth)

	pt.erosion = (gindex + depth) % 20183
	pt.erosioncheck = true
	pt.terrain = pt.erosion % 3

	return pt.erosion
}
