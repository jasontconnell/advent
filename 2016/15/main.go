package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Disc struct {
	Id            int
	StartPosition int
	Position      int
	Positions     int
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
	fmt.Fprintln(w, "--2016 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	discs := parseInput(in)
	return run(discs)
}

func part2(in input) output {
	discs := parseInput(in)
	discs = append(discs, Disc{Id: 7, Positions: 11, Position: 0})
	return run(discs)
}

func run(discs []Disc) int {
	tick := 0
	for {
		for i := 0; i < len(discs); i++ {
			discTime := discs[i].StartPosition + i + 1 + tick
			discs[i].Position = discTime
			discs[i].Position = discs[i].Position % discs[i].Positions
		}

		allzero := true
		for _, d := range discs {
			allzero = allzero && d.Position == 0
		}
		if allzero {
			break
		}
		tick++
	}
	return tick
}

func parseInput(in input) []Disc {
	reg := regexp.MustCompile("^Disc #([0-9]+) has ([0-9]+) positions?; at time=0, it is at position ([0-9]+)\\.")

	discs := []Disc{}
	for _, line := range in {
		groups := reg.FindStringSubmatch(line)
		if len(groups) == 4 {
			id, _ := strconv.Atoi(groups[1])
			poss, _ := strconv.Atoi(groups[2])
			pos, _ := strconv.Atoi(groups[3])

			d := Disc{Id: id, Positions: poss, StartPosition: pos, Position: pos}

			discs = append(discs, d)
		}
	}
	return discs
}
