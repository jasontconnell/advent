package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "20.txt"

type block struct {
	id        int
	grid      [][]bool
	neighbors map[edge]*block
	oriented  bool
}

func (b *block) String() string {
	return fmt.Sprintf("id: %d", b.id)
}

func newBlock(id, d int) *block {
	b := &block{id: id}
	b.neighbors = make(map[edge]*block)
	b.grid = make([][]bool, d)
	return b
}

type neighbor struct {
	block *block
	side  edge
}

func (n neighbor) String() string {
	return fmt.Sprintf("[id: %d]", n.block.id)
}

type xy struct {
	x, y int
}

type edge int

const (
	undefined edge = iota
	top
	bottom
	left
	right
)

func (e edge) String() string {
	s := ""
	switch e {
	case top:
		s = "top"
	case bottom:
		s = "bottom"
	case left:
		s = "left"
	case right:
		s = "right"
	case undefined:
		s = "nil side"
	}
	return s
}

func (e edge) opposite() edge {
	switch e {
	case top:
		return bottom
	case bottom:
		return top
	case left:
		return right
	case right:
		return left
	}
	return undefined
}

var seaMonster []xy = []xy{
	{0, 0}, {1, 1}, {4, 1}, {5, 0}, {6, 0}, {7, 1}, {10, 1}, {11, 0},
	{12, 0}, {13, 1}, {16, 1}, {17, 0}, {18, 0}, {19, 0}, {18, -1},
}

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
	d := 10

	ms := readGrids(lines, d)
	p1 := solve(ms, d)

	fmt.Println("Part 1:", p1)

	joined := joinGrids(ms, d)
	mcount := findSeaMonsters(joined)

	smcount := len(seaMonster) * mcount

	tcount := 0
	for _, r := range joined {
		for _, c := range r {
			if c {
				tcount++
			}
		}
	}

	p2 := tcount - smcount
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func findSeaMonsters(grid [][]bool) int {
	c := 0
	rotates := 0
	for c == 0 {
		for y := 0; y < len(grid); y++ {
			for x := 0; x < len(grid[y]); x++ {
				ismonster := true
				for i := 0; i < len(seaMonster); i++ {
					smxy := seaMonster[i]

					if y+smxy.y >= len(grid) || y+smxy.y < 0 || x+smxy.x >= len(grid[y]) {
						ismonster = false
						break
					}

					if !grid[y+smxy.y][x+smxy.x] {
						ismonster = false
						break
					}
				}
				if ismonster {
					c++
				}
			}
		}
		if c == 0 {
			grid = rotate(grid)
			rotates++
		}

		if rotates == 4 {
			grid = flipVertical(grid)
		}
	}
	return c
}

func solve(m map[int]*block, d int) int {
	q := []int{}
	for k := range m {
		if len(q) == 0 {
			q = append(q, k)
			m[k].oriented = true
			break
		}
	}

	done := false
	for !done {
		k := q[0]
		q = q[1:]
		b := m[k]
		b.oriented = true

		exclude := []int{k}
		for _, n := range b.neighbors {
			if (n.neighbors[top] != nil && n.neighbors[bottom] != nil) ||
				(n.neighbors[left] != nil && n.neighbors[right] != nil) {
				exclude = append(exclude, n.id)
			}
		}
		ns := getNeighbors(b, exclude, m)
		for _, n := range ns {
			n.block.oriented = b.oriented
			b.neighbors[n.side] = n.block
			n.block.neighbors[n.side.opposite()] = b
			if len(n.block.neighbors) < 2 {
				q = append(q, n.block.id)
			}
		}
		done = len(q) == 0
	}

	prod := 1
	for _, b := range m {
		if len(b.neighbors) == 2 {
			prod *= b.id
		}
	}
	return prod
}

