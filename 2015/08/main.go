package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

var hex *regexp.Regexp = regexp.MustCompile(`\\x[abcdef0-9]{2}`)
var q *regexp.Regexp = regexp.MustCompile(`\\"`)
var bs *regexp.Regexp = regexp.MustCompile(`\\\\`)

var bse *regexp.Regexp = regexp.MustCompile(`\\`)
var qe *regexp.Regexp = regexp.MustCompile(`"`)

type input []string

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
	total, adjusted := getCount(in)
	return total - adjusted
}

func part2(in input) int {
	total, expanded := getCountExpand(in)
	return expanded - total
}

func getCount(in input) (int, int) {
	total, adjusted := 0, 0
	for _, line := range in {
		txt := line
		txt = txt[1 : len(txt)-1]
		txt = bs.ReplaceAllString(txt, `|`)
		txt = q.ReplaceAllString(txt, `|`)
		txt = hex.ReplaceAllString(txt, "|")

		total += len(line)
		adjusted += len(txt)
	}
	return total, adjusted
}

func getCountExpand(in []string) (int, int) {
	total, adjusted := 0, 0
	for _, line := range in {
		txt := line
		txt = bse.ReplaceAllString(txt, `||`)
		txt = qe.ReplaceAllString(txt, `|"`)

		total += len(line)
		adjusted += len(txt) + 2
	}
	return total, adjusted
}
