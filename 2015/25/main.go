package main

import (
	"fmt"
	"time"
)

var inputFilename = "input.txt"

type output int64

func main() {
	startTime := time.Now()

	p1 := part1()

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func part1() int64 {
	return getCode(3019, 3010)
}

func getCode(x, y int) int64 {
	ordinal := getOrdinalForCoords(x, y)
	return getValueAtOrdinal(ordinal)
}

func getValueAtOrdinal(x int) int64 {
	if x == 1 {
		return int64(20151125)
	}

	prev := getValueAtOrdinal(x - 1)
	mult := prev * int64(252533)
	mod := mult % int64(33554393)
	return mod
}

func getOrdinalForCoords(x, y int) int {
	c := numInCol(y+x-2) + x
	return c
}

func numInCol(x int) int {
	ans := 0
	for i := 1; i <= x; i++ {
		ans += i
	}
	return ans
}
