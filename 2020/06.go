package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "06.txt"

type group struct {
	qmap map[string][]int
	yes  map[string]int
	size int
}

func newGroup() *group {
	g := new(group)
	g.qmap = make(map[string][]int)
	g.yes = make(map[string]int)
	return g
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	groups := getGroups(lines)
	p1 := 0
	p2 := 0
	for _, grp := range groups {
		p1 += len(grp.qmap)

		for q, _ := range grp.qmap {
			if grp.size == grp.yes[q] {
				p2++
			}
		}
	}

	fmt.Println("part 1:", p1)
	fmt.Println("part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func getGroups(lines []string) []*group {
	cgrp := newGroup()
	groups := []*group{}
	for i, line := range lines {

		for _, c := range line {
			cs := string(c)
			cgrp.qmap[cs] = append(cgrp.qmap[cs], int(c))
			cgrp.yes[cs]++
		}

		if line != "" && i < len(lines) {
			cgrp.size++
		}

		done := i == len(lines)-1
		if line == "" || done {
			groups = append(groups, cgrp)
			if !done {
				cgrp = newGroup()
			}
		}
	}

	return groups
}
