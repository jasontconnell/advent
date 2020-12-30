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
	for _, s := range grid {
		fmt.Println(s)
	}

	list := getIntersections(grid)

	p1 := 0
	for _, pt := range list {
		p1 += (pt.x * pt.y)
	}

	fmt.Println("Part 1:", p1)

	opcodes[0] = 2 // update program
	p2 := getDust(grid, opcodes)

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getDust(grid []string, prog []int) int {
	var start xy
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			ch := grid[y][x]
			if ch == 'v' || ch == '^' || ch == '<' || ch == '>' {
				start = xy{x, y}
				break
			}
		}
	}

	afunc := []int{1, 10}
	bfunc := []int{int('L'), 10}
	cfunc := []int{int('R'), 10}

	c := intcode.NewComputer(prog)

	c.OnOutput = func(out int) {
		fmt.Print(string(rune(out)))
	}

	fmt.Println(start, afunc, bfunc, cfunc)

	input = `B,A,B,C,A,C,A,C,B,C
R,6,R,10,L,10
L,10,R,12,R,12
R,10,L,10,L,12,R,6
`
	fmt.Println("inputs")
	for _, ch := range input {
		fmt.Println(int(ch))
		c.Ins = append(c.Ins, int(ch))
	}

	c.Exec()

	dust := 0

	for _, i := range c.Outs {
		dust += i
	}

	fmt.Println(dust)

	return dust
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
