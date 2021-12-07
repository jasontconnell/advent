package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []int
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadIntCsv(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return getOptimalFuel(in, false)
}

func part2(in input) output {
	return getOptimalFuel(in, true)
}

func getOptimalFuel(crabs input, acc bool) int {
	min, max := minmax(crabs)
	optfuel := math.MaxInt32
	for i := min; i < max; i++ {
		freq := 0
		for _, v := range crabs {
			diff := i - v
			if i > v {
				diff = -diff
			}

			abs := int(math.Abs(float64(diff)))
			if acc {
				abs += doAcc(abs)
			}
			freq += abs
		}

		if freq < optfuel {
			optfuel = freq
		}
	}
	return optfuel
}

func doAcc(diff int) int {
	if diff == 0 {
		return 0
	}

	acc := 0
	for i := diff - 1; i > 0; i-- {
		acc += i
	}
	return acc
}

func minmax(in input) (int, int) {
	min := math.MaxInt32
	max := math.MinInt32

	for _, v := range in {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
