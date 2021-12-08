package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type signal struct {
	patterns []string
	outvalue []string
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
	signals := parseInput(in)
	return count1478(signals)
}

func part2(in input) output {
	signals := parseInput(in)
	vals := getOutputs(signals)
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return sum
}

func count1478(sigs []signal) output {
	c := 0
	for _, sig := range sigs {
		for _, out := range sig.outvalue {
			if len(out) == 2 || len(out) == 4 || len(out) == 3 || len(out) == 7 {
				c++
			}
		}
	}
	return c
}

func getOutputs(signals []signal) []output {
	outs := []output{}
	for _, sig := range signals {
		digits := determineDigits(sig)

		v := ""
		for _, s := range sig.outvalue {
			if i, ok := digits[s]; ok {
				v += strconv.Itoa(i)
			}
		}
		vi, _ := strconv.Atoi(v)
		outs = append(outs, vi)
	}
	return outs
}

func determineDigits(sig signal) map[string]int {
	var one, four, seven, eight string
	var zero, two, three, five, six, nine string

	unified := append(sig.patterns, sig.outvalue...)
	for _, s := range unified {
		switch len(s) {
		case 2:
			one = s
		case 4:
			four = s
		case 3:
			seven = s
		case 7:
			eight = s
		}
	}

	digits := map[string]int{one: 1, four: 4, seven: 7, eight: 8}
	maxloops, loops := 15, 0
	for loops < maxloops {
		for _, s := range unified {
			if _, ok := digits[s]; ok {
				continue
			}

			switch len(s) {
			case 5:
				if three == "" && numcommon(four, s) == 3 && numcommon(seven, s) == 3 && numcommon(one, s) == 2 {
					three = s
					digits[three] = 3
				}
				if five == "" && numcommon(four, s) == 3 && numcommon(one, s) == 1 && numcommon(seven, s) == 2 {
					five = s
					digits[five] = 5
				}
				if two == "" && numcommon(four, s) == 2 && numcommon(one, s) == 1 && numcommon(seven, s) == 2 {
					two = s
					digits[two] = 2
				}
			case 6:
				if nine == "" && numcommon(four, s) == 4 {
					nine = s
					digits[nine] = 9
				}
				if zero == "" && numcommon(four, s) == 3 && numcommon(eight, s) == 6 && numcommon(one, s) == 2 {
					zero = s
					digits[zero] = 0
				}
				if six == "" && numcommon(one, s) == 1 {
					six = s
					digits[six] = 6
				}
			}
		}
		loops++
		if loops > maxloops {
			break
		}
	}

	return digits
}

func sortDigits(s string) string {
	var ch []byte = []byte(s)
	sort.Slice(ch, func(i, j int) bool {
		return ch[i] < ch[j]
	})
	return string(ch)
}

func numcommon(s1, s2 string) int {
	c := 0

	m1 := make(map[rune]rune)

	for _, ch := range s1 {
		m1[ch] = ch
	}

	for _, ch := range s2 {
		if _, ok := m1[ch]; ok {
			c++
		}
	}

	return c
}

func parseInput(in input) []signal {
	reg := regexp.MustCompile("([a-g]+)")
	signals := []signal{}
	for _, s := range in {
		sig := signal{}
		sp := strings.Split(s, "|")
		for i, p := range sp {
			m := reg.FindAllStringSubmatch(p, -1)
			for _, g := range m {
				if i == 0 {
					sig.patterns = append(sig.patterns, sortDigits(g[1]))
				} else {
					sig.outvalue = append(sig.outvalue, sortDigits(g[1]))
				}
			}
		}
		signals = append(signals, sig)
	}
	return signals
}
