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

type instr struct {
	op    string
	param int
}

type xy struct {
	x, y int
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
	fmt.Fprintln(w, "--2022 day 10 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "-----Part 2-----")
	fmt.Fprintln(w, p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	instrs := parseInput(in)
	sum, _ := run(instrs, 20, 60, 100, 140, 180, 220)
	return sum
}

func part2(in input) string {
	instrs := parseInput(in)
	_, strs := run(instrs, 240)
	str := ""
	for _, s := range strs {
		str += fmt.Sprintf("%s\n", s)
	}
	return str
}

func run(instrs []instr, checkcycles ...int) (int, []string) {
	X := 1
	xqueue := []int{}
	cycle := 0
	crt := make(map[xy]bool)

	observe := make(map[int]int)
	for _, c := range checkcycles {
		observe[c] = 0
	}

	i := 0
	working := 0
	for i < len(instrs) {
		cycle++
		if _, ok := observe[cycle]; ok {
			observe[cycle] = cycle * X
		}

		if working == 0 {
			ins := instrs[i]
			switch ins.op {
			case "noop":
				break
			case "addx":
				xqueue = append(xqueue, []int{0, ins.param}...)
				working = 1
			}
			i++
		} else {
			working--
		}

		if len(xqueue) > 0 {
			X += xqueue[0]
			xqueue = xqueue[1:]
		}

		col := cycle % 40
		pt := xy{x: col, y: cycle / 40}
		if _, ok := crt[pt]; !ok {
			crt[pt] = (col >= X-1 && col <= X+1)
		}
	}

	s := 0
	for _, v := range observe {
		s += v
	}

	list := []string{}
	for y := 0; y < 6; y++ {
		s := ""
		for x := 0; x < 40; x++ {
			if b, ok := crt[xy{x, y}]; ok && b {
				s += "#"
			} else {
				s += "."
			}
		}
		list = append(list, s)
	}

	return s, list
}

func parseInput(in input) []instr {
	instrs := []instr{}
	for _, line := range in {
		sp := strings.Split(line, " ")

		if len(sp) > 0 {
			s := instr{op: sp[0]}
			if len(sp) == 2 {
				v, _ := strconv.Atoi(sp[1])
				s.param = v
			}
			instrs = append(instrs, s)
		}
	}
	return instrs
}
