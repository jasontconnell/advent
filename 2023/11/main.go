package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

type space struct {
	pt xy
	ch rune
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 11 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parseInput(in)
	exx, exy := getExpandingSpaces(m)
	return getDistances(m, exx, exy, 2)
}

func part2(in input) output {
	m := parseInput(in)
	exx, exy := getExpandingSpaces(m)
	return getDistances(m, exx, exy, 1000000)
}

func getDistances(m map[xy]space, exx, exy []int, balloon int) int {
	list := []xy{}
	for k := range m {
		list = append(list, k)
	}
	total := 0
	for len(list) > 0 {
		for i := 1; i < len(list); i++ {
			total += getDistance(list[0], list[i], exx, exy, balloon)
		}
		list = list[1:]
	}
	return total
}

func getDistance(p1, p2 xy, exx, exy []int, balloon int) int {
	d := int(math.Abs(float64(p2.x-p1.x)) + math.Abs(float64(p2.y-p1.y)))
	for _, x := range exx {
		if (x > p1.x && x < p2.x) || (x > p2.x && x < p1.x) {
			d += balloon - 1
		}
	}
	for _, y := range exy {
		if (y > p1.y && y < p2.y) || (y > p2.y && y < p1.y) {
			d += balloon - 1
		}
	}
	return d
}

func getExpandingSpaces(m map[xy]space) ([]int, []int) {
	maxx, maxy := getMaxes(m)
	xs, ys := []int{}, []int{}

	for y := 0; y <= maxy; y++ {
		found := false
		for x := 0; x <= maxx && !found; x++ {
			cp := xy{x, y}
			if _, ok := m[cp]; ok {
				found = true
			}
		}
		if !found {
			ys = append(ys, y)
		}
	}
	for x := 0; x <= maxx; x++ {
		found := false
		for y := 0; y <= maxy && !found; y++ {
			cp := xy{x, y}
			if _, ok := m[cp]; ok {
				found = true
			}
		}
		if !found {
			xs = append(xs, x)
		}
	}
	return xs, ys
}

func getMaxes(m map[xy]space) (int, int) {
	maxx, maxy := 0, 0

	for k := range m {
		if k.x > maxx {
			maxx = k.x
		}
		if k.y > maxy {
			maxy = k.y
		}
	}
	return maxx, maxy
}

func parseInput(in input) map[xy]space {
	m := make(map[xy]space)
	for y, line := range in {
		for x, c := range line {
			if c == '.' {
				continue
			}
			pt := xy{x, y}
			sp := space{pt: pt, ch: c}
			m[pt] = sp
		}
	}
	return m
}
