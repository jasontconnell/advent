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
	id      int
	next    *cup
	removed bool
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

	result := play(first, 100)
	for result.next.id != 1 {
		result = result.next
	}

	fmt.Println("Part 1:", strings.Replace(prints(result), "1", "", 1))

	fmt.Println("Time", time.Since(startTime))
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

func play(c *cup, moves int) *cup {
	mv := 0
	cur := c
	for mv < moves {
		v := cur.id
		start := cur
		remove := cur

		var after *cup
		var firstremoved *cup
		var lastremoved *cup
		for j := 0; j < 3; j++ {
			remove = remove.next
			if firstremoved == nil {
				firstremoved = remove // keep track of removed
			}

			if j == 2 {
				after = remove.next
				lastremoved = remove
				lastremoved.next = nil
			}
		}

		start.next = after

		v = v - 1
		if v == 0 {
			v = 9
		}

		ptr := cur
		for {
			ptr = ptr.next
			if ptr.id == v {
				break
			}

			if ptr == cur {
				v = v - 1
				if v == 0 {
					v = 9
				}
			}
		}

		cur = ptr
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
