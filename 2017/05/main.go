package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []int
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadInts(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 05 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	inst := make(input, len(in))
	copy(inst, in)
	return process(inst)
}

func part2(in input) output {
	inst := make(input, len(in))
	copy(inst, in)
	return processPart2(inst)
}

func process(inst []int) int {
	i := 0
	steps := 0
	max := len(inst)
	for i > -1 && i < max {
		steps++
		cur := inst[i]

		inst[i] = cur + 1
		i = i + cur
	}

	return steps
}

func processPart2(inst []int) int {
	i := 0
	steps := 0
	max := len(inst)
	for i > -1 && i < max {
		steps++
		cur := inst[i]
		jmp := 1
		if cur >= 3 {
			jmp = -1
		}

		inst[i] = cur + jmp
		i = i + cur
	}

	return steps
}
