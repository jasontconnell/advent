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
	p1ary := solve(100, false, nums)

	fmt.Println("Part 1:", printnums(p1ary, 8))

	offset, _ := strconv.Atoi(printnums(nums, 7))

	nums = repeat(nums, 10000)
	p2ary := solve(100, true, nums[offset:])

	fmt.Println("Part 2:", printnums(p2ary, 8))

	fmt.Println("Time", time.Since(startTime))
}

func repeat(nums []int, c int) []int {
	dnums := []int{}

	for i := 0; i < c; i++ {
		for i := 0; i < len(nums); i++ {
			dnums = append(dnums, nums[i])
		}
	}

	return dnums
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

func solve(phases int, rpt bool, nums []int) []int {
	input := make([]int, len(nums))
	copy(input, nums)

	for p := 0; p < phases; p++ {
		if !rpt {
			input = runPhase(input)
		} else {
			input = runPhase2(input)
		}
	}
	return input
}

func runPhase2(nums []int) []int {
	out := make([]int, len(nums))
	copy(out, nums)

	sum := 0
	for i := len(out) - 1; i >= 0; i-- {
		sum += out[i]
		out[i] = sum % 10
	}
	return out
}

func runPhase(nums []int) []int {
	results := []int{}
	var steps int = len(nums)

	for i := 0; i < steps; i++ {
		mults := getPatternStep(basePattern, i)
		stepSum := 0

		for j := 0; j < len(nums); j++ {
			nidx := j % len(nums)
			midx := j % len(mults)
			x := nums[nidx] * mults[midx]
			stepSum += x
		}
		results = append(results, int(math.Abs(float64(stepSum)))%10)
	}

	return results
}

func getPatternStep(pattern []int, step int) []int {
	pstep := 0
	stepPattern := []int{}
	tolen := (step + 1) * len(pattern)

	// add one extra so we can trim it from the front
	for i := 0; i < len(pattern)+1; i++ {
		for j := 0; j < step+1; j++ {
			stepPattern = append(stepPattern, pattern[pstep])
		}
		pstep = (pstep + 1) % len(pattern)
	}
	return stepPattern[1 : tolen+1]
}

func getNums(line string) []int {
	nums := []int{}
	for _, ch := range line {
		nums = append(nums, int(ch-'0'))
	}
	return nums
}
