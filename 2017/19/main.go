package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}
type block struct {
	passable   bool
	horizontal bool
	vertical   bool
	cross      bool
	letter     bool
	char       rune
}

type move struct {
	pt    xy
	dir   xy
	steps int
}

func (m move) key() string {
	d := "h"
	if m.dir.y != 0 {
		d = "v"
	}

	return fmt.Sprintf("%v|%v|%v", m.pt.x, m.pt.y, d)
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) string {
	grid := parseInput(in)
	str, _ := traverse(grid)
	return str
}

func part2(in input) output {
	grid := parseInput(in)
	_, val := traverse(grid)
	return val

}

func fixBlocks(blocks [][]block) [][]block {
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			if blocks[y][x].letter {
				blocks[y][x].horizontal = (x > 0 && blocks[y][x-1].horizontal) || (x < len(blocks[y])-1 && blocks[y][x+1].horizontal) ||
					(x > 0 && x < len(blocks[y])-1 && blocks[y][x-1].vertical && blocks[y][x+1].vertical)
				blocks[y][x].vertical = (y > 0 && blocks[y-1][x].vertical) || (y < len(blocks)-1 && blocks[y+1][x].vertical) ||
					(y > 0 && y < len(blocks)-1 && blocks[y-1][x].horizontal && blocks[y+1][x].horizontal)
			}
		}
	}

	return blocks
}

func getStart(blocks [][]block) (int, int) {
	for i := 0; i < len(blocks[0]); i++ {
		if blocks[0][i].char == '|' {
			return i, 0
		}
	}
	return -1, -1
}

func traverse(blocks [][]block) (string, int) {
	startx, starty := getStart(blocks)
	visited := make(map[string]bool)

	first := move{pt: xy{startx, starty}, dir: xy{0, 1}, steps: 0}
	visited[first.key()] = true
	queue := getMoves(visited, blocks, first)
	path := ""

	var last move

	for len(queue) > 0 {
		mv := queue[0]
		queue = queue[1:]

		mvs := getMoves(visited, blocks, mv)

		queue = append(queue, mvs...)

		if blocks[mv.pt.y][mv.pt.x].letter {
			path = path + string(blocks[mv.pt.y][mv.pt.x].char)
		}

		if len(mvs) == 0 {
			last = mv
			last.steps = mv.steps + 1
		}
	}

	return path, last.steps
}

func getHorizontalMoves(blocks [][]block, from move) []move {
	mvs := []move{}
	if from.pt.x > 0 && blocks[from.pt.y][from.pt.x-1].passable {
		mvs = append(mvs, move{pt: xy{from.pt.x - 1, from.pt.y}, dir: xy{-1, 0}, steps: from.steps + 1})
	}

	if from.pt.x < len(blocks[from.pt.y])-1 && blocks[from.pt.y][from.pt.x+1].passable {
		mvs = append(mvs, move{pt: xy{from.pt.x + 1, from.pt.y}, dir: xy{1, 0}, steps: from.steps + 1})
	}

	return mvs
}

func getVerticalMoves(blocks [][]block, from move) []move {
	mvs := []move{}
	if from.pt.y > 0 && blocks[from.pt.y-1][from.pt.x].passable {
		mvs = append(mvs, move{pt: xy{from.pt.x, from.pt.y - 1}, dir: xy{0, 1}, steps: from.steps + 1})
	}

	if from.pt.y < len(blocks)-1 && blocks[from.pt.y+1][from.pt.x].passable {
		mvs = append(mvs, move{pt: xy{from.pt.x, from.pt.y + 1}, dir: xy{0, 1}, steps: from.steps + 1})
	}

	return mvs
}

func getValidMoves(blocks [][]block, from move) []move {
	b := blocks[from.pt.y][from.pt.x]

	mvs := []move{}

	switch b.char {
	case '+':
		mvs = getHorizontalMoves(blocks, from)
		mvs = append(mvs, getVerticalMoves(blocks, from)...)
	default:
		if from.dir.x != 0 {
			mvs = getHorizontalMoves(blocks, from)
		} else {
			mvs = getVerticalMoves(blocks, from)
		}
	}

	return mvs
}

func getMoves(visited map[string]bool, blocks [][]block, from move) []move {
	valid := getValidMoves(blocks, from)
	mvs := []move{}

	for _, mv := range valid {
		if _, ok := visited[mv.key()]; !ok {
			visited[mv.key()] = true
			mvs = append(mvs, mv)
		}
	}

	return mvs
}

func parseInput(in input) [][]block {
	grid := [][]block{}
	for _, line := range in {
		bline := []block{}
		for _, c := range line {
			b := block{char: c, passable: c != ' ', vertical: c == '|' || c == '+', horizontal: c == '-' || c == '+', cross: c == '+'}
			b.letter = b.passable && !b.vertical && !b.horizontal && !b.cross
			bline = append(bline, b)

		}
		grid = append(grid, bline)
	}
	return grid
}
