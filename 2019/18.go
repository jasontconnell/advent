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
	open bool
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
	p1 := solve(blocks, entrypoint)
	fmt.Println("Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
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

		mvs := getValidMoves(blocks, cur.position, cur)

		for _, mv := range mvs {
			cp := copyState(cur)
			vk := getVisitedKey(mv, cp.keys)

			if _, ok := visited[vk]; ok {
				continue
			}

			visited[vk] = true

			doMove := len(cur.moves)+1 < minpath

			b := blocks[mv.y][mv.x]

			if b.kind == key {
				if _, ok := cp.keys[b.ch]; !ok {
					cp.keys[b.ch] = true
					cp.opendoors[b.ch-tolower] = true
				}
			}

			if !doMove {
				continue
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

func getValidMoves(blocks [][]block, pt xy, s *state) []xy {
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
			add = s.opendoors[b.ch] || s.keys[rune(b.ch+tolower)]
		case entrance:
			add = true
		}

		if add {
			valid = append(valid, dp)
		}
	}
	return valid
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

			b := block{pt: pt, ch: rune(ch), open: bt != wall, kind: bt}
			blocks[y][x] = b
		}
	}
	return blocks, entrypoint
}
