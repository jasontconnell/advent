package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = int
type output = int

type State struct {
	Pos int
	Val int
}

type Node struct {
	Left  *Node
	Right *Node
	Value int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadInt(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	val := spinlockv3(in, 2017, 2017)
	return val
}

func part2(in input) output {
	val := spinlockv3fast(in, 50_000_000, 1)
	return val
}

func spinlockv3fast(skip, nums, pos int) int {
	length := 1

	n := 1
	inpos := 0
	idx := 0
	for n < nums+1 {
		idx := (idx + skip) % length
		if idx == pos-1 {
			inpos = n
		}
		length++
		idx = (idx + 1) % length
		n++
	}
	return inpos
}

func spinlockv3(skip, nums, find int) int {
	list := []int{0}
	n := 1
	idx := 0

	for n < nums+1 {
		idx = (idx + skip) % len(list)
		list = append(list[:idx+1], append([]int{n}, list[idx+1:]...)...)
		idx = (idx + 1) % len(list)
		n++
	}

	val := -1
	for i := 0; i < len(list); i++ {
		if list[i] == find {
			val = list[(i+1)%len(list)]
			break
		}
	}

	return val
}
