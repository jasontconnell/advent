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
	registers := []int{0, 0, 0, 0, 0, 0}

	registers = run(instreg, instrs, prog, registers)

	fmt.Println("Part 1:", registers[0])

	fmt.Println("Time", time.Since(startTime))
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

func cp(a []int) []int {
	return append([]int{}, a...)
}

func updateIP(ip int, val int, registers []int) []int {
	registers[ip] = val
	return registers
}

func getInstructions() map[string]instruction {
	imap := make(map[string]instruction)
	imap["addr"] = instruction{name: "addr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] + r[in.b]
		rc[out] = res
		return rc
	}}

	imap["addi"] = instruction{name: "addi", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] + in.b
		rc[out] = res
		return rc
	}}

	imap["mulr"] = instruction{name: "mulr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] * rc[in.b]
		rc[out] = res
		return rc
	}}

	imap["muli"] = instruction{name: "muli", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] * in.b
		rc[out] = res
		return rc
	}}

	imap["banr"] = instruction{name: "banr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] & r[in.b]
		rc[out] = res
		return rc
	}}

	imap["bani"] = instruction{name: "bani", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] & in.b
		rc[out] = res
		return rc
	}}

	imap["borr"] = instruction{name: "borr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] | r[in.b]
		rc[out] = res
		return rc
	}}

	imap["bori"] = instruction{name: "bori", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] | in.b
		rc[out] = res
		return rc
	}}

	imap["setr"] = instruction{name: "setr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		rc[out] = r[in.a]
		return rc
	}}

	imap["seti"] = instruction{name: "seti", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		rc[out] = in.a
		return rc
	}}

	imap["gtir"] = instruction{name: "gtir", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if in.a > r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["gtri"] = instruction{name: "gtri", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] > in.b {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["gtrr"] = instruction{name: "gtrr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] > r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqir"] = instruction{name: "eqir", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if in.a == r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqri"] = instruction{name: "eqri", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] == in.b {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqrr"] = instruction{name: "eqrr", ops: make(map[int]bool), f: func(ip int, r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] == r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
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
