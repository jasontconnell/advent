package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "15.txt"

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

	p1 := fill(getStart(lines), 2020)
	fmt.Println("Part 1:", p1[len(p1)-1])

	p2 := fill(getStart(lines), 30000000)
	fmt.Println("Part 2:", p2[len(p2)-1])

	fmt.Println("Time", time.Since(startTime))
}

func fill(ss []int, c int) []int {
	lasts := make(map[int]int)
	for i := 0; i < len(ss)-1; i++ {
		lasts[ss[i]] = i
	}

	idx := len(ss) - 1
	for len(ss) < c {
		prev := ss[idx]
		pos, ok := lasts[prev]
		lasts[prev] = idx

		num := 0
		if ok {
			num = idx - pos
		}

		ss = append(ss, num)

		idx++
	}

	return ss
}

func getStart(lines []string) []int {
	if len(lines) > 1 {
		panic("invalid input")
	}

	nums := []int{}
	ss := strings.Split(lines[0], ",")
	for _, s := range ss {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}
