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

var input = "05.txt"

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

	cp := make([]int, len(opcodes))
	cp2 := make([]int, len(opcodes))
	copy(cp, opcodes)
	copy(cp2, opcodes)

	c1 := intcode.NewComputer(cp)
	c1.Ins = []int{1}
	c1.Exec()
	fmt.Println("Part 1: ", c1.Outs[0])

	c2 := intcode.NewComputer(cp2)
	c2.Ins = []int{5}
	c2.Exec()
	fmt.Println("Part 2: ", c2.Outs[0])

	fmt.Println("Time", time.Since(startTime))
}
