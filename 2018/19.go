package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	"strings"
	//"math"
)

var input = "19.txt"

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

	p1reg := run(instreg, instrs, prog, []int{0, 0, 0, 0, 0, 0})
	fmt.Println("Part 1:", p1reg[0])

	//p2reg := run(instreg, instrs, prog, []int{1, 0, 0, 0, 0, 0})
	fmt.Println("Part 2:", p2crap(10551364))


	fmt.Println("Time", time.Since(startTime))
}

func p2crap(n int) int {
	s := 0

	for i := 1; i < n; i++ {
		if n % i == 0 {
			s += i
		}
	}
	return s+n
}

func run(ip int, instrs map[string]instruction, prog []op, registers []int) []int {
	idx := 0

	for idx >= 0 && idx < len(prog) {
		opr := prog[idx]
		inst := instrs[opr.opcode]
		registers = inst.f(ip, registers, opr.in, opr.out)
		registers[ip]++
		idx = registers[ip]
	}
	return registers
}

func ipcheck(r []int, ip, out int) []int {
	if ip == out {
		r[ip]--
	}
	return r
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
