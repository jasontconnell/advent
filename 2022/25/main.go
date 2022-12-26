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
type output = string

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Printf("Time %v ", p1time)
}

func part1(in input) output {
	nums := parseInput(in)
	return asSnafu(sum(nums))
}

func sum(list []int) int {
	total := 0
	for _, n := range list {
		total += n
	}
	return total
}

func asSnafu(n int) string {
	x := n
	digits := []int{}
	for x > 0 {
		dm := x % 5
		digits = append([]int{dm}, digits...)
		x /= 5
	}

	for i := len(digits) - 1; i >= 0; i-- {
		nv := math.MinInt32
		if digits[i] == 3 {
			nv = -2
		} else if digits[i] == 4 {
			nv = -1
		} else if digits[i] == 5 {
			nv = 0
		}
		if nv != math.MinInt32 {
			digits[i] = nv
			if i == 0 {
				digits = append([]int{1}, digits...)
			} else {
				digits[i-1]++
			}
		}
	}

	s := ""
	for _, d := range digits {
		sd := strconv.Itoa(d)
		if d == -1 {
			sd = "-"
		} else if d == -2 {
			sd = "="
		}
		s += sd
	}
	return s
}

func getSnafu(s string) int {
	total := 0
	place := 1
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		x, err := strconv.Atoi(string(c))
		mult := x
		if err != nil {
			mult = -1
			if c == '=' {
				mult = -2
			}
		}
		total += mult * place

		place *= 5
	}
	return total
}

func parseInput(in input) []int {
	ns := []int{}
	for _, line := range in {
		ns = append(ns, getSnafu(line))
	}
	return ns
}
