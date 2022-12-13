package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/2017/knot"
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
	fmt.Fprintln(w, "--2017 day 10 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	lns := parseInput(in)
	res := knot.PrimitiveKnotHash(lns)
	return res[0] * res[1]
}

func part2(in input) string {
	return knot.KnotHash(in)
}

func parseInput(in input) []int {
	n := []int{}
	sp := strings.Split(in, ",")
	for _, s := range sp {
		i, _ := strconv.Atoi(s)
		n = append(n, i)
	}
	return n
}
