package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

var input = "18.txt"

const tolower rune = 32

type blocktype int

const (
	wall blocktype = iota
	key
	door
	entrance
	path
)

type xy struct {
	x, y int
}

var deltas []xy = []xy{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
}

type block struct {
	pt   xy
	ch   rune
	kind blocktype
}

type state struct {
	moves     []xy
	keys      map[rune]bool
	opendoors map[rune]bool
	position  xy
}

type vkey struct {
	xy
	keys string
}

type queue []*state

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	blocks, entrypoint := readGrid(lines)
	p1 := 0 //solve(blocks, entrypoint)
	fmt.Println("Part 1:", p1)

	b := blocks[entrypoint.y][entrypoint.x]
	b.ch = '#'
	b.kind = wall
	blocks[entrypoint.y][entrypoint.x] = b

	for _, d := range deltas {
		p := xy{entrypoint.x + d.x, entrypoint.y + d.y}
		dp := blocks[p.y][p.x]
		dp.ch = '#'
		dp.kind = wall
		blocks[p.y][p.x] = dp
	}

	robots := []xy{
		{entrypoint.x + 1, entrypoint.y + 1},
		{entrypoint.x + 1, entrypoint.y - 1},
		{entrypoint.x - 1, entrypoint.y - 1},
		{entrypoint.x - 1, entrypoint.y + 1},
	}
	for _, pt := range robots {
		blocks[pt.y][pt.x].kind = entrance
	}

	p2 := fourRobotSolve(blocks, robots)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func fourRobotSolve(blocks [][]block, entrypoints []xy) int {
	allkeys := getKeyMap(blocks)
	foundKeys := make(map[rune]bool)
	visited := make(map[vkey]bool)

	robot := 0
	robots := len(entrypoints)
	rstates := []*state{}
	for _, pt := range entrypoints {
		rstates = append(rstates, newState(pt))
	}

	rq := make([]queue, robots)
	for i, s := range rstates {
		rq[i] = append(rq[i], s)
	}

	minpaths := []int{}
	for range rstates {
		minpaths = append(minpaths, 100000)
	}

	goalReached := make([]bool, robots)
	goalStates := make([]*state, robots)

	alllen := 4
	for alllen > 0 {
		if len(rq[robot]) == 0 {
			robot = (robot + 1) % robots
			continue
		}
		cur := rq[robot][0]
		rq[robot] = rq[robot][1:]
		alllen--

		if len(cur.moves) >= minpaths[robot] {
			continue
		}

		curb := blocks[cur.position.y][cur.position.x]
		if _, ok := foundKeys[curb.ch+tolower]; curb.kind == door && !ok {
			fmt.Println(robot, "door. waiting for", string(curb.ch+tolower))
			rq[robot] = append(rq[robot], cur)
			alllen++
			robot = (robot + 1) % robots
			continue
		}

		mvs := getValidMoves(blocks, cur.position, cur, false)
		// fmt.Println(mvs)

		for _, mv := range mvs {
			cp := copyState(cur)
			vk := getVisitedKey(mv, cp.keys)

			if _, ok := visited[vk]; ok {
				continue
			}

			// fmt.Println(robot, "visiting", mv, len(cp.moves), "found keys", len(foundKeys))
			visited[vk] = true

			b := blocks[mv.y][mv.x]

			nextRobot := false

			if b.kind == key {
				if _, ok := cp.keys[b.ch]; !ok {
					fmt.Println("found a key", string(b.ch))
					cp.keys[b.ch] = true
					cp.opendoors[b.ch-tolower] = true
					foundKeys[b.ch] = true
					// nextRobot = true
				}
			} else if b.kind == door {
				fmt.Println("came across a door", string(b.ch))
				if _, ok := foundKeys[b.ch+tolower]; !ok {
					nextRobot = true
					fmt.Println("no key for", string(b.ch))
				}
			}

			if !nextRobot {
				cp.position = mv
				cp.moves = append(cp.moves, mv)
				goalReached[robot] = len(foundKeys) == len(allkeys)

				if goalReached[robot] {
					if len(cp.moves) < minpaths[robot] {
						goalStates[robot] = cp
						minpaths[robot] = len(cp.moves)

						fmt.Println("new goal reached. minpath", minpaths[robot])
					}
				} else {
					rq[robot] = append(rq[robot], cp)
					alllen++
				}
			} else {
				rq[robot] = append(rq[robot], cp)
				alllen++
				robot = (robot + 1) % robots
			}
		}
	}

	x := 0
	for _, gs := range goalStates {
		if gs == nil {
			continue
		}
		x += len(gs.moves)
	}
	return x
}

