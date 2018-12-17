package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	"strings"
	//"math"
)

var input = "08.txt"

type node struct {
	id       int
	ccount   int
	mcount   int
	children []*node
	parent   *node
	meta     []int
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

	n := &node{id: id, ccount: cc, mcount: mc, parent: parent}

	idx := 0
	for i := 0; i < cc; i++ {
		sub := ints[idx+2 : len(ints)-mc]
		c := parseNode(n, sub, id+1+i)
		n.children = append(n.children, c)
		idx += nodeLen(c)
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
