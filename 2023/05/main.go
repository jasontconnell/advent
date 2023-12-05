package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type fertmap struct {
	source string
	dest   string
	ranges []seedmap
}

type seedmap struct {
	dest   int
	source int
	length int
}

type categoryresult struct {
	category string
	value    int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 05 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	seeds, m := parseInput(in)
	r := mapValues(seeds, m)
	return findMin("location", r)
}

func part2(in input) output {
	return 0
}

func findMin(cat string, seedresults map[int][]categoryresult) int {
	min := math.MaxInt64
	for _, cr := range seedresults {
		for _, catval := range cr {
			if catval.category == cat && catval.value < min {
				min = catval.value
			}
		}
	}
	return min
}

func mapValues(seeds []int, m []fertmap) map[int][]categoryresult {
	results := map[int][]categoryresult{}

	for _, s := range seeds {
		val := s
		for _, f := range m {
			res := mapTo(val, f)

			if _, ok := results[s]; !ok {
				results[s] = []categoryresult{}
			}

			results[s] = append(results[s], res)
			val = res.value
		}
	}

	return results
}

func mapTo(val int, fm fertmap) categoryresult {
	cres := categoryresult{category: fm.dest}
	found := false
	for _, v := range fm.ranges {
		if val >= v.source && val < v.source+v.length {
			shift := val - v.source
			cres.value = v.dest + shift
			found = true
		}
	}

	if !found {
		cres.value = val
	}
	return cres
}

func parseInput(in input) ([]int, []fertmap) {
	seedstr := strings.Fields(strings.Replace(in[0], "seeds:", "", 1))
	seeds := []int{}
	for _, s := range seedstr {
		x, _ := strconv.Atoi(s)
		seeds = append(seeds, x)
	}

	freg := regexp.MustCompile(`([a-z]+)-to-([a-z]+) map:`)

	var fm fertmap
	var list []fertmap
	for i, line := range in {
		if i < 2 {
			continue
		}
		if line == "" || i == len(in)-1 {
			list = append(list, fm)
			fm = fertmap{}
			continue
		}

		matches := freg.FindStringSubmatch(line)
		if len(matches) > 0 {
			fm.source = matches[1]
			fm.dest = matches[2]
			continue
		}

		flds := strings.Fields(line)
		if len(flds) == 3 {
			dest, _ := strconv.Atoi(flds[0])
			src, _ := strconv.Atoi(flds[1])
			length, _ := strconv.Atoi(flds[2])

			fm.ranges = append(fm.ranges, seedmap{dest: dest, source: src, length: length})
		}
	}

	return seeds, list
}
