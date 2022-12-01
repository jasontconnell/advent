package main

import (
	"fmt"
	"github.com/jasontconnell/advent/common"
	"math/rand"
	"testing"
)

func TestCheckIn(t *testing.T) {
	lines, _ := common.ReadStrings("input.txt")
	p := parseInput(lines)

	in := []int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}

	testStep := 5

	duplicateInputs := map[int]int{}

	for x := 0; x < 1000; x++ {
		m := map[int]int{}
		for j := 0; j < 14; j++ {
			if j == testStep {
				continue
			}
			in[j] = (rand.Int()%9 + 1)
		}
		for i := 1; i <= 9; i++ {
			p.reset()

			in[testStep] = i
			// fmt.Println(in)
			p.in(in...)
			p.run()

			z := p.registers["z"].value
			if z == 0 {
				fmt.Println("GOALLLLL", in)
				return
			}
			if _, ok := m[z]; ok {
				duplicateInputs[i]++
			}
			m[z] = i
			if (z%26)-6 < 10 {
				fmt.Println((z%26)-6, in, z)
			}
		}
		if len(m) != 9 {
			fmt.Println(in)
			fmt.Println()
			for k := range m {
				fmt.Print(k, " ")
			}
			fmt.Println()
		}
	}

	fmt.Println("duplicateInputs", duplicateInputs, len(duplicateInputs))

}

func TestTheory(t *testing.T) {
	lines, _ := common.ReadStrings("testpart.txt")
	p := parseInput(lines)

	in := []int{2, 9, 9, 3, 4, 9, 6, 7, 9, 8, 9, 9, 9, 1}
	in2 := []int{2, 5, 5, 3, 4, 5, 6, 7, 5, 8, 5, 5, 5, 1}

	p.in(in...)
	p.run()

	z1 := p.registers["z"].value

	p.reset()
	p.in(in2...)
	p.run()
	z2 := p.registers["z"].value

	if z1 == z2 {
		fmt.Println("Your theory is valid")
	}
}

func TestLoop(t *testing.T) {
	lines, _ := common.ReadStrings("testpart.txt")

	m := map[[3]int][]int{}
	p := parseInput(lines)
	for i := 0; i <= 500; i++ {
		p.reset()
		in := i%9 + 1
		x := rand.Int() % 100
		y := rand.Int() % 100
		zi := in + 12 //rand.Int() % 1
		zi = (20 * zi) + (7 * zi)
		zi = (zi * 26) + (14 + in)
		zi = zi*26 + in
		zi = zi + 3
		zi = zi*26 + (in + 15)
		zi = zi*26 + (in + 11)
		zi = zi + 1
		zi = (zi * 26) + 1 + in
		zi = zi + 11
		zi = zi - 9
		zi = zi + 7
		p.registers["x"].value = x
		p.registers["y"].value = y
		p.registers["z"].value = zi
		p.in(in)
		p.run()
		if true {
			k := [3]int{x, y, zi}
			m[k] = []int{in, p.registers["x"].value, p.registers["y"].value, p.registers["z"].value}
		}
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println()
}

func TestPart(t *testing.T) {
	lines, _ := common.ReadStrings("testpart.txt")

	in := []int{9, 1, 9}

	p := parseInput(lines)
	testStep := 0

	m := map[int]int{}
	for j := 0; j < len(in); j++ {
		for i := 1; i <= 9; i++ {
			p.reset()

			in[testStep] = i
			// fmt.Println(in)
			p.in(in...)
			p.run()

			z := p.registers["z"].value
			if z == 0 {
				fmt.Println("GOALLLLL", in)
			}

			m[z] = i
		}
	}

	fmt.Println(m)
}
