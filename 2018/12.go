package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	//"strconv"
	"strings"
	//"math"
)

type rule struct {
	in  []bool
	out bool
}

var input = "12.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}

	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	pad := 50
	initial := parseInit(lines[0], pad)
	rules := parseRules(lines[2:])
	result := sim(initial, rules, 20)
	sum := sumIndices(result, pad)

	fmt.Println("part 1:", sum)

	fmt.Println("Time", time.Since(startTime))
}

func sumIndices(result []bool, pad int) int {
	sum := 0
	for i := 0; i < len(result); i++ {
		index := i - pad
		if result[i] {
			sum += index
		}
	}
	return sum
}

func sim(initial []bool, rules []rule, turns int) []bool {
	gen := initial
	for i := 0; i < turns; i++ {
		gen = processMatches(rules, gen)
	}
	return gen
}

func processMatches(rules []rule, val []bool) []bool {
	result := make([]bool, len(val))
	for i := 0; i < len(val)-5; i++ {
		chunk := val[i : i+5]
		for _, r := range rules {
			m := match(r, chunk)
			if m {
				result[i+2] = r.out
			}
		}
	}
	return result
}

func match(r rule, val []bool) bool {
	for i, b := range r.in {
		if val[i] != b {
			return false
		}
	}
	return true
}

func parseInit(init string, pad int) []bool {
	s := len("initial state: ")
	grid := parse(string(init[s:]))
	// pad out grid on both sides
	for i := 0; i < pad; i++ {
		grid = append([]bool{false}, grid...)
		grid = append(grid, false)
	}
	return grid
}

func parseRules(r []string) []rule {
	rules := []rule{}
	for _, s := range r {
		parts := strings.Split(s, " => ")
		rl := rule{in: parse(parts[0]), out: getVal(rune(parts[1][0]))}

		rules = append(rules, rl)
	}
	return rules
}

func parse(s string) []bool {
	b := make([]bool, len(s))
	for i, c := range s {
		b[i] = getVal(c)
	}
	return b
}

func getVal(r rune) bool {
	return r == '#'
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
