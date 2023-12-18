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

func (pt xy) add(p2 xy) xy {
	return xy{pt.x + p2.x, pt.y + p2.y}
}

func (pt xy) adddir(p2 dir) xy {
	return pt.add(xy(p2))
}

type diginstr struct {
	dir   dir
	dist  int
	color string
}

type block struct {
	color string
}

type dir xy

var (
	north dir = dir{0, -1}
	south dir = dir{0, 1}
	east  dir = dir{1, 0}
	west  dir = dir{-1, 0}
)

func print(m map[xy]block) {
	min, max := minmax(m)
	fmt.Println(len(m), min, max)
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			pt := xy{x, y}
			if _, ok := m[pt]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func minmax(m map[xy]block) (xy, xy) {
	min := xy{math.MaxInt32, math.MaxInt32}
	max := xy{math.MinInt32, math.MinInt32}

	for k := range m {
		if k.x < min.x {
			min.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y > max.y {
			max.y = k.y
		}
	}

	return min, max
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	instrs := parseInput(in)
	m := dig(instrs)
	return getEnclosedArea(m)
}

func part2(in input) output {
	return 0
}

func dig(instrs []diginstr) map[xy]block {
	m := make(map[xy]block)
	cur := xy{0, 0}
	for _, instr := range instrs {
		for i := 0; i < instr.dist; i++ {
			b := block{color: instr.color}
			m[cur] = b
			cur = cur.adddir(instr.dir)
		}
	}
	return m
}

func getEnclosedArea(m map[xy]block) int {
	loop := getLoop(m)
	area := 0
	for i := 0; i < len(loop); i++ {
		cur := loop[i]
		next := loop[(i+1)%len(loop)]
		area += cur.y*cur.x - cur.x*next.y
	}
	return area - len(loop)/2 + 1 + len(loop)
}

func getLoop(m map[xy]block) []xy {
	var start xy
	startfound := false
	min, max := minmax(m)
	for y := min.y; y <= max.y && !startfound; y++ {
		for x := min.x; x <= max.x; x++ {
			pt := xy{x, y}
			if _, ok := m[pt]; ok {
				start = pt
				startfound = true
				break
			}
		}
	}

	loop := []xy{start}
	cur := start
	v := make(map[xy]bool)
	v[start] = true
	for i := 0; i < len(m); i++ {
		for _, d := range []dir{north, south, east, west} {
			np := cur.adddir(d)
			_, vok := v[np]
			if _, ok := m[np]; ok && !vok {
				loop = append(loop, np)
				cur = np
				v[np] = true
				break
			}
		}
	}
	return loop
}

func parseInput(in input) []diginstr {
	instrs := []diginstr{}
	for _, line := range in {
		sp := strings.Fields(line)
		var d dir
		switch sp[0] {
		case "U":
			d = north
		case "D":
			d = south
		case "R":
			d = east
		case "L":
			d = west
		}
		dist, _ := strconv.Atoi(sp[1])
		instr := diginstr{dir: d, dist: dist, color: sp[2]}
		instrs = append(instrs, instr)
	}
	return instrs
}
