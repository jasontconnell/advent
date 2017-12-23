package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "19.txt"

type Block struct {
	X, Y       int
	Passable   bool
	Horizontal bool
	Vertical   bool
	Cross      bool
	Letter     bool
	Rune       rune
}

func (b Block) String() string {
	return fmt.Sprintf("(%v,%v) Passable: %v  H: %v V: %v C: %v L: %v R: %v", b.X, b.Y, b.Passable, b.Horizontal, b.Vertical, b.Cross, b.Letter, string(b.Rune))
}

type Move struct {
	X, Y       int
	DirX, DirY int
	Steps      int
}

func (m Move) String() string {
	return fmt.Sprintf("(%v,%v) - dir (%v, %v)", m.X, m.Y, m.DirX, m.DirY)
}

func (m Move) Key() string {
	d := "h"
	if m.DirY != 0 {
		d = "v"
	}

	return fmt.Sprintf("%v|%v|%v", m.X, m.Y, d)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	blocks := [][]Block{}
	y := 0

	for scanner.Scan() {
		var txt = scanner.Text()
		blocks = append(blocks, getBlocks(y, txt))
		y++
	}

	blocks = fixBlocks(blocks)
	path, steps := traverse(blocks)

	fmt.Println("Path :", path)
	fmt.Println("Steps:", steps)

	fmt.Println("Time", time.Since(startTime))
}

func getBlocks(y int, line string) []Block {
	lineBlocks := []Block{}
	for x, c := range line {
		b := Block{X: x, Y: y, Rune: c, Passable: c != ' ', Vertical: c == '|' || c == '+', Horizontal: c == '-' || c == '+', Cross: c == '+'}
		b.Letter = b.Passable && !b.Vertical && !b.Horizontal && !b.Cross
		lineBlocks = append(lineBlocks, b)
	}
	return lineBlocks
}

func fixBlocks(blocks [][]Block) [][]Block {
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			if blocks[y][x].Letter {
				blocks[y][x].Horizontal = (x > 0 && blocks[y][x-1].Horizontal) || (x < len(blocks[y])-1 && blocks[y][x+1].Horizontal) ||
					(x > 0 && x < len(blocks[y])-1 && blocks[y][x-1].Vertical && blocks[y][x+1].Vertical)
				blocks[y][x].Vertical = (y > 0 && blocks[y-1][x].Vertical) || (y < len(blocks)-1 && blocks[y+1][x].Vertical) ||
					(y > 0 && y < len(blocks)-1 && blocks[y-1][x].Horizontal && blocks[y+1][x].Horizontal)
			}
		}
	}

	return blocks
}

func getStart(blocks [][]Block) (int, int) {
	for i := 0; i < len(blocks[0]); i++ {
		if blocks[0][i].Rune == '|' {
			return i, 0
		}
	}
	return -1, -1
}

func traverse(blocks [][]Block) (string, int) {
	startx, starty := getStart(blocks)
	visited := make(map[string]bool)

	first := Move{X: startx, Y: starty, DirX: 0, DirY: 1, Steps: 0}
	visited[first.Key()] = true
	queue := getMoves(visited, blocks, first)
	path := ""

	var last Move

	for len(queue) > 0 {
		mv := queue[0]
		queue = queue[1:]

		mvs := getMoves(visited, blocks, mv)

		queue = append(queue, mvs...)

		if blocks[mv.Y][mv.X].Letter {
			path = path + string(blocks[mv.Y][mv.X].Rune)
		}

		if len(mvs) == 0 {
			last = mv
			last.Steps = mv.Steps + 1
		}
	}

	return path, last.Steps
}

func getHorizontalMoves(blocks [][]Block, from Move) []Move {
	mvs := []Move{}
	if from.X > 0 && blocks[from.Y][from.X-1].Passable {
		mvs = append(mvs, Move{X: from.X - 1, Y: from.Y, DirX: -1, DirY: 0, Steps: from.Steps + 1})
	}

	if from.X < len(blocks[from.Y])-1 && blocks[from.Y][from.X+1].Passable {
		mvs = append(mvs, Move{X: from.X + 1, Y: from.Y, DirX: 1, DirY: 0, Steps: from.Steps + 1})
	}

	return mvs
}

func getVerticalMoves(blocks [][]Block, from Move) []Move {
	mvs := []Move{}
	if from.Y > 0 && blocks[from.Y-1][from.X].Passable {
		mvs = append(mvs, Move{X: from.X, Y: from.Y - 1, DirX: 0, DirY: -1, Steps: from.Steps + 1})
	}

	if from.Y < len(blocks)-1 && blocks[from.Y+1][from.X].Passable {
		mvs = append(mvs, Move{X: from.X, Y: from.Y + 1, DirX: 0, DirY: 1, Steps: from.Steps + 1})
	}

	return mvs
}

func getValidMoves(blocks [][]Block, from Move) []Move {
	b := blocks[from.Y][from.X]

	mvs := []Move{}

	switch b.Rune {
	case '+':
		mvs = getHorizontalMoves(blocks, from)
		mvs = append(mvs, getVerticalMoves(blocks, from)...)
	default:
		if from.DirX != 0 {
			mvs = getHorizontalMoves(blocks, from)
		} else {
			mvs = getVerticalMoves(blocks, from)
		}
	}

	return mvs
}

func getMoves(visited map[string]bool, blocks [][]Block, from Move) []Move {
	valid := getValidMoves(blocks, from)
	mvs := []Move{}

	for _, mv := range valid {
		if _, ok := visited[mv.Key()]; !ok {
			visited[mv.Key()] = true
			mvs = append(mvs, mv)
		}
	}

	return mvs
}
