package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	v := getCalibrations(in)
	return sum(v)
}

func part2(in input) output {
	v := getWordCalibrations(in)
	return sum(v)
}

func sum(ints []int) int {
	s := 0
	for _, v := range ints {
		s += v
	}
	return s
}

func getCalibrations(list []string) []int {
	ints := []int{}
	for _, s := range list {
		f := -1
		l := -1
		for i := 0; i < len(s); i++ {
			cs := s[i]
			ce := s[len(s)-i-1]
			if f == -1 && strings.ContainsAny(string(cs), "0123456789") {
				f, _ = strconv.Atoi(string(cs))
			}
			if l == -1 && strings.ContainsAny(string(ce), "0123456789") {
				l, _ = strconv.Atoi(string(ce))
			}
		}

		ints = append(ints, f*10+l)
	}

	return ints
}

func getWordCalibrations(list []string) []int {
	ints := []int{}
	for _, s := range list {
		f := -1
		l := -1
		for i := 0; i < len(s); i++ {
			cs := s[i]
			ce := s[len(s)-i-1]
			if f == -1 && strings.ContainsAny(string(cs), "0123456789") {
				f, _ = strconv.Atoi(string(cs))
			} else if f == -1 {
				for k, v := range nums {
					if strings.HasPrefix(s[i:], k) {
						f = v
					}
				}
			}
			if l == -1 && strings.ContainsAny(string(ce), "0123456789") {
				l, _ = strconv.Atoi(string(ce))
			} else if l == -1 {
				for k, v := range nums {
					if strings.HasPrefix(s[len(s)-i-1:], k) {
						l = v
					}
				}
			}
		}

		ints = append(ints, f*10+l)
	}

	return ints
}

var nums map[string]int = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
	"zero":  0,
}

func getNumber(s string) int {
	if v, ok := nums[s]; ok {
		return v
	}
	return -1
}
