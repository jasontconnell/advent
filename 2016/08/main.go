package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type rect struct {
	w, h int
}
type rotaterow struct {
	y   int
	val int
}
type rotatecol struct {
	x   int
	val int
}
type instruction struct {
	rect      *rect
	rotaterow *rotaterow
	rotatecol *rotatecol
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:")
	part2(in)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	insts := parseInput(in)
	display := initDisplay(50, 6)
	display = apply(display, insts)
	return countOn(display)
}

func part2(in input) {
	insts := parseInput(in)
	display := initDisplay(50, 6)
	display = apply(display, insts)
	print(display)
}

func countOn(display [][]bool) int {
	oncount := 0
	for y := 0; y < len(display); y++ {
		for x := 0; x < len(display[y]); x++ {
			if display[y][x] {
				oncount++
			}
		}
	}
	return oncount
}

func apply(display [][]bool, insts []instruction) [][]bool {
	applied := make([][]bool, len(display))
	for y, line := range display {
		applied[y] = make([]bool, len(line))
	}
	for _, inst := range insts {
		if inst.rect != nil {
			initRect(applied, inst.rect.w, inst.rect.h)
			continue
		}

		if inst.rotatecol != nil {
			rotateColumn(applied, inst.rotatecol.x, inst.rotatecol.val)
			continue
		}

		if inst.rotaterow != nil {
			rotateRow(applied, inst.rotaterow.y, inst.rotaterow.val)
			continue
		}
	}
	return applied
}

func initDisplay(w, h int) [][]bool {
	display := make([][]bool, h)
	for i := 0; i < h; i++ {
		display[i] = make([]bool, w)
	}
	return display
}

func print(display [][]bool) {
	for y := 0; y < len(display); y++ {
		for x := 0; x < len(display[y]); x++ {
			ch := " "
			if display[y][x] {
				ch = "#"
			}

			fmt.Print(ch)
		}
		fmt.Println("")
	}
}

func parseInput(in input) []instruction {
	rectreg := regexp.MustCompile("^rect (\\d+)x(\\d+)$")
	rotreg := regexp.MustCompile("^rotate (row|column) (x|y)=(\\d+) by (\\d+)$")

	instructions := []instruction{}
	for _, line := range in {
		inst := instruction{}
		if groups := rectreg.FindStringSubmatch(line); groups != nil {
			w, _ := strconv.Atoi(groups[1])
			h, _ := strconv.Atoi(groups[2])
			r := &rect{w, h}
			inst.rect = r
		}

		if groups := rotreg.FindStringSubmatch(line); groups != nil {
			ordinal, _ := strconv.Atoi(groups[3])
			count, _ := strconv.Atoi(groups[4])

			switch groups[1] {
			case "column":
				inst.rotatecol = &rotatecol{x: ordinal, val: count}
				break
			case "row":
				inst.rotaterow = &rotaterow{y: ordinal, val: count}
				break
			}
		}

		instructions = append(instructions, inst)
	}
	return instructions
}

func shift(slice []bool, count int) []bool {
	cp := make([]bool, len(slice))
	copy(cp, slice)

	for i := 0; i < len(slice); i++ {
		dstIndex := (i + count) % len(slice)
		cp[dstIndex] = slice[i]
	}
	return cp
}

func rotateColumn(display [][]bool, col, count int) {
	shifted := []bool{}

	for i := 0; i < len(display); i++ {
		shifted = append(shifted, display[i][col])
	}

	shifted = shift(shifted, count)

	for i := 0; i < len(display); i++ {
		display[i][col] = shifted[i]
	}
}

func rotateRow(display [][]bool, row, count int) {
	shifted := []bool{}

	for i := 0; i < len(display[row]); i++ {
		shifted = append(shifted, display[row][i])
	}

	shifted = shift(shifted, count)

	for i := 0; i < len(display[row]); i++ {
		display[row][i] = shifted[i]
	}
}

func initRect(display [][]bool, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			display[y][x] = true
		}
	}
}
