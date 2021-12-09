package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = string

type runeCount struct {
	ch    rune
	count int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
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
	signal, _ := getSignal(in)
	return signal
}

func part2(in input) output {
	_, signal := getSignal(in)
	return signal
}

func getSignal(in input) (output, output) {
	if len(in) == 0 {
		return "", ""
	}

	list := make([]map[rune]int, len(in[0]))
	for i := 0; i < len(list); i++ {
		list[i] = make(map[rune]int)
	}

	for _, line := range in {
		for i, c := range line {
			list[i][c]++
		}
	}

	maxes := make(map[int]runeCount)
	mins := make(map[int]runeCount)
	for col, m := range list {
		for r, i := range m {
			if rc, ok := maxes[col]; !ok {
				maxes[col] = runeCount{ch: r, count: i}
			} else {
				if i > rc.count {
					maxes[col] = runeCount{ch: r, count: i}
				}
			}

			if rc, ok := mins[col]; !ok {
				mins[col] = runeCount{ch: r, count: i}
			} else {
				if i < rc.count {
					mins[col] = runeCount{ch: r, count: i}
				}
			}
		}
	}

	max, min := "", ""
	for i := 0; i < len(list); i++ {
		if c, ok := maxes[i]; ok {
			max += string(c.ch)
		}

		if c, ok := mins[i]; ok {
			min += string(c.ch)
		}
	}
	return max, min
}