func solve(blocks [][]block, entrypoint xy) int {
	allkeys := getKeyMap(blocks)
	s := newState(entrypoint)

	minpath := 1000000
	q := queue{s}
	goalReached := false
	goalState := &state{}

	visited := make(map[vkey]bool)

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if len(cur.moves) >= minpath {
			continue
		}

		mvs := getValidMoves(blocks, cur.position, cur, true)

		for _, mv := range mvs {
			cp := copyState(cur)
			vk := getVisitedKey(mv, cp.keys)

			if _, ok := visited[vk]; ok {
				continue
			}

			visited[vk] = true

			b := blocks[mv.y][mv.x]

			if b.kind == key {
				if _, ok := cp.keys[b.ch]; !ok {
					cp.keys[b.ch] = true
					cp.opendoors[b.ch-tolower] = true
				}
			}

			cp.position = mv
			cp.moves = append(cp.moves, mv)
			goalReached = len(cp.keys) == len(allkeys)

			if goalReached {
				if len(cp.moves) < minpath {
					goalState = cp
					minpath = len(cp.moves)

					fmt.Println("new goal reached. minpath", minpath, len(q))
				}
			} else {
				q = append(q, cp)
			}
		}
	}

	return len(goalState.moves)
}

func getVisitedKey(pt xy, m map[rune]bool) vkey {
	s := []int{}
	for k, _ := range m {
		s = append(s, int(k))
	}

	sort.Ints(s)

	ret := ""
	for _, ch := range s {
		ret += string(rune(ch))
	}
	return vkey{xy: pt, keys: ret}
}

func newState(pos xy) *state {
	s := &state{position: pos}
	s.keys = make(map[rune]bool)
	s.opendoors = make(map[rune]bool)
	return s
}

func copyState(s *state) *state {
	sc := newState(s.position)
	xyc := make([]xy, len(s.moves))
	copy(xyc, s.moves)
	sc.moves = xyc

	for k, v := range s.keys {
		sc.keys[k] = v
	}

	for k, v := range s.opendoors {
		sc.opendoors[k] = v
	}

	return sc
}

func loopBlocks(blocks [][]block, f func(y, x int, b block)) {
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			b := blocks[y][x]
			f(y, x, b)
		}
	}
}

func getKeyMap(blocks [][]block) map[xy]rune {
	km := make(map[xy]rune)
	loopBlocks(blocks, func(y, x int, b block) {
		if b.kind == key {
			km[xy{x, y}] = b.ch
		}
	})
	return km
}

func getDoorMap(blocks [][]block) map[xy]rune {
	dm := make(map[xy]rune)
	loopBlocks(blocks, func(y, x int, b block) {
		if b.kind == door {
			dm[xy{x, y}] = b.ch
		}
	})
	return dm
}

func getValidMoves(blocks [][]block, pt xy, s *state, keycheck bool) []xy {
	maxx, maxy := len(blocks[0]), len(blocks)
	valid := []xy{}
	for _, d := range deltas {
		dp := xy{d.x + pt.x, d.y + pt.y}

		if dp.x < 0 || dp.y < 0 || dp.x >= maxx || dp.y >= maxy {
			continue
		}

		add := false
		b := blocks[dp.y][dp.x]
		switch b.kind {
		case path, key:
			add = true
		case door:
			if keycheck {
				add = s.opendoors[b.ch] || s.keys[rune(b.ch+tolower)]
			} else {
				add = true
			}
		case entrance:
			add = true
		}

		if add {
			valid = append(valid, dp)
		}
	}
	return valid
}

func drawGrid(blocks [][]block) {
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			ch := "#"
			b := blocks[y][x]

			if b.kind == key {
				ch = string(b.ch)
			} else if b.kind == entrance {
				ch = "@"
			} else if b.kind == path {
				ch = "."
			}

			fmt.Print(ch)
		}
		fmt.Println()
	}
	fmt.Println()
}

func readGrid(lines []string) ([][]block, xy) {
	blocks := make([][]block, len(lines))
	var entrypoint xy
	for y := 0; y < len(lines); y++ {
		blocks[y] = make([]block, len(lines[y]))
		for x := 0; x < len(lines[y]); x++ {
			pt := xy{x, y}
			ch := lines[y][x]
			bt := path

			if ch >= 97 && ch <= 122 {
				bt = key
			} else if ch >= 65 && ch <= 90 {
				bt = door
			} else if ch == 64 {
				bt = entrance
				entrypoint = xy{x, y}
			} else if ch == 35 {
				bt = wall
			}

			b := block{pt: pt, ch: rune(ch), kind: bt}
			blocks[y][x] = b
		}
	}
	return blocks, entrypoint
}
