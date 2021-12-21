package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []int
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadDigits(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 01 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return getSum(in, false)
}

func part2(in input) output {
	return getSum(in, true)
}

func getSum(digits []int, fromAcross bool) int {
	sum := 0

	var nextIndex func(x int, l int) int
	if fromAcross {
		nextIndex = func(x int, l int) int {
			return ((l / 2) + x) % l
		}
	} else {
		nextIndex = func(x int, l int) int {
			ix := x - 1
			if ix < 0 {
				ix = l - 1
			}
			return ix
		}
	}

	for i := 0; i < len(digits); i++ {
		idx := i
		cur := digits[idx]
		cmpidx := nextIndex(idx, len(digits))
		cmp := digits[cmpidx]
		if cur == cmp {
			sum += cmp
		}
	}

	return sum
}
