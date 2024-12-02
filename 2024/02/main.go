package main

import (
	"fmt"
	"log"
	"math"
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
	fmt.Fprintln(w, "--2024 day 02 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	rs := parse(in)
	return getSafes(rs, false)
}

func part2(in input) output {
	rs := parse(in)
	return getSafes(rs, true)
}

func getSafes(rs [][]int, scrub bool) int {
	c := 0
	for _, r := range rs {
		if !scrub {
			if isSafe(r) {
				c++
			}
		} else {
			vars := removeOne(r)
			for _, v := range vars {
				if isSafe(v) {
					c++
					break
				}
			}
		}
	}
	return c
}

func isSafe(rpt []int) bool {
	if !isConsistent(rpt) {
		return false
	}
	errs := 0
	for i, x := range rpt {
		if i == 0 {
			continue
		}
		prev := rpt[i-1]
		v := int(math.Abs(float64(prev - x)))
		if v == 0 {
			errs++
		}
		if v > 3 {
			errs++
		}
	}
	return errs == 0
}

func isConsistent(rpt []int) bool {
	sign := false
	errs := 0
	for i, x := range rpt {
		if i == 0 {
			sign = math.Signbit(float64(x - rpt[1]))
			continue
		}
		diff := rpt[i-1] - x
		if diff == 0 {
			errs++
		}
		nsign := math.Signbit(float64(diff))
		if sign != nsign {
			errs++
		}
	}
	return errs == 0
}

func removeOne(rpt []int) [][]int {
	vars := [][]int{}
	for x := 0; x < len(rpt); x++ {
		d := []int{}
		for i := 0; i < len(rpt); i++ {
			if x == i {
				continue
			}
			d = append(d, rpt[i])
		}
		vars = append(vars, d)
	}
	return vars
}

func parse(lines []string) [][]int {
	val := [][]int{}
	for _, line := range lines {
		sp := strings.Fields(line)
		x := []int{}
		for _, s := range sp {
			v, _ := strconv.Atoi(s)
			x = append(x, v)
		}
		val = append(val, x)
	}
	return val
}
