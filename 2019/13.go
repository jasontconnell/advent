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

type gametile struct {
	x, y int
	tile tile
}

type tile int

const (
	empty tile = iota
	wall
	block
	h_paddle
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
	p1 := part1(c)

	fmt.Println("Part 1: ", p1)

	fmt.Println("Time", time.Since(startTime))
}

func part1(c *intcode.Computer) int {
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

	blocks := 0
	for _, gt := range gametiles {
		if gt.tile == block {
			blocks++
		}
	}

	return blocks
}
