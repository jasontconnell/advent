package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "08.txt"

type node struct {
	id       int
	ccount   int
	mcount   int
	children []*node
	parent   *node
	meta     []int
	idmap    map[int]*node
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	var n *node
	for scanner.Scan() {
		var txt = scanner.Text()
		n, err = parseTree(txt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println(sumMeta(n))
	fmt.Println(nodeValue(n))

	fmt.Println("Time", time.Since(startTime))
}

func sumMeta(n *node) int {
	sum := 0
	for _, v := range n.meta {
		sum += v
	}

	for _, c := range n.children {
		sum += sumMeta(c)
	}

	return sum
}

func nodeValue(n *node) int {
	var val int
	if len(n.children) == 0 {
		val = sumMeta(n)
	} else {
		for _, m := range n.meta {
			if c, ok := n.idmap[m]; ok {
				val += nodeValue(c)
			}
		}
	}
	return val
}

func parseTree(txt string) (*node, error) {
	flds := strings.Fields(txt)
	var d []int
	for _, fld := range flds {
		i, err := strconv.Atoi(fld)
		if err != nil {
			return nil, err
		}
		d = append(d, i)
	}

	return parseNode(nil, d, 0), nil
}

func parseNode(parent *node, ints []int, id int) *node {
	cc := ints[0]
	mc := ints[1]

	n := &node{id: id, ccount: cc, mcount: mc, parent: parent, idmap: make(map[int]*node)}

	idx := 0
	for i := 0; i < cc; i++ {
		id := i + 1
		sub := ints[idx+2 : len(ints)-mc]
		c := parseNode(n, sub, id)
		n.children = append(n.children, c)
		idx += nodeLen(c)
		n.idmap[id] = c
	}

	end := nodeLen(n)
	n.meta = ints[end-mc : end]

	return n
}

func nodeLen(n *node) int {
	l := 2
	for _, c := range n.children {
		l += nodeLen(c)
	}
	return l + n.mcount
}
