package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = int
type output = int

type Elf struct {
	Number   int
	Presents int
	Next     *Elf
	Prev     *Elf
}

func main() {
	startTime := time.Now()

	in, err := common.ReadInt(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	elf := getElves(in)
	return solve(elf, in, false)
}

func part2(in input) output {
	elf := getElves(in)
	return solve(elf, in, true)
}

func solve(elf *Elf, numElves int, solveAcross bool) int {
	solved := false
	var across *Elf
	if solveAcross {
		across = elf
		for i := 0; i < numElves/2; i++ {
			across = across.Next
		}
	}
	count := numElves

	for !solved {
		if elf.Presents > 0 {
			if solveAcross {
				elf.Presents += across.Presents
				across.Next.Prev = across.Prev
				across.Prev.Next = across.Next
			} else {
				elf.Presents += elf.Next.Presents
				elf.Next = elf.Next.Next
				elf.Next.Prev = elf
			}
		}
		count--
		elf = elf.Next
		if solveAcross {
			across = across.Next
			if count%2 == 0 {
				across = across.Next
			}
		}
		solved = elf.Next == elf
	}

	return elf.Number
}

func getElves(in input) *Elf {
	elf := &Elf{Number: 1, Presents: 1, Next: nil, Prev: nil}
	first := elf
	for i := 1; i < in; i++ {
		next := Elf{Number: i + 1, Presents: 1, Next: nil, Prev: elf}
		elf.Next = &next
		elf = &next
		if i == in-1 {
			next.Next = first
			first.Prev = &next
		}
	}
	return first
}
