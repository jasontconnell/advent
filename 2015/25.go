package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	f := getCode(3019, 3010)
	fmt.Println(f)

	fmt.Println("Time", time.Since(startTime))
}

func getCode(x, y int) int64 {
	ordinal := getOrdinalForCoords(x, y)
	return getValueAtOrdinal(ordinal)
}

func getValueAtOrdinal(x int) int64 {
	if x == 1 {
		return int64(20151125)
	}

	prev := getValueAtOrdinal(x-1)
	mult := prev * int64(252533)
	mod := mult % int64(33554393)
	return mod
}

func getOrdinalForCoords(x, y int) int {
	c := numInCol(y + x - 2) + x
	return c
}

func numInCol(x int) int {
	ans := 0
	for i := 1; i <= x; i++ {
		ans += i
	}
	return ans
}