package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type dir int

type nodeinstr struct {
	node     string
	nodepair []string
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 08 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	dirs, instr := parseInput(in)
	return find(dirs, instr, "AAA", "ZZZ")
}

func part2(in input) output {
	dirs, instr := parseInput(in)
	return ghostTraverse(dirs, instr, "A", "Z")
}

func find(dirs []int, nodes []nodeinstr, start, search string) int {
	m := make(map[string]nodeinstr)
	for _, n := range nodes {
		m[n.node] = n
	}
	found := false
	cur := start
	steps := 0
	diridx := 0

	for !found {
		cdir := dirs[diridx]
		cnode := m[cur]
		cur = cnode.nodepair[cdir]
		steps++
		diridx = (diridx + 1) % len(dirs)
		found = strings.HasSuffix(cur, search)
	}
	return steps
}

func ghostTraverse(dirs []int, nodes []nodeinstr, startsuffix, searchsuffix string) int {
	ghosts := []string{}
	for _, s := range nodes {
		if strings.HasSuffix(s.node, startsuffix) {
			ghosts = append(ghosts, s.node)
		}
	}
	res := []int{}
	for _, g := range ghosts {
		steps := find(dirs, nodes, g, searchsuffix)
		res = append(res, steps)
	}

	return lcm(res[0], res[1], res[2:]...)
}

func parseInput(in input) ([]int, []nodeinstr) {
	dirs := []int{}
	for _, c := range in[0] {
		d := 1
		if c == 'L' {
			d = 0
		}
		dirs = append(dirs, d)
	}

	reg := regexp.MustCompile(`([A-Z]+) += +\(([A-Z]+), ?([A-Z]+)\)`)

	list := []nodeinstr{}
	for _, line := range in[2:] {
		m := reg.FindStringSubmatch(line)
		n := nodeinstr{node: m[1], nodepair: []string{m[2], m[3]}}
		list = append(list, n)
	}
	return dirs, list
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
