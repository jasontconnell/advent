package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type stone struct {
	digits string
	next   *stone
	prev   *stone
}

func main() {
	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 11 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	start := parse(in)
	return count(blink(start, 25))
}

func part2(in input) output {
	return 0
}

func count(start *stone) int {
	c := 0
	cur := start
	for cur != nil {
		c++
		cur = cur.next
	}
	return c
}

func blink(start *stone, n int) *stone {
	for i := 0; i < n; i++ {
		start = change(start)
	}
	return start
}

func change(start *stone) *stone {
	cur := start
	for cur != nil {
		x, _ := strconv.Atoi(cur.digits)
		if cur.digits == "0" {
			cur.digits = "1"
		} else if len(cur.digits)%2 == 0 {
			left, right := &stone{}, &stone{}
			ldigits := cur.digits[:len(cur.digits)/2]
			rdigits := cur.digits[len(cur.digits)/2:]

			ln, _ := strconv.Atoi(ldigits)
			rn, _ := strconv.Atoi(rdigits)

			left.digits = strconv.Itoa(ln)
			right.digits = strconv.Itoa(rn)

			repl := cur

			if cur.prev == nil {
				start = left
			}

			left.prev = repl.prev
			left.next = right
			right.prev = left
			right.next = repl.next
			if repl.prev != nil {
				repl.prev.next = left
			}
			if repl.next != nil {
				cur.next.prev = right
			}
			repl.next = nil
			repl.prev = nil
			cur = right
		} else {
			n := x * 2024
			cur.digits = strconv.Itoa(n)
		}
		cur = cur.next
	}
	return start
}

func print(msg string, start *stone) {
	cur := start
	fmt.Print(msg + "   ")
	for cur != nil {
		fmt.Print(cur.digits, " ")
		cur = cur.next
	}
	fmt.Println()
}

func parse(in string) *stone {
	sp := strings.Fields(in)
	var last *stone
	var first *stone
	for _, s := range sp {
		st := &stone{digits: s, prev: last}
		if last != nil {
			last.next = st
		}
		if first == nil {
			first = st
		}
		last = st
	}
	return first
}
