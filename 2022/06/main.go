package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 06 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return findNonRepeat(in, 4)
}

func part2(in input) output {
	return findNonRepeat(in, 14)
}

func findNonRepeat(in input, n int) int {
	idx := 0
	for i := n; i < len(in); i++ {
		if !hasRepeats(in[i-n : i]) {
			idx = i
			break
		}
	}
	return idx
}

func hasRepeats(rs string) bool {
	rep := false
	for i, r := range rs {
		for j, r2 := range rs {
			if i == j {
				continue
			}
			if r == r2 {
				rep = true
				break
			}
		}
	}
	return rep
}
