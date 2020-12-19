package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

var input = "16.txt"

var basePattern []int = []int{0, 1, 0, -1}

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

	nums := getNums(lines[0])
	p1ary := solve(100, nums)

	fmt.Println("Part 1:", printnums(p1ary, 8))

	fmt.Println("Time", time.Since(startTime))
}

func printnums(ns []int, nlen int) string {
	s := ""
	for i, n := range ns {
		if i == nlen {
			break
		}
		s += strconv.Itoa(n)
	}
	return s
}

func solve(phases int, nums []int) []int {
	input := make([]int, len(nums))
	copy(input, nums)

	for p := 0; p < phases; p++ {
		input = runPhase(input)
	}
	return input
}

func runPhase(nums []int) []int {
	results := []int{}
	var steps int = len(nums)

	for i := 0; i < steps; i++ {
		mults := getPatternStep(basePattern, i, len(nums))
		stepSum := 0

		for j := 0; j < len(nums); j++ {
			x := nums[j] * mults[j]
			stepSum += x
		}
		results = append(results, int(math.Abs(float64(stepSum)))%10)
	}

	return results
}

func getPatternStep(pattern []int, step, length int) []int {
	pstep := 0
	stepPattern := []int{}

	// add one extra so we can trim it from the front
	for i := 0; len(stepPattern) < length+1; i++ {
		for j := 0; j < step+1 && len(stepPattern) < length+1; j++ {
			stepPattern = append(stepPattern, pattern[pstep])
		}
		pstep = (pstep + 1) % len(pattern)
	}
	return stepPattern[1:]
}

func getNums(line string) []int {
	nums := []int{}
	for _, ch := range line {
		n, _ := strconv.Atoi(string(ch))
		nums = append(nums, n)
	}
	return nums
}
