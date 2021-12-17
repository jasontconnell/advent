package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int64

type IPRange struct {
	Low, High int64
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	ips := parseInput(in)
	merged := mergeRanges(ips)
	return merged[0].High + 1
}

func part2(in input) output {
	ips := parseInput(in)
	merged := mergeRanges(ips)
	return countValid(merged)
}

func countValid(merged []IPRange) int64 {
	var count int64 = 0
	for i := 0; i < len(merged); i++ {
		count += (merged[i].High + 1 - merged[i].Low) // inclusive of high for total count
	}
	return 4294967295 + 1 - count
}

func mergeRanges(ipranges []IPRange) []IPRange {
	merged := []IPRange{}
	for _, ip := range ipranges {
		extended := false
		for i := 0; i < len(merged); i++ {
			m := merged[i]
			if withinRange(ip.Low, m.Low, m.High) && withinRange(ip.High, m.Low, m.High) { // fully contained already
				extended = true
				break
			}

			if !withinRange(ip.Low, m.Low, m.High) && withinRange(ip.High, m.Low, m.High) { // merged low needs extension down to ip low
				merged[i].Low = ip.Low
				extended = true
				break
			} else if withinRange(ip.Low, m.Low, m.High) && !withinRange(ip.High, m.Low, m.High) { // merged high needs extension to ip high
				merged[i].High = ip.High
				extended = true
				break
			}
		}

		if len(merged) == 0 || !extended {
			merged = append(merged, ip)
		}
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Low < merged[j].Low
	})

	return merged
}

func withinRange(num, low, high int64) bool {
	return num >= low-1 && num <= high+1
}

func parseInput(in input) []IPRange {
	reg := regexp.MustCompile("^([0-9]+)-([0-9]+)$")
	ranges := []IPRange{}
	for _, line := range in {
		if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
			low, _ := strconv.ParseInt(groups[1], 10, 64)
			high, _ := strconv.ParseInt(groups[2], 10, 64)

			ranges = append(ranges, IPRange{Low: low, High: high})
		}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Low < ranges[j].Low
	})
	return ranges
}
