package main

import (
	"fmt"
	"time"
)

var input = 343

type State struct {
	Pos int
	Val int
}

func main() {
	startTime := time.Now()

	val := spinlockv3(input, 0, 2017, 2017)
	val2 := spinlockv3(input, 0, 50000000, 0)

	fmt.Println("value after 2017            :", val)
	fmt.Println("value after 0               :", val2)

	fmt.Println("Time", time.Since(startTime))
}

type Node struct {
	Left  *Node
	Right *Node
	Value int
}

func spinlockv3(skip, start, nums, find int) int {
	s := &Node{Value: start}
	s.Left = s
	s.Right = s

	ln := 1
	cur := s

	for i := start + 1; i <= nums; i++ {
		for j := 0; j < skip%ln; j++ {
			cur = cur.Right
		}

		n := &Node{Value: i, Left: cur, Right: cur.Right}
		tmp := cur.Right

		tmp.Left = n
		cur.Right = n

		cur = n
		ln++
	}

	loops := 0
	for cur.Value != find && loops <= nums {
		loops++
		cur = cur.Right
	}

	return cur.Right.Value
}
