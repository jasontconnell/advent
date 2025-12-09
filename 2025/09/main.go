package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 09 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	points := parseInput(in)
	return largestArea(points)
}

func part2(in input) output {
	return 0
}

func largestArea(pts []xy) int {
	max := 0
	for i := 0; i < len(pts)-1; i++ {
		for j := i + 1; j < len(pts); j++ {
			p1, p2 := pts[i], pts[j]
			area := getArea(p1.x, p2.x, p1.y, p2.y)
			if area > max {
				max = area
			}
		}
	}
	return max
}

func getArea(x1, x2, y1, y2 int) int {
	return int((math.Abs(float64(x2-x1)) + 1) * (math.Abs(float64(y2-y1)) + 1))
}

func parseInput(in input) []xy {
	list := []xy{}
	for _, line := range in {
		sp := strings.Split(line, ",")
		x, _ := strconv.Atoi(sp[0])
		y, _ := strconv.Atoi(sp[1])
		list = append(list, xy{x, y})
	}
	return list
}
