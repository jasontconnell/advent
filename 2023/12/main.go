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
	return solve(list)
}

func part2(in input) output {
	return 0
}

func solve(list []record) int {
	total := 0
	for _, r := range list {
		total += getPermutations(r)
	}
	return total
}

func sum(i []int) int {
	s := 0
	for _, x := range i {
		s += x
	}
	return s
}

func group(list []condition) [][]condition {
	groups := [][]condition{}
	last := Undefined
	for _, c := range list {
		if c != last {
			g := []condition{c}
			groups = append(groups, g)
		} else {
			groups[len(groups)-1] = append(groups[len(groups)-1], c)
		}
		last = c
	}
	return groups
}

func firstoccurrence(c condition, list []condition) int {
	idx := -1
	for i := range list {
		if list[i] == c {
			idx = i
			break
		}
	}
	return idx
}

func occurrences(c condition, list []condition) int {
	x := 0
	for _, cc := range list {
		if cc == c {
			x++
		}
	}
	return x
}

func replace(c condition, list []condition, idx int) []condition {
	cp := make([]condition, len(list))
	copy(cp, list)
	cp[idx] = c
	return cp
}

func getPermutations(r record) int {
	maxperms := 0

	queue := [][]condition{}
	queue = append(queue, r.conditions)
	permutations := [][]condition{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		fst := firstoccurrence(Unknown, cur)
		if fst == -1 {
			permutations = append(permutations, cur)
			continue
		}

		// permutate
		n1 := replace(Operational, cur, fst)
		n2 := replace(Damaged, cur, fst)
		queue = append(queue, n1, n2)
	}

	for _, p := range permutations {
		if occurrences(Damaged, p) != sum(r.arrangement) {
			continue
		}
		gg := group(p)
		pass := true
		gidx := 0
		for _, arr := range r.arrangement {
			if gidx >= len(gg)-1 {
				break
			}
			for gg[gidx][0] != Damaged {
				gidx++
			}
			if len(gg[gidx]) != arr {
				pass = false
			}
			gidx++
		}
		if pass {
			maxperms++
		}
	}
	return maxperms
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
