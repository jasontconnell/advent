package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 06 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	nums, ops := parseInput(in)
	return calcAll(nums, ops, false)
}

func part2(in input) output {
	nums, ops := parseInputPart2(in)
	return calcAll(nums, ops, true)
}

func calcAll(nums [][]int, ops []string, rev bool) int {
	total := 0

	idx := 0
	incr := 1
	done := false
	if rev {
		idx = len(nums) - 1
		incr = -1
	}

	for !done {
		var cur []int
		op := ops[idx]

		if !rev {
			for j := 0; j < len(nums); j++ {
				cur = append(cur, nums[j][idx])
			}
		} else {
			cur = nums[idx]
		}

		total += calculate(cur, op)

		idx += incr
		if rev {
			done = idx < 0
		} else {
			done = idx >= len(nums[0])
		}
	}
	return total
}

func calculate(nums []int, op string) int {
	total := 0
	if op == "*" {
		total = 1
	}

	for i := 0; i < len(nums); i++ {
		switch op {
		case "+":
			total += nums[i]
		case "*":
			total *= nums[i]
		}
	}
	return total
}

func parseInput(in input) ([][]int, []string) {
	nums := [][]int{}
	ops := []string{}
	for _, line := range in {
		flds := strings.Fields(line)
		if flds[0] == "+" || flds[0] == "*" {
			ops = flds
		} else {
			nline := []int{}
			for _, f := range flds {
				x, _ := strconv.Atoi(f)
				nline = append(nline, x)
			}
			nums = append(nums, nline)
		}
	}
	return nums, ops
}

func parseInputPart2(in input) ([][]int, []string) {
	nums := [][]int{}
	ops := []string{}

	cur := []int{}
	for j := 0; j < len(in[0]); j++ {
		digits := []int{}
		allblank := true
		for i := 0; i < len(in); i++ {
			isop := in[i][j] == '+' || in[i][j] == '*'
			if in[i][j] != ' ' && !isop {
				x, _ := strconv.Atoi(string(in[i][j]))
				digits = append(digits, x)
				allblank = false
			} else if isop || i == len(in)-1 {
				if in[i][j] != ' ' {
					ops = append(ops, string(in[i][j]))
				}

				if len(digits) > 0 {
					val := getNum(digits)
					cur = append(cur, val)
				}
			}

			if !allblank {
				allblank = i == len(in)-1 && j == len(in[0])-1
			}
		}

		if allblank {
			if len(cur) > 0 {
				nums = append(nums, cur)
				cur = []int{}
			}
		}
	}

	return nums, ops
}

func getNum(n []int) int {
	pow := 0
	val := 0
	for i := len(n) - 1; i >= 0; i-- {
		val += n[i] * int(math.Pow10(pow))
		pow++
	}
	return val
}
