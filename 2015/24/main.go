package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []int
type output int

func main() {
	startTime := time.Now()

	in, err := common.ReadInts(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	qe := getQuantumEntanglement(in, 3)
	sort.Ints(qe)
	return output(qe[0])
}

func part2(in input) output {
	qe := getQuantumEntanglement(in, 4)
	sort.Ints(qe)
	return output(qe[0])
}

func getQuantumEntanglement(in input, div int) []int {
	sort.Ints(in)
	sum := 0
	for _, v := range in {
		sum += v
	}
	res := []int{}
	each := sum / div
	for i := 0; i < len(in); i++ {
		cp := make(input, len(in))
		copy(cp, in)
		cp = append(cp[:i], cp[i+1:]...)
		filled := Fill(cp, 0, each)
		product := 1
		for _, num := range filled {
			product = product * num
		}
		res = append(res, product)
	}
	return res
}

func Fill(list []int, bucket, max int) []int {
	nums := []int{}

	for i := len(list) - 1; i >= 0; i-- {
		if bucket+list[i] <= max {
			bucket += list[i]
			nums = append(nums, list[i])
		}
	}

	return nums
}
