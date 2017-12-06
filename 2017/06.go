package main

import (
	"fmt"
	"time"
	//"regexp"
	"strconv"
	//"math"
)

var input = []int{0, 5, 10, 0, 11, 14, 13, 4, 11, 8, 8, 7, 1, 4, 12, 11}

func main() {
	startTime := time.Now()

	c1, c2 := solve(input)

	fmt.Println("solve", c1, c2)

	fmt.Println("Time", time.Since(startTime))
}

func solve(nums []int) (int, int) {
	c1 := 0
	c2 := 0
	solved1 := false
	solved2 := false
	seen := make(map[string]bool)
	var find string

	for !(solved1 && solved2) {
		seenVal := ""
		for _, i := range nums {
			seenVal += strconv.Itoa(i) + "-"
		}

		if seenVal == find {
			solved2 = true
			break
		}

		if _, ok := seen[seenVal]; ok && !solved1 {
			solved1 = true
			find = seenVal
			c2 = 0
		}
		seen[seenVal] = true

		if !(solved1 && solved2) {
			processOne(nums)
		}

		if !solved1 {
			c1++
		}

		if !solved2 {
			c2++
		}
	}

	return c1, c2
}

func maxIndex(nums []int) int {
	index := 0
	value := -1
	for i := 0; i < len(nums); i++ {
		if nums[i] > value {
			value = nums[i]
			index = i
		}
	}

	return index
}

func processOne(nums []int) {
	index := maxIndex(nums)
	distribute(nums, index)
}

func distribute(nums []int, takeIndex int) {
	count := nums[takeIndex]
	nums[takeIndex] = 0
	for i := 0; i < count; i++ {
		index := (takeIndex + i + 1) % len(nums)
		nums[index]++
	}
}
