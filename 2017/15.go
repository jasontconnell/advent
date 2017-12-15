package main

import (
	"fmt"
	"sync"
	"time"
)

type Generator struct {
	Previous int
	Factor   int
}

type Judge struct {
	Div  int
	GenA []int
	GenB []int
}

var genA Generator = Generator{Previous: 516, Factor: 16807}
var genB Generator = Generator{Previous: 190, Factor: 48271}

var genAtest Generator = Generator{Previous: 65, Factor: 16807}
var genBtest Generator = Generator{Previous: 8921, Factor: 48271}

var judge Judge = Judge{Div: 2147483647}

func main() {
	startTime := time.Now()

	matched := compute(genA, genB, judge, 40000000)

	matched2 := computePart2(genA, genB, judge, 5000000)

	fmt.Println("matched         :", matched)
	fmt.Println("matched part 2  :", matched2)

	fmt.Println("Time", time.Since(startTime))
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

	go parallelCompute(&wg, &a, j, &j.GenA, c, 4)
	go parallelCompute(&wg, &b, j, &j.GenB, c, 8)

	wg.Wait()

	matched := 0
	for i := 0; i < c; i++ {
		if j.GenA[i] == j.GenB[i] {
			matched++
		}
	}

	return matched
}

func parallelCompute(wg *sync.WaitGroup, gen *Generator, j Judge, vals *[]int, c int, m int) {
	for len(*vals) < c {
		processGen(gen, j)
		p := gen.Previous & 65535

		if p%m == 0 {
			*vals = append(*vals, p)
		}
	}

	fmt.Println("done")
	wg.Done()
}

func processGen(gen *Generator, j Judge) {
	n := gen.Previous * gen.Factor
	gen.Previous = n % j.Div
}
