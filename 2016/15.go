package main

import (
	"fmt"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

type Disc struct {
	Position  int
	Positions int
}

func main() {
	startTime := time.Now()
	part2 := true

	discs := []Disc{}
	discs = append(discs, Disc{Positions: 13, Position: 1})
	discs = append(discs, Disc{Positions: 19, Position: 10})
	discs = append(discs, Disc{Positions: 3, Position: 2})
	discs = append(discs, Disc{Positions: 7, Position: 1})
	discs = append(discs, Disc{Positions: 5, Position: 3})
	discs = append(discs, Disc{Positions: 17, Position: 5})
	if part2 {
		discs = append(discs, Disc{Positions: 11, Position: 0})
	}

	solved := false
	tick := 0

	for !solved {
		cp := copyDiscs(discs)
		if test(cp, tick) {
			fmt.Println("Got a capsule at time", tick)
			solved = true
		}
		tick++
	}

	fmt.Println("Time", time.Since(startTime))
}

func copyDiscs(discs []Disc) []Disc {
	nd := make([]Disc, len(discs))
	copy(nd, discs)
	return nd
}

func test(discs []Disc, tick int) bool {
	for i := 0; i < len(discs); i++ {
		discTime := i + 1 + tick
		discs[i].Position += discTime
		discs[i].Position = discs[i].Position % discs[i].Positions
	}

	allzero := true
	for _, d := range discs {
		allzero = allzero && d.Position == 0
	}
	return allzero
}
