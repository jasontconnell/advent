package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

var input = "10.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	adapters := getAdapters(lines)
	p1map := getRatings(0, 3, adapters)
	j1, j3 := p1map[1], p1map[3]+1
	fmt.Println("Part 1:", j1*j3)

	fmt.Println("read", len(lines), "lines")
	fmt.Println("Time", time.Since(startTime))
}

func getRatings(start, maxdiff int, adapters []int) map[int]int {
	m := make(map[int]int)
	cur := start
	for i := 0; i < len(adapters); i++ {
		idx, jolt, diff := nextAdapter(cur, i, maxdiff, adapters)
		if jolt > 0 && diff <= maxdiff {
			m[diff]++
			cur = jolt
			i = idx
		}
	}
	return m
}

func getCombinations(start, maxdiff int, adapters []int) int64 {
	var x int64

	return x
}

func permutateCount(adapters []int) int64 {
	perms := [][]int{}

	if len(adapters) == 2 {
		perms = append(perms, []int{adapters[0], adapters[1]})
		perms = append(perms, []int{adapters[1], adapters[0]})
	} else {

	}

	return int64(len(perms))
}

func nextAdapter(curjolt, start, maxdiff int, adapters []int) (index, joltage, diff int) {
	loop := math.Min(float64(start+maxdiff), float64(len(adapters)))
	for i := start; i < int(loop); i++ {
		if adapters[i]-curjolt <= 3 {
			return i, adapters[i], adapters[i] - curjolt
		}
	}
	return start, -1, -1
}

func getAdapters(lines []string) []int {
	adapters := []int{}
	for _, line := range lines {
		x, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println(err)
		}
		adapters = append(adapters, x)
	}
	sort.Ints(adapters)
	return adapters
}
