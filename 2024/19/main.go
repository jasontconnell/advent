package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type state struct {
	pos    int
	sub    string
	total  string
	pieces string
}

type statekey struct {
	pos    int
	sub    string
	pieces string
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	p, d := parse(in)
	return findValidDesigns(p, d, false)
}

func part2(in input) output {
	p, d := parse(in)
	return findValidDesigns(p, d, true)
}

func findValidDesigns(pmap map[string]bool, designs []string, countall bool) int {
	lookup := make(map[string]int)
	min, max := minmaxlen(pmap)
	count := 0
	for _, design := range designs {
		var c int
		c, lookup = countVariations(pmap, design, min, max, len(design), lookup)
		if c == 0 {
			continue
		}
		if !countall && c > 1 {
			count++
		} else if x, ok := lookup[design]; ok && countall {
			count += x
		}
	}
	return count
}

func minmaxlen(pmap map[string]bool) (int, int) {
	min, max := math.MaxInt32, math.MinInt32
	for s := range pmap {
		if len(s) < min {
			min = len(s)
		}
		if len(s) > max {
			max = len(s)
		}
	}
	return min, max
}

func countVariations(pmap map[string]bool, design string, min, max, orig int, lookup map[string]int) (int, map[string]int) {
	w, ok := lookup[design]
	if ok {
		return w, lookup
	}

	count := 0
	for k := range pmap {
		if k == design {
			count++
		} else if strings.HasPrefix(design, k) {
			var local int
			local, lookup = countVariations(pmap, design[len(k):], min, max, orig, lookup)
			count += local
		}
	}
	lookup[design] = count
	return lookup[design], lookup
}

func parse(in []string) (map[string]bool, []string) {
	patterns := strings.Fields(strings.ReplaceAll(in[0], ",", ""))
	designs := in[2:]
	pmap := make(map[string]bool)
	for _, p := range patterns {
		pmap[p] = true
	}

	return pmap, designs
}
