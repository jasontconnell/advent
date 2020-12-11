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
	p1map := getRatings(3, adapters)
	j1, j3 := p1map[1], p1map[3]+1
	fmt.Println("Part 1:", j1*j3)
	p2 := getCombinations(3, adapters)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getRatings(maxdiff int, adapters []int) map[int]int {
	m := make(map[int]int)
	cur := 0
	for i := 0; i < len(adapters); i++ {
		jolt, diff := nextAdapter(cur, i, maxdiff, adapters)
		if jolt > 0 && diff <= maxdiff {
			m[diff]++
			cur = jolt
		}
	}
	return m
}

func getCombinations(maxdiff int, adapters []int) int64 {
	total := append([]int{0}, adapters...)
	total = append(total, adapters[len(adapters)-1])
	var x int64
	allOnes := [][]int{}
	allOnes = append(allOnes, []int{})
	cur := 0
	for i := 0; i < len(total); i++ {
		jolt, diff := nextAdapter(cur, i, maxdiff, total)
		cur = jolt
		if diff == 1 {
			allOnes[len(allOnes)-1] = append(allOnes[len(allOnes)-1], jolt)
		} else {
			allOnes = append(allOnes, []int{jolt})
		}
	}
	x = 1
	for _, a := range allOnes {
		x = x * combos(a)
	}
	return x
}

func combos(a []int) int64 {
	var c int64 = 1
	switch len(a) {
	case 3:
		c = 2
	case 4:
		c = 4
	case 5:
		c = 7
	}
	return c
}

func fact(i int) int {
	if i < 3 {
		return i
	}
	return i * fact(i-1)
}

func nextAdapter(curjolt, curidx, maxdiff int, adapters []int) (joltage, diff int) {
	loop := math.Min(float64(curidx+maxdiff), float64(len(adapters)))
	for i := curidx; i < int(loop); i++ {
		if adapters[i]-curjolt <= maxdiff {
			return adapters[i], adapters[i] - curjolt
		}
	}
	return -1, -1
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
