package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input string

func main() {
	startTime := time.Now()

	line, err := common.ReadString(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(line))
	p2 := part2(input(line))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) int {
	vals := getInput(string(in))
	return process(vals, 40)
}

func part2(in input) int {
	vals := getInput(string(in))
	return process(vals, 50)
}

func process(vals []int, count int) int {
	nv := make([]int, len(vals))
	copy(nv, vals)
	for i := 0; i < count; i++ {
		nv = lookAndSay(nv)
	}
	return len(nv)
}

func getInput(line string) []int {
	vals := []int{}
	for _, c := range line {
		i, _ := strconv.Atoi(string(c))

		vals = append(vals, i)
	}
	return vals
}

func lookAndSay(ints []int) (output []int) {
	count, digit := 0, ints[0]
	for _, c := range ints {
		if digit == c {
			count++
		} else {
			output = append(output, count, digit)
			count = 1
		}
		digit = c
	}
	output = append(output, count, digit)
	return
}
