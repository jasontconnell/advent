package main

import (
	"fmt"
	"log"
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
	return sim(in, 80)
}

func part2(in input) output {
	return sim(in, 256)
}

func sim(fish []int, days int) int {
	m := make(map[int]int)

	for _, f := range fish {
		m[f]++
	}

	for i := 0; i < days; i++ {
		add := m[0]
		for age := 1; age < 9; age++ {
			m[age-1] = m[age]
		}
		m[6] += add
		m[8] = add
	}

	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}
