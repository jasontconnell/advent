package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("--2021 day 10 solution--")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	scores := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
	match := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
	score := 0
	for _, s := range in {
		_, sc := getSyntaxError(s, scores, match)
		score += sc
	}
	return score
}

func part2(in input) output {
	scores := map[rune]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}
	match := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
	results := []int{}
	for _, line := range in {
		corrupt, _ := getSyntaxError(line, scores, match)
		if !corrupt {
			results = append(results, complete(line, scores, match))
		}
	}
	sort.Ints(results)
	score := results[len(results)/2]

	return score
}

func complete(line string, scores map[rune]int, match map[rune]rune) int {
	stack := []rune{}
	for _, c := range line {
		switch c {
		case ')', ']', '}', '>':
			stack = stack[1:]
		default:
			stack = append([]rune{c}, stack...)
		}
	}

	score := 0
	for j, _ := range stack {
		c := stack[j]
		m := match[c]
		score = 5*score + scores[m]
	}
	return score
}

func getSyntaxError(line string, scores map[rune]int, match map[rune]rune) (bool, int) {
	score := 0
	stack := []rune{}
	found := false
	revmatch := make(map[rune]rune)
	for k, v := range match {
		revmatch[v] = k
	}
	for _, c := range line {
		switch c {
		case ')', ']', '}', '>':
			m := revmatch[c]
			if stack[0] != m {
				score = scores[c]
				found = true
			}
			stack = stack[1:]
		default:
			stack = append([]rune{c}, stack...)
		}

		if found {
			break
		}
	}
	return found, score
}

func check(depths map[rune]int) (bool, rune) {
	var ch rune
	found := false
	for c, i := range depths {
		if i < 0 {
			found = true
			ch = c
			break
		}
	}
	return found, ch
}
