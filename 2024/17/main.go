package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = string

type register struct {
	name string
	data int
}

type op struct {
	id      int
	operand int
}

type program struct {
	ptr int
	ops []int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	regs, prog := parse(in)
	return runProgram(prog, regs)
}

func part2(in input) output {
	return ""
}

func runProgram(prog program, regs []register) string {
	b := bytes.NewBufferString("")
	rm := make(map[string]register)
	for _, r := range regs {
		rm[r.name] = r
	}

	getVal := func(i int) int {
		switch i {
		case 0, 1, 2, 3, 7:
			return i
		case 4:
			return rm["A"].data
		case 5:
			return rm["B"].data
		case 6:
			return rm["C"].data
		}
		return -1
	}

	setVal := func(r string, i int) {
		reg := rm[r]
		reg.data = i
		rm[r] = reg
	}

	outcount := 0
	for {
		if prog.ptr >= len(prog.ops) {
			break
		}
		cur := prog.ops[prog.ptr]
		operand := prog.ops[prog.ptr+1]

		dojump := true
		jmp := 2
		switch cur {
		case 0:
			n := rm["A"].data
			d := int(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("A", v)
			}
		case 1:
			n := rm["B"].data
			v := n ^ operand
			setVal("B", v)
		case 2:
			n := getVal(operand)
			n = n % 8
			setVal("B", n)
		case 3:
			jnz := rm["A"].data
			if jnz != 0 {
				dojump = false
				prog.ptr = operand
			}
		case 4:
			b := rm["B"].data
			c := rm["C"].data
			v := b ^ c
			setVal("B", v)
		case 5:
			v := getVal(operand) % 8
			if outcount > 0 {
				fmt.Fprint(b, ",")
			}
			fmt.Fprintf(b, "%d", v)
			outcount++
		case 6:
			n := rm["A"].data
			d := int(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("B", v)
			}
		case 7:
			n := rm["A"].data
			d := int(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("C", v)
			}
		}
		if dojump {
			prog.ptr += jmp
			fmt.Println("ptr", prog.ptr)
		}

	}
	return b.String()
}

func parse(in []string) ([]register, program) {
	rreg := regexp.MustCompile("^Register ([ABC]): ([0-9]+)$")
	preg := regexp.MustCompile("^Program: ([0-9,]+)$")

	regs := []register{}
	prog := program{}

	for _, line := range in {
		rm := rreg.FindStringSubmatch(line)
		if len(rm) != 0 {
			name := rm[1]
			val, _ := strconv.Atoi(rm[2])
			reg := register{name: name, data: val}
			regs = append(regs, reg)
		}
		pm := preg.FindStringSubmatch(line)
		if len(pm) != 0 {
			ops := strings.Split(pm[1], ",")
			vals := []int{}
			for _, s := range ops {
				v, _ := strconv.Atoi(s)
				vals = append(vals, v)
			}
			prog.ptr = 0
			prog.ops = vals
		}
	}
	return regs, prog
}
