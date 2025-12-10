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

type machine struct {
	diagram int
	toggles []int
	reqs    []int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 10 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	parseInput(in)
	log.Println(getBinaryPos([]int{0, 1, 2, 3}))
	return 0
}

func part2(in input) output {
	return 0
}

func parseInput(in input) []machine {
	machines := []machine{}
	for _, line := range in {
		sp := strings.Fields(line)
		diagram := []int{}
		for _, c := range sp[0][1 : len(sp[0])-1] {
			d := 1
			if c == '.' {
				d = 0
			}
			diagram = append(diagram, d)
		}

		toggles := []int{}
		for i := 1; i < len(sp)-1; i++ {
			pos := []int{}
			for _, c := range sp[i] {
				if c == '(' || c == ')' || c == ',' {
					continue
				}
				x, _ := strconv.Atoi(string(c))
				pos = append(pos, x)
			}
			toggles = append(toggles, getBinaryPos(pos))
		}

		last := sp[len(sp)-1]
		csv := strings.Split(last[1:len(last)-1], ",")
		reqs := []int{}
		for _, c := range csv {
			x, _ := strconv.Atoi(c)
			reqs = append(reqs, x)
		}

		dint := getBinary(diagram)
		m := machine{diagram: dint, toggles: toggles, reqs: reqs}

		machines = append(machines, m)
	}
	return machines
}

func getBinary(d []int) int {
	b := 0
	cur := 0
	for i := len(d) - 1; i >= 0; i-- {
		b += d[i] << cur
		cur++
	}
	return b
}

func getBinaryPos(pos []int) int {
	b := 0
	for i := len(pos) - 1; i >= 0; i-- {
		b += 1 << pos[i]
	}
	return b
}
