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

func run(ops, amps []int) int {
	var out int
	names := "ABCDE"
	comps := []*intcode.Computer{}
	for i, a := range amps {
		prog := make([]int, len(ops))
		copy(prog, ops)

		c := intcode.NewComputer(prog)
		c.Name = string(names[i])
		c.Ins =  []int{a}

		comps = append(comps, c)
	}

	for i := 0; i < len(comps); i++ {
		dp := i-1
		dn := i+1
		if dp < 0 {
			dp = len(comps)-1
		}
		if dn > len(comps)-1 {
			dn = 0
		}
		comps[i].Prev = comps[dp]
		comps[i].Next = comps[dn]
	}

	comps[0].Ins = append(comps[0].Ins, 0)

	for _, c := range comps {
		c.Exec()
	}

	out = comps[len(comps)-1].Outs[0]
	return out
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
	perms := Permutate([]int{5, 6, 7, 8, 9})
	output := 0

	for _, perm := range perms {
		o := run(ops, perm)
		if o > output {
			output = o
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
