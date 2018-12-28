package main

import (
	"bufio"
	"fmt"
	"os"
	//"strings"
	"time"
)

var input = "20.txt"

const (
	East       rune = 'E'
	West       rune = 'W'
	South      rune = 'S'
	North      rune = 'N'
	OpenParen  rune = '('
	CloseParen rune = ')'
	Pipe       rune = '|'
	Start      rune = '^'
	End        rune = '$'
)

type block struct {
	contents int
	dist     int
}

const (
	Unknown int = iota
	Open
	HDoor
	VDoor
	Wall
	StartPos
)

type node struct {
	val      string
	children []*node
	open     bool
}

type xy struct {
	x, y int
}

type minmax struct {
	minx, maxx, miny, maxy int
}

type stack struct {
	ary []*node
}

func (s *stack) push(n *node) {
	s.ary = append([]*node{n}, s.ary...)
}

func (s *stack) pop() *node {
	n := s.ary[0]
	s.ary = s.ary[1:]
	return n
}

func (s *stack) popOpen() {
	for !s.ary[0].open {
		s.pop()
	}
}

func (s *stack) peek(n int) *node {
	return s.ary[n]
}

func (s *stack) peekOpen() *node {
	i := 0
	n := s.ary[i]
	for !n.open {
		i++
		n = s.ary[i]
	}
	return n
}

func (n *node) String() string {
	s := "<not nil> " + n.val
	return s
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	var root *node
	for scanner.Scan() {
		var txt = scanner.Text()
		root = parse(txt)
	}

	dim := 400
	sy, sx := dim/2, dim/2

	start := xy{sx, sy}
	grid := makeGrid(dim, dim)
	grid[sy][sx].contents = StartPos
	grid[sy][sx].dist = 0
	grid = traverse(grid, start, root, 0)
	//draw(grid)

	p1 := getLongest(grid)
	p2 := getCountMin(grid, 1000)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func getLongest(grid [][]block) int {
	max := 0
	for _, line := range grid {
		for _, b := range line {
			if b.contents == Open && b.dist > max {
				max = b.dist
			}
		}
	}
	return max
}

func getCountMin(grid [][]block, min int) int {
	cc := 0
	for _, line := range grid {
		for _, b := range line {
			if b.contents == Open && b.dist >= min {
				cc++
			}
		}
	}
	return cc
}

func draw(grid [][]block) {
	for y := 0; y < len(grid); y++ {
		line := ""
		vy := len(grid) - 1 - y
		for x := 0; x < len(grid[y]); x++ {
			var c rune
			switch grid[vy][x].contents {
			case Unknown:
				c = '?'
			case Wall:
				c = '#'
			case Open:
				c = '.'
			case VDoor:
				c = '|'
			case HDoor:
				c = '-'
			case StartPos:
				c = 'X'
			}
			line += string(c)
		}
		fmt.Println(line)
	}
	fmt.Println("-------------------------------")
}

func makeGrid(w, h int) [][]block {
	grid := make([][]block, h)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			grid[i] = append(grid[i], block{dist: 10000, contents: Wall})
		}
	}
	return grid
}

func traverse(grid [][]block, start xy, n *node, d int) [][]block {
	for _, c := range n.children {
		p, dt := path(grid, c.val, start, d)
		grid = traverse(grid, p, c, dt)
	}
	return grid
}

func path(grid [][]block, dirs string, start xy, d int) (xy, int) {
	p := start
	dist := d
	for _, c := range dirs {
		switch rune(c) {
		case North:
			grid[p.y+1][p.x].contents = HDoor
			grid[p.y+2][p.x].contents = Open
			if dist+1 < grid[p.y+2][p.x].dist {
				grid[p.y+2][p.x].dist = dist + 1
			}
			p.y += 2
		case South:
			grid[p.y-1][p.x].contents = HDoor
			grid[p.y-2][p.x].contents = Open
			if dist+1 < grid[p.y-2][p.x].dist {
				grid[p.y-2][p.x].dist = dist + 1
			}
			p.y -= 2
		case East:
			grid[p.y][p.x+1].contents = VDoor
			grid[p.y][p.x+2].contents = Open
			if dist+1 < grid[p.y][p.x+2].dist {
				grid[p.y][p.x+2].dist = dist + 1
			}
			p.x += 2
		case West:
			grid[p.y][p.x-1].contents = VDoor
			grid[p.y][p.x-2].contents = Open
			if dist+1 < grid[p.y][p.x-2].dist {
				grid[p.y][p.x-2].dist = dist + 1
			}
			p.x -= 2
		}
		dist++
	}
	return p, dist
}

func parse(line string) *node {
	root := &node{}
	curnode := root
	cur := 0
	stk := new(stack)
	stk.push(curnode)
	for cur < len(line) {
		c := rune(line[cur])

		switch c {
		case Start, End:
			break
		case North, South, East, West:
			//curnode.val = getOrds(line, cur)
			n := &node{val: getOrds(line, cur)}
			curnode.children = append(curnode.children, n)
			curnode = n
			stk.push(curnode)
			cur += len(curnode.val) - 1
		case Pipe:
			curnode = stk.peekOpen()
		case OpenParen:
			n := &node{open: true}
			curnode.children = append(curnode.children, n)
			curnode = n
			stk.push(curnode)
		case CloseParen:
			stk.popOpen()
			curnode = stk.pop()
		}

		cur++
	}

	return root
}

func getOrds(line string, i int) string {
	s := ""
	done := false

	for !done {
		switch rune(line[i]) {
		case North, South, East, West:
			s += string(line[i])
		default:
			done = true
		}
		i++
	}
	return s
}
