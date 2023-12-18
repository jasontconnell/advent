package main

import (
	"fmt"
	"log"
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

func (pt xy) mult(x int) xy {
	return xy{pt.x * x, pt.y * x}
}

func (pt xy) adddir(d dir) xy {
	return pt.add(xy(d))
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
	return getArea(instrs, false)
}

func part2(in input) output {
	instrs := parseInput(in)
	return getArea(instrs, true)
}

func hexvalue(instr diginstr) (dir, int) {
	dr := rune(instr.color[len(instr.color)-1])
	hex := instr.color[1 : len(instr.color)-1]
	val, _ := strconv.ParseInt(hex, 16, 32)

	var d dir
	switch dr {
	case '0':
		d = east
	case '1':
		d = south
	case '2':
		d = west
	case '3':
		d = north
	}

	return d, int(val)
}

func getArea(instrs []diginstr, hexmode bool) int {
	area := 0
	totalholes := 0
	var cur, next xy

	for i := 0; i < len(instrs); i++ {
		var instrdir dir
		var val int
		var instr diginstr

		instr = instrs[i]
		instrdir, val = instr.dir, instr.dist
		if hexmode {
			instrdir, val = hexvalue(instr)
		}
		totalholes += val

		next = cur.add(xy(instrdir).mult(val))
		area += cur.x*next.y - cur.y*next.x + val
		cur = next
	}

	return area/2 + 1
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
		color := strings.TrimRight(strings.TrimLeft(sp[2], "("), ")")
		instr := diginstr{dir: d, dist: dist, color: color}
		instrs = append(instrs, instr)
	}
	return instrs
}
