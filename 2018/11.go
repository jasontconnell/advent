package main

import (
	"fmt"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)
var input = 2694

func main() {
	startTime := time.Now()

	grid := make([][]int, 300)
	for i := 0; i < 300; i++ {
		grid[i] = make([]int, 300)
	}

	mapPowers(grid, input)
	x, y, _ := largestCube(grid, 3)
	fmt.Printf("%d,%d\n", x, y)

	x1, y1, size := largestVariableCube(grid)
	fmt.Printf("%d,%d,%d\n", x1, y1, size)

	fmt.Println("Time", time.Since(startTime))
}

func largestVariableCube(grid [][]int) (int, int, int) {
	maxval := -10000
	x, y := -1, -1
	size := 0

	for i := 2; i < 40; i++ {
		x1, y1, m := largestCube(grid, i)
		if m > maxval {
			maxval = m
			x = x1
			y = y1
			size = i
		}
	}
	return x, y, size
}

func largestCube(grid [][]int, size int) (int, int, int) {
	x, y := -1, -1
	max := -10000
	for i := 0; i < len(grid) - size; i++ {
		for j := 0; j < len(grid[i]) - size; j++ {
			sum := sumcube(grid, size, i, j)
			if sum > max {
				max = sum
				x = i+1
				y = j+1
			}
		}
	}

	return x, y, max
}

func sumcube(grid [][]int, size, x, y int) int {
	sum := 0
	for i := x; i < x + size; i++ {
		for j := y; j < y + size; j++ {
			sum += grid[i][j]
		}
	}
	return sum
}

func mapPowers(grid [][]int, serial int) {
	for i := 0; i < 300; i++ {
		for j := 0; j < 300; j++ {
			grid[i][j] = cellPower(serial, i+1, j+1)
		}
	}
}

func cellPower(serial, x, y int) int { // x and y are 0 based
	rackId := x + 10
	h := int((rackId * y + serial) * rackId / 100)
	if h == 0 {
		return 0
	}

	return h % 10 - 5
}