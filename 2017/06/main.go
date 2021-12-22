package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 06 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	vals := parseInput(in)
	c, _ := solve(vals)
	return c
}

func part2(in input) output {
	vals := parseInput(in)
	_, c := solve(vals)
	return c
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

func parseInput(in input) []int {
	flds := strings.Fields(in)
	vals := []int{}

	for _, f := range flds {
		val, _ := strconv.Atoi(f)
		vals = append(vals, val)
	}
	return vals
}
