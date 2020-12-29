package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "23.txt"

type cup struct {
	id   int
	next *cup
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

	cupvals := getCups(lines[0])
	first := fillCups(cupvals)

	result := play(first, 100, 9)
	for result.next.id != 1 {
		result = result.next
	}

	fmt.Println("Part 1:", strings.Replace(prints(result), "1", "", 1))

	max := 0
	for _, cv := range cupvals {
		if cv > max {
			max = cv
		}
	}

	cur := max + 1
	for len(cupvals) < 1_000_000 {
		cupvals = append(cupvals, cur)
		cur++
	}

	p2first := fillCups(cupvals)
	p2result := play(p2first, 10_000_000, cur-1)
	p2cur := p2result
	for p2cur.id != 1 {
		p2cur = p2cur.next
	}
	p2 := p2cur.next.id * p2cur.next.next.id

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func fillCups(cupvals []int) *cup {
	var first *cup
	var cur *cup
	for _, cv := range cupvals {
		c := &cup{id: cv}

		if first == nil {
			first = c
		}

		if cur != nil {
			cur.next = c
		}
		cur = c
	}

	cur.next = first

	return first
}

func prints(c *cup) string {
	s := ""
	p := c.next
	done := false
	for !done {
		s += fmt.Sprint(p.id)
		p = p.next

		done = p == c.next
	}
	return s
}

func play(c *cup, moves, max int) *cup {
	mv := 0
	cur := c
	m := make(map[int]*cup)

	for cur.next != c {
		m[cur.id] = cur
		cur = cur.next
	}
	m[cur.id] = cur

	cur = c

	for mv < moves {
		// if mv%1000 == 0 {
		// 	fmt.Println(mv)
		// }
		v := cur.id
		start := cur

		var firstremoved *cup = cur.next
		var lastremoved *cup = cur.next.next.next

		start.next = lastremoved.next
		lastremoved.next = nil

		v = v - 1
		if v == 0 {
			v = max
		}

		ptr := m[v]
		for ptr == nil || ptr.next == nil || ptr.next.next == nil || ptr.next.next.next == nil {
			v = v - 1
			if v == 0 {
				v = max
			}
			ptr = m[v]
		}

		curnext := ptr.next
		ptr.next = firstremoved
		lastremoved.next = curnext

		cur = start.next

		mv++
	}

	return cur
}

func getCups(line string) []int {
	cups := []int{}
	for _, ch := range line {
		v, _ := strconv.Atoi(string(ch))
		cups = append(cups, v)
	}
	return cups
}
