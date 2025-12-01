package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type instr struct {
	quant int
	dir   int
}

type dial struct {
	prev, next *dial
	val        int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 01 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	instrs := parseInput(in)
	safe := makeSafe(100, 50)

	z := 0
	turn(instrs, safe, func(val int) {
		if val == 0 {
			z++
		}
	}, nil)

	return z
}

func part2(in input) output {
	instrs := parseInput(in)
	safe := makeSafe(100, 50)

	z := 0
	turn(instrs, safe, nil, func(val int) {
		if val == 0 {
			z++
		}
	})

	return z
}

func turn(instrs []instr, safe *dial, land func(val int), turn func(val int)) {
	cur := safe
	for _, i := range instrs {

		next := i.dir == 1

		for j := 0; j < i.quant; j++ {
			if next {
				cur = cur.next
			} else {
				cur = cur.prev
			}

			if turn != nil {
				turn(cur.val)
			}
		}

		if land != nil {
			land(cur.val)
		}
	}
}

func parseInput(lines input) []instr {
	instrs := []instr{}
	for _, line := range lines {
		dir := 1
		if line[0] == 'L' {
			dir = -1
		}

		q, err := strconv.Atoi(line[1:])
		if err != nil {
			log.Println("parsing", line)
		}

		instrs = append(instrs, instr{quant: q, dir: dir})
	}
	return instrs
}

func makeSafe(n, start int) *dial {
	list := []*dial{}

	for i := 0; i < n; i++ {
		list = append(list, &dial{val: i})
	}

	for i, v := range list {
		var prev, next *dial
		if i == 0 {
			prev = list[len(list)-1]
		} else {
			prev = list[i-1]
		}
		if i == len(list)-1 {
			next = list[0]
		} else {
			next = list[i+1]
		}

		v.prev = prev
		v.next = next
	}

	return list[start]
}
