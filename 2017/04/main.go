package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 04 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return countValid(in, false)
}

func part2(in input) output {
	return countValid(in, true)
}

func countValid(in input, anagram bool) int {
	c := 0
	for _, line := range in {
		if !anagram && isValid(line) {
			c++
		} else if anagram && hasNoAnagrams(line) {
			c++
		}
	}
	return c
}

func isValid(line string) bool {
	sp := strings.Split(line, " ")
	sort.Strings(sp)

	valid := true
	for i, s := range sp {
		if i > 0 && s == sp[i-1] {
			valid = false
			break
		}
	}
	return valid
}

func hasNoAnagrams(line string) bool {
	sp := strings.Split(line, " ")
	strs := []string{}

	for i := 0; i < len(sp); i++ {
		runes := make([]rune, len(sp[i]))
		for j, c := range sp[i] {
			if c > 96 && c < 123 || c > 64 && c < 91 {
				runes[j] = c
			}
		}

		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})

		str := ""
		for _, c := range runes {
			str += string(c)
		}

		strs = append(strs, str)
	}

	sort.Strings(strs)

	valid := true
	for i := range strs {
		if i > 0 && strs[i] == strs[i-1] {
			valid = false
			break
		}
	}
	return valid
}
