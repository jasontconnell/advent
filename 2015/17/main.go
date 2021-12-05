package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []int
type output int

func main() {
	startTime := time.Now()

	in, err := common.ReadInts(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return output(count(in, 150, len(in), 0))
}

func part2(in input) output {
	result := 0
	i := 0
	for result == 0 {
		result = count(in, 150, i, 0)
		i++
	}
	return output(result)
}

func count(list []int, total, n, i int) int {
	if i < 0 {
		i = 0
	}

	if n < 0 {
		return 0
	} else if total == 0 {
		return 1
	} else if i == len(list) || total < 0 {
		return 0
	} else {
		return count(list, total, n, i+1) + count(list, total-list[i], n-1, i+1)
	}
}
