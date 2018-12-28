package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	//"sort"
)

var input = "21.txt"

type op struct {
	opcode string
	in     inputs
	out    int
}

type inputs struct {
	a, b int
}

type instruction struct {
	opcode int
	name   string
	f      func(ip int, registers []int, in inputs, out int) []int // returns modified registers

	ops map[int]bool
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	instreg := getInstReg(lines[0])
	prog := getProgram(lines[1:])
	instrs := getInstructions()

	p1, cycles := solvep1(instreg, instrs, prog, 2000)
	//p2 := solvep2(instreg, instrs, prog, cycles)

	fmt.Println("Part 1:", p1, cycles)
	//fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func solvep1(ip int, instrs map[string]instruction, prog []op, maxcycles int) (int, int) {
	i := 0
	_, c := run(ip, instrs, prog, []int{i, 0, 0, 0, 0, 0}, maxcycles)
	return i, c
}

func solvep2(ip int, instrs map[string]instruction, prog []op, maxcycles int) int {
	r := []int{0, 0, 0, 0, 0, 0}
	i := 0
	max := 0
	r[0] = i
	//solves := []int{}
	for {
		reg, c := run(ip, instrs, prog, []int{i, 0, 0, 0, 0, 0}, maxcycles*1000)
		if c > maxcycles*1000 {
			maxcycles = c
			fmt.Println(c)
		} else if c == -1 {
			fmt.Println("force halted", i, "at", maxcycles*1000)
		}
		i++
		fmt.Println("solved", reg, c, " ----- ", max)
	}
	return i
}

func run(ip int, instrs map[string]instruction, prog []op, registers []int, maxcycles int) ([]int, int) {
	idx := 0
	cycles := 0
	forceHalt := false
	m := make(map[int]int)
	d := make(map[int]int)
	for idx >= 0 && idx < len(prog) && !forceHalt {
		opr := prog[idx]
		inst := instrs[opr.opcode]
		registers = inst.f(ip, registers, opr.in, opr.out)
		if registers[ip] == 28 {
			if _, ok := m[registers[2]]; !ok {
				fmt.Println(registers)
				m[registers[2]] = registers[3]
			}

			if _, ok := d[registers[3]]; ok {
				fmt.Println("DONE")
				os.Exit(0)
			}
			d[registers[3]] = 1
		}

		if registers[ip] < 15 {
			fmt.Println(registers)
		} else {
			os.Exit(0)

		}
		registers[ip]++
		idx = registers[ip]
		cycles++

		forceHalt = cycles == maxcycles
	}

	if forceHalt {
		cycles = -1
	}

	return registers, cycles
}

func getInstructions() map[string]instruction {
	imap := make(map[string]instruction)
	imap["addr"] = instruction{name: "addr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] + r[in.b]
		r[out] = res
		return r
	}}

	imap["addi"] = instruction{name: "addi", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] + in.b
		r[out] = res
		return r
	}}

	imap["mulr"] = instruction{name: "mulr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] * r[in.b]
		r[out] = res
		return r
	}}

	imap["muli"] = instruction{name: "muli", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] * in.b
		r[out] = res
		return r
	}}

	imap["banr"] = instruction{name: "banr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] & r[in.b]
		r[out] = res
		return r
	}}

	imap["bani"] = instruction{name: "bani", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] & in.b
		r[out] = res
		return r
	}}

	imap["borr"] = instruction{name: "borr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] | r[in.b]
		r[out] = res
		return r
	}}

	imap["bori"] = instruction{name: "bori", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		res := r[in.a] | in.b
		r[out] = res
		return r
	}}

	imap["setr"] = instruction{name: "setr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		r[out] = r[in.a]
		return r
	}}

	imap["seti"] = instruction{name: "seti", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		r[out] = in.a
		return r
	}}

	imap["gtir"] = instruction{name: "gtir", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if in.a > r[in.b] {
			v = 1
		}
		r[out] = v
		return r
	}}

	imap["gtri"] = instruction{name: "gtri", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if r[in.a] > in.b {
			v = 1
		}
		r[out] = v
		return r
	}}

	imap["gtrr"] = instruction{name: "gtrr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if r[in.a] > r[in.b] {
			v = 1
		}
		r[out] = v
		return r
	}}

	imap["eqir"] = instruction{name: "eqir", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if in.a == r[in.b] {
			v = 1
		}
		r[out] = v
		return r
	}}

	imap["eqri"] = instruction{name: "eqri", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if r[in.a] == in.b {
			v = 1
		}
		r[out] = v
		return r
	}}

	imap["eqrr"] = instruction{name: "eqrr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		v := 0
		if r[in.a] == r[in.b] {
			v = 1
		}
		r[out] = v
		return r
	}}

	return imap
}

func getInstReg(txt string) int {
	a, _ := strconv.Atoi(strings.Replace(txt, "#ip ", "", -1))
	return a
}

func getProgram(lines []string) []op {
	prog := []op{}
	for _, line := range lines {
		prog = append(prog, getOp(line))
	}
	return prog
}

func getOp(txt string) op {
	ins := strings.Fields(txt)
	opcode := ins[0]
	in1, _ := strconv.Atoi(ins[1])
	in2, _ := strconv.Atoi(ins[2])
	out, _ := strconv.Atoi(ins[3])
	return op{opcode: opcode, in: inputs{a: in1, b: in2}, out: out}
}
