package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xyz struct {
	x, y, z int
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
	fmt.Fprintln(w, "--2022 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	pts := parseInput(in)
	return getNotConnectedArea(pts)
}

func part2(in input) output {
	return 0
}

func getNotConnectedArea(pts []xyz) int {
	area := 0
	for _, pt := range pts {
		area += getExposed(pts, pt)
	}
	return area
}

func getExposed(pts []xyz, check xyz) int {
	exposed := 6
	for _, pt := range pts {
		if pt == check {
			continue
		}
		if dist(pt, check) == 1 {
			exposed--
		}
	}
	return exposed
}

func dist(pt1, pt2 xyz) int {
	d := abs(pt1.x-pt2.x) + abs(pt1.y-pt2.y) + abs(pt1.z-pt2.z)
	return d
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func parseInput(in input) []xyz {
	pts := []xyz{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		if len(sp) != 3 {
			continue
		}
		var x, y, z int

		x, _ = strconv.Atoi(sp[0])
		y, _ = strconv.Atoi(sp[1])
		z, _ = strconv.Atoi(sp[2])

		pts = append(pts, xyz{x, y, z})
	}
	return pts
}
