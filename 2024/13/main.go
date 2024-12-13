package main

import (
	"fmt"
	"log"
	"math"
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
	return solveAll(modifyMachines(machines))
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
	return machines
	for i := 0; i < len(machines); i++ {
		machines[i].prize.x += 10000000000000
		machines[i].prize.y += 10000000000000
	}
	return machines
}

func solveOne(mach clawmachine) (int, bool) {

	ab := mach.prize.x / mach.a.xdelta
	bb := mach.prize.x / mach.b.xdelta
	start := common.Min(ab, bb)

	left := start == ab
	queue := []state{{presses: start}}
	visit := make(map[state]bool)
	min := math.MaxInt32
	solved := false
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.presses <= 0 {
			continue
		}

		if _, ok := visit[cur]; ok {
			continue
		}
		visit[cur] = true

		var abutton, bbutton button

		if left {
			abutton = mach.a
			bbutton = mach.b
		} else {
			abutton = mach.b
			bbutton = mach.a
		}

		avx := abutton.xdelta * cur.presses
		avy := abutton.ydelta * cur.presses

		if (mach.prize.x-avx)%bbutton.xdelta == 0 && (mach.prize.y-avy)%bbutton.ydelta == 0 {
			// even divide
			rem := (mach.prize.x - avx) / bbutton.xdelta
			if avy+rem*bbutton.ydelta == mach.prize.y {
				solved = true

				apress := cur.presses
				bpress := rem

				if !left {
					apress = rem
					bpress = cur.presses
				}

				total := apress*3 + bpress
				if total < min {
					min = total
				}
				continue
			}
		}

		queue = append(queue, state{presses: cur.presses - 1})
	}
	return min, solved
}

func parse(in []string) []clawmachine {
	breg := regexp.MustCompile(`^Button ([AB]): X\+([0-9]+), Y\+([0-9]+)$`)
	preg := regexp.MustCompile(`^Prize: X=([0-9]+), Y=([0-9]+)$`)
	machines := []clawmachine{}
	mach := clawmachine{}
	for i, line := range in {
		if line == "" || i == len(in)-1 {
			machines = append(machines, mach)
			continue
		}

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
		}
	}
	return machines
}
