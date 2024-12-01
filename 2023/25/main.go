package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type node struct {
	name        string
	connnames   []string
	connections map[string]*node
}

type state struct {
	name  string
	steps int
}

func (n *node) String() string {
	cnns := "["
	for _, c := range n.connections {
		cnns += c.name + ","
	}
	cnns = strings.TrimRight(cnns, ",") + "]"
	return fmt.Sprintf("%s -> %s", n.name, cnns)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 25 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Printf("Time %v", p1time)
}

func part1(in input) output {
	g := parseInput(in)
	for _, v := range g {
		fmt.Println(v.name, len(v.connections), v.connections)
	}
	fmt.Println(len(g))
	return 0
}

func determineCuts(m map[string]*node) (int, int) {
	return 0, 0
}

func parseInput(in input) map[string]*node {
	nodes := []*node{}
	for _, line := range in {
		sp := strings.Split(line, ":")

		nname := sp[0]
		cnames := strings.Fields(sp[1])

		n := node{name: nname, connections: make(map[string]*node)}
		for _, c := range cnames {
			n.connnames = append(n.connnames, c)
		}
		nodes = append(nodes, &n)
	}

	m := make(map[string]*node)
	for _, n := range nodes {
		m[n.name] = n
	}

	for _, n := range m {
		for _, c := range n.connnames {
			cn, ok := m[c]
			if !ok {
				cn = &node{name: c, connections: make(map[string]*node)}
			}
			if _, ok := n.connections[c]; !ok {
				n.connections[c] = cn
			}
			if _, ok := cn.connections[n.name]; !ok {
				cn.connections[n.name] = n
			}
			m[c] = cn
		}
	}

	return m
}
