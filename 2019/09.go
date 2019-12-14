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

var input = "09.txt"

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

	p1 := part1(opcodes)
	fmt.Println("Part 1: ", p1)

	p2 := part2(opcodes)
	fmt.Println("Part 2: ", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(opcodes []int) int {
	prog := make([]int, len(opcodes))

	copy(prog, opcodes)
	c := intcode.NewComputer(prog)
	c.AddInput(1)
	c.Exec()

	return c.Outs[0]
}

func part2(opcodes []int) int {
	prog := make([]int, len(opcodes))

	copy(prog, opcodes)
	c := intcode.NewComputer(prog)
	c.AddInput(2)
	c.Exec()

	return c.Outs[0]
}
