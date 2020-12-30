package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/2019/intcode"
)

var input = "17.txt"

type xy struct {
	x, y int
}

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

	grid := getMap(prog)

	list := getIntersections(grid)

	p1 := 0
	for _, pt := range list {
		p1 += (pt.x * pt.y)
	}

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func getIntersections(grid []string) []xy {
	list := []xy{}
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			if grid[y][x] == '.' {
				continue
			}
			pt := xy{x, y}

			if grid[pt.y-1][pt.x] == '#' && grid[pt.y+1][pt.x] == '#' &&
				grid[pt.y][pt.x+1] == '#' && grid[pt.y][pt.x-1] == '#' {
				list = append(list, pt)
			}

		}
	}
	return list
}

func getMap(prog []int) []string {
	c := intcode.NewComputer(prog)

	s := []string{}
	s = append(s, "")

	c.OnOutput = func(out int) {
		if out != 10 {
			s[len(s)-1] += string(rune(out))
		} else {
			s = append(s, "")
		}
	}

	c.Exec()

	s = s[:len(s)-2]

	return s
}
