package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "13.txt"

const (
	Blank      int = 0
	Horizontal int = iota
	Vertical
	CurveFront //   /
	CurveBack  //   \
	Intersection
)

const (
	Left int = iota
	Right
	Up
	Down
)

const (
	TurnLeft int = iota
	GoStraight
	TurnRight
)

type xy struct {
	x, y int
}

type cart struct {
	xy
	id, dir      int
	decision     int
}

func (c cart) String() string {
	dir := ""
	switch c.dir {
	case Left:
		dir = "left"
	case Right:
		dir = "right"
	case Up:
		dir = "up"
	case Down:
		dir = "down"
	}

	return fmt.Sprintf("%d,%d - %s", c.x, c.y, dir)
}

func (c cart) char() string {
	dir := ""
	switch c.dir {
	case Left:
		dir = "<"
	case Right:
		dir = ">"
	case Up:
		dir = "^"
	case Down:
		dir = "v"
	}
	return dir
}

func clone(carts []cart) []cart {
	return append([]cart{}, carts...)
}

type path struct {
	xy
	dir int // one of - | + / \
	c   rune
}

func (p path) String() string {
	return fmt.Sprintf("%d,%d - %c", p.x, p.y, p.c)
}

func (p path) char() string {
	return string(p.c)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	grid := [][]path{}
	carts := []cart{}
	i := 0
	for scanner.Scan() {
		var txt = scanner.Text()
		g, c := getLine(i, txt)
		grid = append(grid, g)
		carts = append(carts, c...)
		i++
	}

	for i := 0; i < len(carts); i++ {
		carts[i].id = i
	}

	p2carts := clone(carts)
	crash := sim(grid, carts)
	last := lastStanding(grid, p2carts)
	fmt.Printf("first crash   :   %d,%d\n", crash.x, crash.y)
	fmt.Printf("last standing :   %d,%d\n", last.x, last.y)
	fmt.Println("Time", time.Since(startTime))
}

func print(grid [][]path, carts []cart) {
	m := make(map[xy]cart)
	crash := make(map[xy]rune)
	for _, c := range carts {
		key := xy{c.x, c.y}
		_, ok := m[key]
		m[key] = c
		if ok {
			crash[key] = 'X'
		}
	}

	for y := 0; y < len(grid); y++ {
		line := ""
		for x := 0; x < len(grid[y]); x++ {
			key := xy{x: x, y: y}
			p := grid[y][x]
			c, ok := m[key]
			x, crok := crash[key]

			str := p.char()
			if ok {
				str = c.char()
			}

			if crok {
				str = string(x)
			}
			line += str
		}

		fmt.Println(line)
	}
	fmt.Println("---------------------------")
}

func lastStanding(grid [][]path, carts []cart) cart {
	done := false

	for !done {
		for i := 0; i < len(carts); i++ {
			c := carts[i]
			c = move(c, grid)
			c = turn(c, grid)
			carts[i] = c
		}

		crashes := collisions(carts)
		for _, crash := range crashes {
			for i := len(carts)-1; i >= 0; i-- {
				cc := carts[i]
				if cc.x == crash.x && cc.y == crash.y {
					carts = append(carts[:i], carts[i+1:]...)
				}
			}
		}

		done = len(carts) == 1
	}

	return carts[0]
}

func sim(grid [][]path, carts []cart) xy {
	crash := xy{x: 0, y: 0}
	crashed := false

	for !crashed {
		for i := 0; i < len(carts); i++ {
			c := carts[i]
			c = move(c, grid)
			c = turn(c, grid)
			carts[i] = c
		}

		crashes := collisions(carts)
		if len(crashes) > 0 {
			crashed = true
			crash = crashes[0]
		}
	}
	return crash
}

func collisions(carts []cart) []xy {
	m := make(map[xy]bool)
	list := []xy{}
	for _, c := range carts {
		key := xy{x: c.x, y: c.y }
		_, ok := m[key]
		if ok {
			list = append(list, key)
		}
		m[key] = true
	}
	return list
}

func turn(c cart, grid [][]path) cart {
	p := grid[c.y][c.x]
	switch p.dir {
	case Intersection:
		if c.decision != GoStraight {
			switch c.dir {
			case Left:
				if c.decision == TurnLeft {
					c.dir = Down
				} else {
					c.dir = Up
				}
			case Right:
				if c.decision == TurnLeft {
					c.dir = Up
				} else {
					c.dir = Down
				}
			case Up:
				if c.decision == TurnLeft {
					c.dir = Left
				} else {
					c.dir = Right
				}
			case Down:
				if c.decision == TurnLeft {
					c.dir = Right
				} else {
					c.dir = Left
				}
			}
		}
		c.decision = (c.decision+1)%3
	case CurveFront:
		switch c.dir {
		case Left:
			c.dir = Down
		case Right:
			c.dir = Up
		case Down:
			c.dir = Left
		case Up:
			c.dir = Right
		}
	case CurveBack:
		switch c.dir {
		case Left:
			c.dir = Up
		case Right:
			c.dir = Down
		case Down:
			c.dir = Right
		case Up:
			c.dir = Left
		}
	}
	return c
}

func move(c cart, grid [][]path) cart {
	nx, ny := c.x, c.y
	switch c.dir {
	case Left:
		nx--
	case Right:
		nx++
	case Up:
		ny--
	case Down:
		ny++
	}
	c.x = nx
	c.y = ny

	return c
}

func getLine(i int, txt string) ([]path, []cart) {
	paths := []path{}
	carts := []cart{}
	for x, char := range txt {
		coords := xy{x: x, y: i}
		p := path{xy: coords, c: char}
		var c cart
		switch char {
		case '|':
			p.dir = Vertical
		case '-':
			p.dir = Horizontal
		case '/':
			p.dir = CurveFront
		case '\\':
			p.dir = CurveBack
		case '+':
			p.dir = Intersection
		case '>':
			p.dir = Horizontal
			p.c = '-'
			c = cart{xy: coords, dir: Right}
			carts = append(carts, c)
		case '^':
			p.dir = Vertical
			p.c = '|'
			c = cart{xy: coords, dir: Up}
			carts = append(carts, c)
		case 'v':
			p.dir = Vertical
			p.c = '|'
			c = cart{xy: coords, dir: Down}
			carts = append(carts, c)
		case '<':
			p.dir = Horizontal
			p.c = '-'
			c = cart{xy: coords, dir: Left}
			carts = append(carts, c)
		default:
			p.dir = Blank
		}
		paths = append(paths, p)
	}

	return paths, carts
}
