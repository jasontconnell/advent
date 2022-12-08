package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 08 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	grid := parseInput(in)
	return getVisible(grid)
}

func part2(in input) output {
	grid := parseInput(in)
	scores := calcScenicScores(grid)

	high := 0
	for _, r := range scores {
		for _, c := range r {
			if c > high {
				high = c
			}
		}
	}
	return high
}

func getVisible(grid [][]int) int {
	total := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if isVisible(grid, x, y) {
				total++
			}
		}
	}
	return total
}

func isVisible(grid [][]int, x, y int) bool {
	if x == 0 || y == 0 || y == len(grid)-1 || x == len(grid[y])-1 {
		return true
	}

	var visleft, visright, visup, visdown bool = true, true, true, true
	val := grid[y][x]
	for i := 0; i < y; i++ {
		if grid[i][x] >= val {
			visup = false
			break
		}
	}

	for i := y + 1; i < len(grid); i++ {
		if grid[i][x] >= val {
			visdown = false
			break
		}
	}

	for i := 0; i < x; i++ {
		if grid[y][i] >= val {
			visleft = false
			break
		}
	}

	for i := x + 1; i < len(grid[y]); i++ {
		if grid[y][i] >= val {
			visright = false
			break
		}
	}
	return visdown || visup || visleft || visright
}

func calcScenicScores(grid [][]int) [][]int {
	scores := [][]int{}
	for y := 1; y < len(grid)-1; y++ { // don't care about edges
		row := []int{}
		for x := 1; x < len(grid)-1; x++ {
			row = append(row, getScenicScore(grid, x, y))
		}
		scores = append(scores, row)
	}
	return scores
}

func getScenicScore(grid [][]int, x, y int) int {
	var left, right, up, down int = 0, 0, 0, 0
	val := grid[y][x]
	for i := y - 1; i >= 0; i-- {
		up++
		if grid[i][x] >= val {
			break
		}
	}

	for i := y + 1; i < len(grid); i++ {
		down++
		if grid[i][x] >= val {
			break
		}
	}

	for i := x - 1; i >= 0; i-- {
		left++
		if grid[y][i] >= val {
			break
		}
	}

	for i := x + 1; i < len(grid); i++ {
		right++
		if grid[y][i] >= val {
			break
		}
	}

	return left * right * up * down
}

func parseInput(in input) [][]int {
	grid := [][]int{}
	for _, line := range in {
		a := []int{}
		for _, c := range line {
			x, _ := strconv.Atoi(string(c))
			a = append(a, x)
		}
		grid = append(grid, a)
	}
	return grid
}
