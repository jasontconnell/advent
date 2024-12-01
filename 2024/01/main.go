package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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
	fmt.Fprintln(w, "--2024 day 01 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	c1, c2 := parse(in)
	sort.Ints(c1)
	sort.Ints(c2)
	return getDist(c1, c2)
}

func part2(in input) output {
	c1, c2 := parse(in)
	sort.Ints(c1)
	sort.Ints(c2)
	return getSim(c1, c2)
}

func getDist(c1, c2 []int) int {
	td := 0
	for i := 0; i < len(c1); i++ {
		d := int(math.Abs(float64(c1[i] - c2[i])))
		td += d
	}
	return td
}

func getSim(c1, c2 []int) int {
	m2 := make(map[int]int)
	sims := make(map[int]int)
	for _, c := range c2 {
		m2[c]++
	}

	for _, k := range c1 {
		if v2, ok := m2[k]; ok {
			sims[k] += k * v2
		}
	}

	sum := 0
	for _, v := range sims {
		sum += v
	}
	return sum
}

func parse(lines []string) ([]int, []int) {
	c1, c2 := []int{}, []int{}
	for _, line := range lines {
		sp := strings.Fields(line)
		if len(sp) != 2 {
			continue
		}
		i1, _ := strconv.Atoi(sp[0])
		i2, _ := strconv.Atoi(sp[1])
		c1 = append(c1, i1)
		c2 = append(c2, i2)
	}
	return c1, c2
}
