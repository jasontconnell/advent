package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type bank struct {
	val, pos int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 03 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	batteries := parseInput(in)
	sum := 0
	for _, b := range batteries {
		sum += findLargest(b, 2)
	}
	return sum
}

func part2(in input) output {
	batteries := parseInput(in)
	sum := 0
	for _, b := range batteries {
		sum += findLargest(b, 12)
	}
	return sum
}

func findLargest(batteries []int, n int) int {
	res := []int{}
	lidx := 0
	for idx := 0; idx < n; idx++ {
		cmax := 0
		end := len(batteries) - (n - len(res) - 1)
		for i := lidx; i < end; i++ {
			x := batteries[i]
			if x > cmax {
				cmax = x
				lidx = i
			}
		}

		if cmax > 0 {
			res = append(res, cmax)
			lidx++
		}
	}

	ans := 0
	pow := 0
	for i := len(res) - 1; i >= 0; i-- {
		ans += res[i] * int(math.Pow10(pow))
		pow++
	}

	return ans
}

func parseInput(in input) [][]int {
	list := [][]int{}
	for _, line := range in {
		ls := []int{}
		for _, c := range line {
			x, _ := strconv.Atoi(string(c))
			ls = append(ls, x)
		}
		list = append(list, ls)
	}
	return list
}
