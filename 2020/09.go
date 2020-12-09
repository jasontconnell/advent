package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "09.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	nums := []int{}
	for scanner.Scan() {
		var txt = scanner.Text()
		i, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nums = append(nums, i)
	}

	p1 := part1(nums, 25)
	fmt.Println("part 1:", p1)

	fmt.Println("read", len(nums), "numbers")
	fmt.Println("Time", time.Since(startTime))
}

func part1(nums []int, size int) int {
	p1 := 0
	for i := size; i < len(nums); i++ {
		preamble := getPreamble(nums, i-size, size)
		if !hasMatch(nums[i], preamble) {
			p1 = nums[i]
			break
		}
	}
	return p1
}

func hasMatch(num int, preamble map[int]int) bool {
	for _, v := range preamble {
		x := num - v
		_, ok := preamble[x]
		if ok {
			return true
		}
	}
	return false
}

func getPreamble(nums []int, start, size int) map[int]int {
	preamble := make(map[int]int)
	for i := start; i < start+size; i++ {
		preamble[nums[i]] = nums[i]
	}
	return preamble
}
