package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/advent/2019/intcode"
	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []int
type output = int

var prog string = `NOT A J
NOT C T
OR T J
AND D J
WALK
`

var prog2 string = `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
AND E T
OR H T
AND T J
RUN
`

func main() {
	startTime := time.Now()

	in, err := common.ReadIntCsv(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	hullDamage := 0
	c := intcode.NewComputer(in)
	c.OnOutput = func(i int) {
		if i < 256 {
			fmt.Print(string(rune(i)))
		} else {
			hullDamage = i
		}
	}
	c.AddInput(getInts(prog)...)
	c.Exec()
	return hullDamage
}

func part2(in input) output {
	hullDamage := 0
	c := intcode.NewComputer(in)
	c.OnOutput = func(i int) {
		if i < 256 {
			fmt.Print(string(rune(i)))
		} else {
			hullDamage = i
		}
	}
	c.AddInput(getInts(prog2)...)
	c.Exec()
	return hullDamage
}

func getInts(s string) []int {
	d := []int{}
	for _, c := range s {
		d = append(d, int(c))
	}
	fmt.Println(d)
	return d
}
