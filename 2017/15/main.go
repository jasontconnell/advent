package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Generator struct {
	Previous int
	Factor   int
}

type Judge struct {
	Div  int
	GenA []int
	GenB []int
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
	fmt.Fprintln(w, "--2017 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	gens := parseInput(in)
	var judge Judge = Judge{Div: 2147483647}
	matched := compute(gens[0], gens[1], judge, 40000000)
	return matched
}

func part2(in input) output {
	gens := parseInput(in)
	var judge Judge = Judge{Div: 2147483647}
	matched2 := computePart2(gens[0], gens[1], judge, 5000000)
	return matched2
}

func compute(a, b Generator, j Judge, itrs int) int {
	matched := 0
	for i := 0; i < itrs; i++ {
		processGen(&a, j)
		processGen(&b, j)

		ra := a.Previous & 65535
		rb := b.Previous & 65535

		if ra == rb {
			matched++
		}
	}

	return matched
}

func computePart2(a, b Generator, j Judge, c int) int {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		parallelCompute(&a, j, &j.GenA, c, 4)
		wg.Done()
	}()

	go func() {
		parallelCompute(&b, j, &j.GenB, c, 8)
		wg.Done()
	}()

	wg.Wait()

	matched := 0
	for i := 0; i < c; i++ {
		if j.GenA[i] == j.GenB[i] {
			matched++
		}
	}

	return matched
}

func parallelCompute(gen *Generator, j Judge, vals *[]int, c int, m int) {
	for len(*vals) < c {
		processGen(gen, j)
		p := gen.Previous & 65535

		if p%m == 0 {
			*vals = append(*vals, p)
		}
	}
}

func processGen(gen *Generator, j Judge) {
	n := gen.Previous * gen.Factor
	gen.Previous = n % j.Div
}

func parseInput(in input) []Generator {
	factors := map[string]int{
		"A": 16807,
		"B": 48271,
	}

	reg := regexp.MustCompile("Generator (A|B) starts with ([0-9]+)")

	gens := []Generator{}
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) == 3 {
			prev, _ := strconv.Atoi(m[2])
			gen := Generator{Previous: prev, Factor: factors[m[1]]}
			gens = append(gens, gen)
		}
	}
	return gens
}
