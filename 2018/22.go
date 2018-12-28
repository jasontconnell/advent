package main

import (
	"fmt"
	"time"
)

var goal block = block{xy: xy{x: 12, y: 763}, gindex: 0}
var depth int = 7740

var debug bool = true
var dbgoal block = block{xy: xy{x: 10, y: 10}, gindex: 0}
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
	return fmt.Sprintf("%s: %s", xyt.xy, toolstr(xyt.tool))
}

type rescuer struct {
	xy
	tool int
	minutes int
}

type state struct {
	current xytool
	moves   []xytool
	minutes int
}

func main() {
	startTime := time.Now()
	if debug {
		goal = dbgoal
		depth = dbdepth
	}

	grid := makeGrid(goal.x*2, goal.y*2)
	calcTerrain(grid, goal, depth)

	p1 := sumTerrain(grid)

	res := rescuer{xy: xy{0, 0}, tool: Torch, minutes: 0} // You start at 0,0 (the mouth of the cave) with the torch equipped

	fmt.Println("Part 1:", p1)
	
	p2 := solve(res, grid, goal)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
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

func solve(res rescuer, grid [][]*block, goal block) int {
	vmap := make(map[xytool]bool)
	solves := getPath(res, goal, grid, 1, 7, vmap)
	solve := solves[len(solves)-1]
	
	for _, mv := range solve.moves {
		blk := grid[mv.y][mv.x]
		fmt.Println("Moved to", mv.xy, "with tool", toolstr(mv.tool), "through", terrstr(blk.terrain), "terrain.")
	}
	return solve.minutes
}

func getPath(res rescuer, goal block, grid [][]*block, travelTime, toolChangeTime int, visited map[xytool]bool) []state {
	goalxyt := xytool{xy: goal.xy, tool: Torch}
	curxytool := xytool{xy: res.xy, tool: res.tool}
	start := state{current: curxytool, moves: []xytool{curxytool}, minutes: 0}
	queue := []state{start}
	solves := []state{}
	minsolve := 10000

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]

		mvs := getMoves(grid, st.current.xy, goal.xy)

		for _, mv := range mvs {
			minutes := st.minutes

			if st.current.tool != mv.tool {
				minutes += toolChangeTime + travelTime
			} else {
				minutes += travelTime
			}

			if mv.x == goalxyt.x && mv.y == goalxyt.y {
				// if mv.tool != goalxyt.tool {
				// 	mv.tool = goalxyt.tool
				// 	minutes += toolChangeTime
				// }

				if minutes < minsolve {
					fmt.Println(minutes)
					minsolve = minutes
					solves = append(solves, state{current: mv, moves: append(st.moves, mv), minutes: minutes})
				}
			}

			if _, ok := visited[mv]; !ok {
				visited[mv] = true

				if minutes+travelTime < minsolve { // assume just travel time
					queue = append(queue, state{current: mv, moves: append(st.moves, mv), minutes: minutes})
				}
			}
		}
	}

	return solves
}

func getMoves(grid [][]*block, start xy, goal xy) []xytool {
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

		tools := []int{}
		if tool == ClimbingGear|Torch {
			tools = append(tools, ClimbingGear)
			tools = append(tools, Torch)
		} else {
			tools = append(tools, tool)
			tools = append(tools, None)
		}

		for _, tl := range tools {
			res := xytool{xy: pt, tool: tl}
			mvs = append(mvs, res)
		}
	}

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

func calcTerrain(grid [][]*block, goal block, depth int) {
	for y := 0; y < goal.y+1; y++ {
		for x := 0; x < goal.x+1; x++ {
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

func geologicIndex(grid [][]*block, pt *block, goal block, depth int) int {
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

func erosionLevel(grid [][]*block, pt *block, goal block, depth int) int {
	if pt.erosioncheck {
		return pt.erosion
	}
	gindex := geologicIndex(grid, pt, goal, depth)

	pt.erosion = (gindex + depth) % 20183
	pt.erosioncheck = true
	pt.terrain = pt.erosion % 3

	return pt.erosion
}
