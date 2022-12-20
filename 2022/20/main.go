package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []int
type output = int

type node struct {
	value       int
	moved       bool
	left, right *node
	origidx     int
}

func (n *node) String() string {
	return fmt.Sprintf("[%d] = %d moved: %t", n.origidx, n.value, n.moved)
}

func print(front *node, ptr *node) {
	cur := front
	first := true
	for cur != front || first {
		first = false
		f := "%d"
		if cur == ptr {
			f = "[%d]"
		}
		fmt.Print(fmt.Sprintf(f, cur.value))
		if cur.right != front {
			fmt.Print(" -> ")
		}
		cur = cur.right
	}
	fmt.Println()
}

func main() {
	startTime := time.Now()

	in, err := common.ReadInts(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	front := getList(in)
	mix(front, len(in))
	return getGroveCoordinates(front)
}

func part2(in input) output {
	front := getList(in)
	applyDecryptKey(front, 811589153)
	for i := 0; i < 10; i++ {
		mix(front, len(in))
	}
	return getGroveCoordinates(front)
}

func applyDecryptKey(front *node, k int) {
	cur := front
	first := true
	for cur != front || first {
		first = false
		cur.value *= k
		cur = cur.right
	}
}

func mix(front *node, count int) {
	cur := front
	ptr := front
	idx := 0
	for idx < count {
		tomove := cur.value
		dir := 1
		if tomove < 0 {
			dir = -1
			tomove *= -1
		}

		tomove = tomove % (count - 1)
		// this also works.
		// for tomove/count != 0 {
		// 	tomove = tomove/count + tomove%count
		// }

		ptr = cur.right
		cur.moved = true

		for j := 0; j < tomove; j++ {
			if dir == 1 {
				tmpleft := cur.left
				tmpright := cur.right
				tmpright.left = tmpleft
				tmpleft.right = tmpright

				cur.right = tmpright.right
				cur.left = tmpright
				tmpright.right = cur
				cur.right.left = cur

			} else {
				tmpleft := cur.left
				tmpright := cur.right
				tmpright.left = tmpleft
				tmpleft.right = tmpright

				cur.left = tmpleft.left
				cur.right = tmpleft
				tmpleft.left = cur
				cur.left.right = cur
			}
		}

		idx++
		for ptr.origidx != idx && idx < count {
			ptr = ptr.right
		}
		cur = ptr
	}
}

func getGroveCoordinates(front *node) int {
	gc := []int{}

	cur := front
	for cur.value != 0 {
		cur = cur.right
	}

	for len(gc) < 3 {
		for i := 0; i < 1000; i++ {
			cur = cur.right
		}
		gc = append(gc, cur.value)
	}
	return gc[0] + gc[1] + gc[2]
}

func getList(in input) *node {
	var front, prev *node
	for i := 0; i < len(in); i++ {
		if i == 0 {
			front = &node{value: in[i], origidx: i}
			prev = front
			continue
		}

		n := &node{value: in[i], left: prev, origidx: i}
		prev.right = n
		prev = n
	}
	front.left = prev
	prev.right = front
	return front
}
