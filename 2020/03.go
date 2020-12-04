package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

var input = "03.txt"

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

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(lines []string) int {
	return traverse(lines, 3, 1)
}

func part2(lines []string) int {
	slopes := []xy{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	answers := []int{}
	for _, s := range slopes {
		ans := traverse(lines, s.x, s.y)
		answers = append(answers, ans)
	}

	full := 1
	for _, ans := range answers {
		full = full * ans
	}
	return full
}

func traverse(lines []string, dx, dy int) int {
	x := 0
	trees := 0
	for i := 0; i < len(lines); i += dy {
		tree := isTree(lines[i], x)

		if tree {
			trees++
		}

		x += dx
	}
	return trees
}

func isTree(line string, x int) bool {
	return line[x%len(line)] == '#'
}
