package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

var bad []string = []string{"ab", "cd", "pq", "xy"}
var vowels = []rune{'a', 'e', 'i', 'o', 'u'}

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) int {
	return getGoodCount(in)
}

func part2(in input) int {
	return getPairs(in)
}

func getGoodCount(lines input) int {
	goodCount := 0
	for _, txt := range lines {
		if !badCheck(txt, bad) && appearanceCheck(txt, vowels, 3) && doubleCheck(txt) {
			goodCount++
		}
	}
	return goodCount
}

func getPairs(lines input) int {
	goodCount := 0
	allCount := 0

	for _, txt := range lines {
		if pairCheck(txt) && wrapCheck(txt) {
			goodCount++
		}

		allCount++
	}

	return goodCount
}

func appearanceCheck(s string, runes []rune, count int) bool {
	total := 0
	for _, r := range runes {
		total += strings.Count(s, string(r))
	}
	return total >= count
}

func doubleCheck(s string) bool {
	double := false

	for i, _ := range s {
		if i > 0 && s[i-1] == s[i] {
			double = true
			break
		}
	}
	return double
}

func badCheck(s string, bad []string) bool {
	isBad := false
	for _, b := range bad {
		if !isBad {
			isBad = strings.Contains(s, b)
		}
	}
	return isBad
}

func pairCheck(s string) bool {
	pairs := 0
	for i, _ := range s {
		if i > 0 {
			overlaps := false
			current := string(s[i-1]) + string(s[i])
			c := strings.Count(s[i:], current)
			if s[i-1] == s[i] && i < len(s)-1 {
				if s[i] == s[i+1] {
					overlaps = true
				}
			}

			if c > 0 && !overlaps {
				pairs++
			}
		}
	}
	return pairs > 0
}

func wrapCheck(s string) bool {
	wrap := false
	for i, _ := range s {
		if i > 1 && s[i-2] == s[i] {
			wrap = true
			break
		}
	}
	return wrap
}
