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

var input = "07.txt"

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

func part1(ops []int) int {
	perms := Permutate([]int{0, 1, 2, 3, 4})
	output := 0
	for _, perm := range perms {
		o := run(ops, perm)
		if o > output {
			output = o
		}
	}
	return output
}

func part2(ops []int) int {
	perms := Permutate([]int{5,6,7,8,9})
	output := 0
	done := false

	for _, perm := range perms {
		o := run(ops, perm)
		if o > output {
			output = o
		}
	}
	return output
}

func run(ops []int, amps []int, input int) int {
	input2 := input
	output := 0
	for i, amp := range amps {
		prog := make([]int, len(ops))
		copy(prog, ops)
		
		_, outs := intcode.Exec(prog, []int{amp, input2})

		input2 = outs[0]
		if i == len(amps)-1 {
			output = input2
		}
	}
	return output
}

func Permutate(ints []int) [][]int {
	var ret [][]int

	if len(ints) == 2 {
		ret = append(ret, []int{ints[0], ints[1]})
		ret = append(ret, []int{ints[1], ints[0]})
	} else {

		for i := 0; i < len(ints); i++ {
			cp := make([]int, len(ints))
			copy(cp, ints)

			t := cp[i]
			sh := append(cp[:i], cp[i+1:]...)
			perm := Permutate(sh)

			for _, p := range perm {
				p = append([]int{t}, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}
