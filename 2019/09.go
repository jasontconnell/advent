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

	// part1([]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99})
	// return
	p1 := part1(opcodes)

	fmt.Println("Part 1: ", p1)

	fmt.Println("Time", time.Since(startTime))
}

func part1(opcodes []int) int {
	prog := make([]int, len(opcodes))

	copy(prog, opcodes)
	c := intcode.NewComputer(prog)
	c.Ins = append(c.Ins, 1)
	c.Exec()

	return c.Outs[0]
}
