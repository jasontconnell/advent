package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

type rockpath struct {
	list []xy
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	paths := parseInput(in)
	for _, p := range paths {
		fmt.Println(p)
	}
	return 0
}

func part2(in input) output {
	return 0
}

func parseInput(in input) []rockpath {
	p := []rockpath{}
	for _, line := range in {
		coords := strings.Split(line, " -> ")
		list := []xy{}
		for _, coord := range coords {
			ns := strings.Split(coord, ",")
			x, _ := strconv.Atoi(ns[0])
			y, _ := strconv.Atoi(ns[1])

			pt := xy{x, y}
			list = append(list, pt)
		}
		path := rockpath{list: list}
		p = append(p, path)
	}
	return p
}
