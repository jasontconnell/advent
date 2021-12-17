package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

const (
	safe byte = '.'
	trap byte = '^'
)

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return countSafe(in, 40)
}

func part2(in input) output {
	return countSafe(in, 400000)
}

func countSafe(in input, rows int) int {
	count := 0
	prev := in
	for i := 0; i < rows; i++ {
		count += strings.Count(prev, string(safe))
		prev = transform(prev)
	}
	return count
}

func transform(str string) string {
	ns := ""
	for i := 0; i < len(str); i++ {
		nc := '^'
		if isSafe(i, str) {
			nc = '.'
		}
		ns += string(nc)
	}
	return ns
}

func isSafe(index int, prev string) bool {
	tri := getThreePrev(index, prev)
	//fmt.Println(prev, index, tri)
	istrap := tri[0] == trap && tri[1] == trap && tri[2] != trap ||
		tri[0] != trap && tri[1] == trap && tri[2] == trap ||
		tri[0] == trap && tri[1] != trap && tri[2] != trap ||
		tri[0] != trap && tri[1] != trap && tri[2] == trap

	return !istrap
}

func getThreePrev(index int, prev string) string {
	tri := ""
	if index == 0 {
		tri = "." + string(prev[:2])
	} else if index == len(prev)-1 {
		tri = string(prev[len(prev)-2:]) + "."
	} else {
		tri = string(prev[index-1 : index+2])
	}

	return tri
}
