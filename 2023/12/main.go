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

type record struct {
	conditions  []condition
	arrangement []int
}

type condition int

const (
	Damaged condition = iota
	Operational
	Unknown
	Undefined
)

type key struct {
	x, y, z int
}

func (c condition) String() string {
	s := '_'
	switch c {
	case Unknown:
		s = '?'
	case Operational:
		s = '.'
	case Damaged:
		s = '#'
	}
	return string(s)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 12 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	list := parseInput(in)
	return solve(list, 1)
}

func part2(in input) output {
	list := parseInput(in)
	return solve(list, 5)
}

func solve(list []record, mult int) int {
	total := 0
	for _, r := range list {
		pattern := r.conditions
		groups := r.arrangement

		if mult > 1 {
			pattern = multiply(mult, pattern, Unknown, true)
			groups = multiply(mult, groups, 0, false)
		}

		m := make(map[key]int)
		total += getPermutations(m, pattern, groups, 0, 0, 0)
	}
	return total
}

// use dynamic programming
func getPermutations(dp map[key]int, pattern []condition, groups []int, pidx, gidx, curlen int) int {
	if pidx == len(pattern) {
		if (gidx == len(groups)-1 && groups[gidx] == curlen) || (gidx == len(groups) && curlen == 0) {
			return 1
		}
		return 0
	}

	k := key{pidx, gidx, curlen}
	if v, ok := dp[k]; ok {
		return v
	}

	total := 0
	c := pattern[pidx]

	// if damaged or unknown, we're in a group, move to the next pattern index and
	// increment the current length by 1
	if c == Unknown || c == Damaged {
		total += getPermutations(dp, pattern, groups, pidx+1, gidx, curlen+1)
	}

	// if damaged or unknown
	// if we are at the end of the current group length, start a new group
	// if the current group is empty, move the pattern to the next one
	if c == Unknown || c == Operational {
		if curlen > 0 && gidx < len(groups) && groups[gidx] == curlen {
			total += getPermutations(dp, pattern, groups, pidx+1, gidx+1, 0)
		}

		if curlen == 0 {
			total += getPermutations(dp, pattern, groups, pidx+1, gidx, 0)
		}
	}
	dp[k] = total
	return total
}

func multiply[T any](num int, list []T, sep T, dosep bool) []T {
	ret := []T{}
	for i := 0; i < num; i++ {
		ret = append(ret, list...)
		if i < num-1 && dosep {
			ret = append(ret, sep)
		}
	}

	return ret
}

func parseInput(in input) []record {
	records := []record{}
	for _, line := range in {
		sp := strings.Split(line, " ")
		r := record{}
		for _, c := range sp[0] {
			switch c {
			case '#':
				r.conditions = append(r.conditions, Damaged)
			case '.':
				r.conditions = append(r.conditions, Operational)
			case '?':
				r.conditions = append(r.conditions, Unknown)
			}
		}

		arrs := strings.Split(sp[1], ",")
		for _, s := range arrs {
			x, _ := strconv.Atoi(s)
			r.arrangement = append(r.arrangement, x)
		}

		records = append(records, r)
	}
	return records
}
