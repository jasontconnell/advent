package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "01.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	vals := []int{}
	for scanner.Scan() {
		var txt = scanner.Text()
		i, _ := strconv.Atoi(txt)
		vals = append(vals, i)
	}

	p1 := getIncs(vals)
	p2 := getIncSums(vals)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getIncs(vals []int) int {
	if len(vals) < 2 {
		return -1
	}
	incs := 0
	for i := 1; i < len(vals); i++ {
		if vals[i] > vals[i-1] {
			incs++
		}
	}
	return incs
}

func getIncSums(vals []int) int {
	if len(vals) < 2 {
		return -1
	}

	incs := 0

	last := vals[2] + vals[1] + vals[0]
	for i := 3; i < len(vals); i++ {
		cur := vals[i] + vals[i-1] + vals[i-2]

		if cur > last {
			incs++
		}
		last = cur
	}
	return incs
}
