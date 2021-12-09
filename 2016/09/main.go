package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = string
type output = int

var reg = regexp.MustCompile("(\\d+)x(\\d+)\\)")

func main() {
	startTime := time.Now()

	in, err := common.ReadString(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("--2016 day 09 solution--")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return explode(in, false)
}

func part2(in input) output {
	return explode(in, true)
}

func explode(txt string, dec bool) int {
	result := 0
	for i := 0; i < len(txt); i++ {
		cs := string(txt[i])

		if cs == "(" {
			m := i + 10
			if i+m > len(txt) {
				m = len(txt) - i
			}
			if groups := reg.FindStringSubmatch(string(txt[i : i+m])); groups != nil && len(groups) > 1 {
				chars, _ := strconv.Atoi(groups[1])
				repeat, _ := strconv.Atoi(groups[2])

				i += len(groups[0])
				rep := txt[i+1 : i+chars+1]

				if dec && strings.Contains(rep, "(") {
					result += repeat * explode(rep, dec)
				} else {
					result += repeat * len(rep)
				}

				i += chars
			}
		} else {
			result++
		}
	}
	return result
}
