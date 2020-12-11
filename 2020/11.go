package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "11.txt"

type block int

const (
	seat block = iota
	floor
)

type xy struct {
	x, y int
}

type gridblock struct {
	xy
	block    block
	occupied bool
}

func (gb gridblock) String() string {
	b := "seat"
	if gb.block == floor {
		b = "floor"
	}
	return fmt.Sprintf("(%d,%d) %s %v", gb.x, gb.y, b, gb.occupied)
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

	grid := getGrid(lines)
	p1 := sim(grid, 100)

	fmt.Println("Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
}

func sim(grid [][]gridblock, threshold int) int {
	changed := 0
	prev := 0
	g := copyGrid(grid)

	//printGrid(g)

	streak := 0
	done := false
	for !done {
		g = simOne(g)
		//printGrid(g)
		changed = totalOccupied(g)

		if changed == prev {
			streak++
		}

		prev = changed

		done = streak > threshold
	}

	return changed
}

func simOne(grid [][]gridblock) [][]gridblock {
	gcopy := copyGrid(grid)

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			s := grid[y][x]
			if s.block == floor {
				continue
			}

			nocc := filterCount(grid, x, y, func(adj gridblock) bool {
				return adj.block == seat && adj.occupied
			})

			if nocc == 0 && !s.occupied {
				gcopy[y][x].occupied = true
			}

			// nvac := filterCount(grid, x, y, func(adj gridblock) bool {
			// 	return adj.block == seat && !adj.occupied
			// })

			if nocc > 3 && s.occupied {
				gcopy[y][x].occupied = false
			}
		}
	}

	return gcopy
}

func totalOccupied(grid [][]gridblock) int {
	count := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			b := grid[y][x]

			if b.occupied {
				count++
			}
		}
	}
	return count
}

func copyGrid(grid [][]gridblock) [][]gridblock {
	c := make([][]gridblock, len(grid))
	for i := 0; i < len(grid); i++ {
		c[i] = make([]gridblock, len(grid[i]))
		copy(c[i], grid[i])
	}
	return c
}

func filterCount(grid [][]gridblock, x, y int, test func(gb gridblock) bool) int {
	w, h := len(grid[0])-1, len(grid)-1
	pts := []xy{
		{x, y - 1},
		{x, y + 1},
		{x + 1, y},
		{x - 1, y},
		{x + 1, y + 1},
		{x + 1, y - 1},
		{x - 1, y + 1},
		{x - 1, y - 1},
	}

	count := 0
	for _, pt := range pts {
		if pt.x < 0 || pt.x > w ||
			pt.y < 0 || pt.y > h {
			continue
		}

		b := grid[pt.y][pt.x]
		if test(b) {
			count++
		}

	}
	return count
}

func getGrid(lines []string) [][]gridblock {
	grid := [][]gridblock{}
	for y, line := range lines {
		gline := []gridblock{}
		for x, c := range line {
			s := gridblock{xy: xy{x: x, y: y}}
			switch c {
			case 'L':
				s.block = seat
			case '.':
				s.block = floor
			case '#':
				s.block = seat
				s.occupied = true
			}
			gline = append(gline, s)
		}

		grid = append(grid, gline)
	}

	return grid
}

func printGrid(grid [][]gridblock) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			b := grid[y][x]

			c := '.'

			if b.block == seat {
				if b.occupied {
					c = '#'
				} else {
					c = 'L'
				}
			}

			fmt.Print(string(c))
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
