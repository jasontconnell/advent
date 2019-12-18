package main

import (
	"bufio"
	"fmt"
	"github.com/jasontconnell/advent/2019/intcode"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "13.txt"

type xy struct {
}

type gametile struct {
	x, y  int
	tile  tile
	score int
}

type tile int

const (
	empty tile = iota
	wall
	block
	paddle
	ball
)

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	opcodes := []int{}
	if scanner.Scan() {
		var txt = scanner.Text()
		sopcodes := strings.Split(txt, ",")
		for _, s := range sopcodes {
			i, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				continue
			}

			opcodes = append(opcodes, i)
		}
	}

	prog := make([]int, len(opcodes))
	copy(prog, opcodes)

	c := intcode.NewComputer(prog)
	board := getBoard(c)
	p1 := countBlocks(board)
	fmt.Println("Part 1: ", p1)

	c2 := intcode.NewComputer(opcodes)
	p2 := part2(c2)

	fmt.Println("Part 2: ", p2)
	fmt.Println("Time", time.Since(startTime))
}

func countBlocks(gametiles []gametile) int {
	blocks := 0
	for _, gt := range gametiles {
		if gt.tile == block {
			blocks++
		}
	}
	return blocks
}

func getBoard(c *intcode.Computer) []gametile {
	instridx := 0

	gametiles := []gametile{}
	tileIndex := 0

	c.OnOutput = func(val int) {
		var gt gametile

		if len(gametiles) == 0 || tileIndex >= len(gametiles) {
			gametiles = append(gametiles, gametile{})
		}

		gt = gametiles[tileIndex]

		switch instridx {
		case 0:
			gt.x = val

		case 1:
			gt.y = val
		case 2:
			gt.tile = tile(val)
		}

		gametiles[tileIndex] = gt

		instridx = (instridx + 1) % 3
		if instridx == 0 {
			tileIndex++
		}
	}
	c.Exec()

	return gametiles
}

func part2(c *intcode.Computer) int {
	instridx := 0

	// gametiles := []gametile{}
	score := 0

	var current gametile

	// insert two quarters
	c.Ops[0] = 2

	ballx, paddlex := -1, -1

	c.OnOutput = func(val int) {
		switch instridx {
		case 0:
			current.x = val
		case 1:
			current.y = val
		case 2:
			if current.x == -1 && current.y == 0 {
				current.score = val
				score = val
			} else {
				current.tile = tile(val)
				if current.tile == ball {
					ballx = current.x
				} else if current.tile == paddle {
					paddlex = current.x
				}
			}
		}

		instridx = (instridx + 1) % 3
	}

	c.RequestInput = func() int {
		if ballx == -1 && paddlex == -1 {
			panic("input requested before init final")
		}

		if ballx > paddlex {
			return 1
		} else if ballx < paddlex {
			return -1
		} else {
			return 0
		}
	}

	c.Exec()

	return score
}
