package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type cave struct {
	id          string
	connections map[string]*cave
	big         bool
}

func (c *cave) String() string {
	return fmt.Sprintf("id: %s, big: %v", c.id, c.big)
}

type state struct {
	path   []*cave
	visits map[string]int
	goal   bool
	loc    *cave
}

func (s *state) printPath() string {
	val := ""
	for _, p := range s.path {
		val += p.id + "=>"
	}
	return strings.TrimSuffix(val, "=>")
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("--2021 day 12 solution--")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	caves := parseInput(in)
	paths := findPaths(caves, false)
	return len(paths)
}

func part2(in input) output {
	caves := parseInput(in)
	paths := findPaths(caves, true)
	return len(paths)
}

func findPaths(list map[string]*cave, modify bool) [][]*cave {
	paths := [][]*cave{}
	start := []*cave{}
	for _, c := range list {
		if c.id == "start" {
			start = append(start, c)
		}
	}

	queue := []*state{}
	for _, s := range start {
		queue = append(queue, newState(nil, s))
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.loc.id == "end" {
			paths = append(paths, cur.path)
			continue
		}

		for _, cn := range cur.loc.connections {
			if _, ok := cur.visits[cn.id]; !ok || (modify && cn.id != "start" && cn.id != "end" && mmax(cur.visits) == 1) {
				st := newState(cur, cn)
				queue = append(queue, st)
			}
		}
	}
	return paths
}

func newState(cur *state, c *cave) *state {
	st := &state{loc: c, visits: make(map[string]int)}
	if cur != nil {
		st.path = make([]*cave, len(cur.path))
		copy(st.path, cur.path)
		for k, v := range cur.visits {
			st.visits[k] = v
		}
	}
	st.path = append(st.path, c)
	if !c.big {
		st.visits[c.id]++
	}

	return st
}

func mmax(m map[string]int) int {
	max := 0
	for _, i := range m {
		if i > max {
			max = i
		}
	}
	return max
}

func parseInput(in input) map[string]*cave {
	m := make(map[string]*cave)

	for _, s := range in {
		sp := strings.Split(s, "-")
		lid, rid := sp[0], sp[1]

		if c, ok := m[lid]; !ok {
			lbig := strings.ToUpper(lid) == lid
			c = &cave{id: lid, connections: make(map[string]*cave), big: lbig}
			m[c.id] = c
		}

		if c, ok := m[rid]; !ok {
			rbig := strings.ToUpper(rid) == rid
			c = &cave{id: rid, connections: make(map[string]*cave), big: rbig}
			m[c.id] = c
		}
	}
	for _, s := range in {
		sp := strings.Split(s, "-")
		lid, rid := sp[0], sp[1]

		l := m[lid]
		r := m[rid]
		l.connections[rid] = r
		r.connections[lid] = l
	}
	return m
}
