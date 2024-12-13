package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}
type clawmachine struct {
	a, b  button
	prize xy
}

type button struct {
	xdelta, ydelta int
}

type state struct {
	presses int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 13 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	machines := parse(in)
	return solveAll(machines)
}

func part2(in input) output {
	machines := parse(in)
	machines = modifyMachines(machines)
	return solveAll(machines)
}

func solveAll(machines []clawmachine) int {
	total := 0
	for _, mach := range machines {
		if tokens, ok := solveOne(mach); ok {
			total += tokens
		}
	}
	return total
}

func modifyMachines(machines []clawmachine) []clawmachine {
	for i := 0; i < len(machines); i++ {
		machines[i].prize.x += 10000000000000
		machines[i].prize.y += 10000000000000
	}
	return machines
}

func solveOne(mach clawmachine) (int, bool) {
	d := mach.a.xdelta*mach.b.ydelta - mach.b.xdelta*mach.a.ydelta
	d1 := mach.prize.x*mach.b.ydelta - mach.prize.y*mach.b.xdelta
	d2 := mach.prize.y*mach.a.xdelta - mach.prize.x*mach.a.ydelta

	if d1%d != 0 || d2%d != 0 {
		return 0, false
	}
	return (d1/d)*3 + d2/d, true
}

func parse(in []string) []clawmachine {
	breg := regexp.MustCompile(`^Button ([AB]): X\+([0-9]+), Y\+([0-9]+)$`)
	preg := regexp.MustCompile(`^Prize: X=([0-9]+), Y=([0-9]+)$`)
	machines := []clawmachine{}
	mach := clawmachine{}
	for _, line := range in {
		if breg.MatchString(line) {
			m := breg.FindStringSubmatch(line)
			x, _ := strconv.Atoi(m[2])
			y, _ := strconv.Atoi(m[3])
			a := m[1] == "A"
			btn := button{xdelta: x, ydelta: y}
			if a {
				mach.a = btn
			} else {
				mach.b = btn
			}
		}
		if preg.MatchString(line) {
			m := preg.FindStringSubmatch(line)
			x, _ := strconv.Atoi(m[1])
			y, _ := strconv.Atoi(m[2])
			mach.prize = xy{x, y}
			machines = append(machines, mach)
			mach = clawmachine{}
		}
	}
	return machines
}
