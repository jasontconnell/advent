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
	fmt.Fprintln(w, "--2023 day 09 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	all := parseInput(in)
	return extractAllHistory(all)
}

func part2(in input) output {
	return 0
}

func extractAllHistory(all [][]int) int {
	s := 0
	for _, list := range all {
		lvl := make([]int, len(list)-1)
		lst := extractHistory(list, lvl)
		s += list[len(list)-1] + lst
	}
	return s
}

func extractHistory(vals []int, slvl []int) int {
	cp := make([]int, len(vals))
	copy(cp, vals)
	v := make(map[int]int)
	for i := 0; i < len(cp)-1; i++ {
		x := cp[i+1] - cp[i]
		slvl[i] = x
		v[x]++
	}

	if v[0] == len(slvl) {
		return 0
	}

	nhst := make([]int, len(slvl)-1)
	x := extractHistory(slvl, nhst)
	last := slvl[len(slvl)-1] + x

	return last
}

func parseInput(in input) [][]int {
	ret := [][]int{}
	for _, line := range in {
		lvals := []int{}
		ss := strings.Fields(line)
		for _, s := range ss {
			x, _ := strconv.Atoi(s)
			lvals = append(lvals, x)
		}
		ret = append(ret, lvals)
	}
	return ret
}