func joinGrids(m map[int]*block, d int) [][]bool {
	var cur *block
	for _, b := range m {
		if b.neighbors[top] == nil && b.neighbors[left] == nil {
			cur = b
			break
		}
	}

	rows, cols := 1, 1
	ptr := cur

	counted := false
	for !counted {
		if ptr.neighbors[right] != nil {
			cols++
			ptr = ptr.neighbors[right]
		}

		if ptr.neighbors[bottom] != nil {
			rows++
			ptr = ptr.neighbors[bottom]
		}

		counted = ptr.neighbors[right] == nil && ptr.neighbors[bottom] == nil
	}

	ret := make([][]bool, rows*(d-2))
	for i := 0; i < rows*(d-2); i++ {
		ret[i] = make([]bool, cols*(d-2))
	}
	done := false
	row, col := 0, 0
	for !done {
		for y := 1; y < d-1; y++ {
			for x := 1; x < d-1; x++ {
				ny := row*(d-2) + y - 1
				nx := col*(d-2) + x - 1
				ret[ny][nx] = cur.grid[y][x]
			}
		}

		if cur.neighbors[right] != nil {
			cur = cur.neighbors[right]
			col++
		} else {
			for cur.neighbors[left] != nil {
				cur = cur.neighbors[left]
			}
			cur = cur.neighbors[bottom]
			row++
			col = 0
		}

		done = cur == nil
	}
	return ret
}

func getNeighbors(b *block, exclude []int, m map[int]*block) []neighbor {
	ns := []neighbor{}

	exmap := make(map[int]int)
	for _, id := range exclude {
		exmap[id] = id
	}
	for k, v := range m {
		if _, ok := exmap[k]; ok {
			continue
		}

		ok, side := isMatch(b, v)
		if ok {
			n := neighbor{block: v, side: side}
			ns = append(ns, n)
		}
	}
	return ns
}

func isMatch(b, test *block) (bool, edge) {
	match := false
	e := top
	i := 0
	for i < 4 {
		match, e = checkEdges(b.grid, test.grid)

		if !match && !test.oriented {
			test.grid = flipVertical(test.grid)
			match, e = checkEdges(b.grid, test.grid)

			if !match && !test.oriented {
				test.grid = flipVertical(test.grid)
				test.grid = rotate(test.grid)
			}
		}
		i++
	}

	return match, e
}

func checkEdges(b, test [][]bool) (bool, edge) {
	edgemap := make(map[edge]int)

	for i := 0; i < len(b); i++ {
		if b[i][0] == test[i][len(b)-1] {
			edgemap[left]++
		}

		if b[0][i] == test[len(b)-1][i] {
			edgemap[top]++
		}

		if b[len(b)-1][i] == test[0][i] {
			edgemap[bottom]++
		}

		if b[i][len(b)-1] == test[i][0] {
			edgemap[right]++
		}
	}

	matched := false
	side := undefined
	for k, v := range edgemap {
		if v == len(b) {
			matched = true
			side = k
			break
		}
	}

	return matched, side
}

func printGrid(grid [][]bool) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			on := grid[y][x]
			if on {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func flipVertical(m [][]bool) [][]bool {
	m2 := make([][]bool, len(m))
	for y := 0; y < len(m); y++ {
		m2[y] = make([]bool, len(m[y]))
	}
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			m2[len(m)-y-1][x] = m[y][x]
		}
	}
	return m2
}

func rotate(m [][]bool) [][]bool {
	m2 := make([][]bool, len(m))
	for y := 0; y < len(m); y++ {
		m2[y] = make([]bool, len(m[y]))
		for x := 0; x < len(m[y]); x++ {
			m2[y][x] = m[len(m[y])-x-1][y]
		}
	}
	return m2
}

func readGrids(lines []string, d int) map[int]*block {
	m := make(map[int]*block)
	curId := 0
	x, y := 0, 0
	var curb *block
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "Tile ") {
			curId, _ = strconv.Atoi(line[5 : len(line)-1])
			curb = newBlock(curId, d)
			m[curId] = curb
			y = 0
			continue
		}

		x = 0
		curb.grid[y] = make([]bool, d)
		for _, ch := range line {
			curb.grid[y][x] = ch == '#'
			x++
		}
		y++
	}
	return m
}
