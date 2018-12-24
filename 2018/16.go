package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	"strings"
	//"math"
)

var input = "16.txt"

type op struct {
	opcode int
	in inputs
	out int
}

type inputs struct {
	a, b int
}

type cpuobs struct { // cpu observation
	sid         int
	before      []int
	operation op
	after       []int
}

func (obs cpuobs) String() string {
	opr := obs.operation
	s := fmt.Sprintf("Before: %v\nInstruction: %d Inputs: %d %d Output: %d\nAfter: %v\n\n", obs.before, opr.opcode, opr.in.a, opr.in.b, opr.out, obs.after)
	return s
}

type instruction struct {
	opcode   int
	name string
	f    func(registers []int, in inputs, out int) []int // returns modified registers

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

	observations, progstart := getCPUResults(lines)
	instmap := getInstructions()

	p1 := behaviorCounts(observations, instmap, 3)
	fmt.Println("Part 1:", p1)

	determineOpcodes(instmap)
	prog := getProgram(lines, progstart)
	p2 := runProgram(prog, instmap)

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func eq(a, b []int) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func runProgram(prog []op, imap map[string]instruction) int {
	r := []int{0,0,0,0}
	opmap := make(map[int]instruction)
	for _, inst := range imap {
		opmap[inst.opcode] = inst
	}

	for _, o := range prog {
		inst := opmap[o.opcode]
		r = inst.f(r, o.in, o.out)
	}

	return r[0]
}

func allDetermined(imap map[string]instruction) bool {
	res := true
	for _, inst := range imap {
		if len(inst.ops) > 1 {
			res = false
		}
	}
	return res
}

func determineOpcodes(imap map[string]instruction) {
	for !allDetermined(imap) {
		for k, inst := range imap {
			if len(inst.ops) == 1 {
				for opcode, _ := range inst.ops {
					inst.opcode = opcode
					removeOpcodes(imap, opcode)
				}
				imap[k] = inst
			}
		}
	}
}

func removeOpcodes(imap map[string]instruction, opcode int) {
	for k, inst := range imap {
		delete(inst.ops, opcode)
		imap[k] = inst
	}
}

func behaviorCounts(obs []cpuobs, imap map[string]instruction, c int) int {
	srmap := make(map[int]int)
	for _, o := range obs {
		for k, inst := range imap {
			rc := inst.f(o.before, o.operation.in, o.operation.out)
			if eq(rc, o.after) {
				imap[k].ops[o.operation.opcode] = true
				srmap[o.sid]++
			}
		}
	}

	count := 0
	for _, s := range srmap {
		if s >= c {
			count++
		}
	}
	return count
}

func cp(a []int) []int {
	return append([]int{}, a...)
}

func getInstructions() map[string]instruction {
	imap := make(map[string]instruction)
	imap["addr"] = instruction{name: "addr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] + r[in.b]
		rc[out] = res
		return rc
	}}

	imap["addi"] = instruction{name: "addi", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] + in.b
		rc[out] = res
		return rc
	}}

	imap["mulr"] = instruction{name: "mulr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] * rc[in.b]
		rc[out] = res
		return rc
	}}

	imap["muli"] = instruction{name: "muli", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] * in.b
		rc[out] = res
		return rc
	}}

	imap["banr"] = instruction{name: "banr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] & r[in.b]
		rc[out] = res
		return rc
	}}

	imap["bani"] = instruction{name: "bani", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] & in.b
		rc[out] = res
		return rc
	}}

	imap["borr"] = instruction{name: "borr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] | r[in.b]
		rc[out] = res
		return rc
	}}

	imap["bori"] = instruction{name: "bori", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		res := r[in.a] | in.b
		rc[out] = res
		return rc
	}}

	imap["setr"] = instruction{name: "setr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		rc[out] = r[in.a]
		return rc
	}}

	imap["seti"] = instruction{name: "seti", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		rc[out] = in.a
		return rc
	}}

	imap["gtir"] = instruction{name: "gtir", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if in.a > r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["gtri"] = instruction{name: "gtri", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] > in.b {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["gtrr"] = instruction{name: "gtrr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] > r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqir"] = instruction{name: "eqir", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if in.a == r[in.b] {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqri"] = instruction{name: "eqri", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
		rc := cp(r)
		v := 0
		if r[in.a] == in.b {
			v = 1
		}
		rc[out] = v
		return rc
	}}

	imap["eqrr"] = instruction{name: "eqrr", ops: make(map[int]bool), f: func(r []int, in inputs, out int) []int {
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

func getCPUResults(lines []string) ([]cpuobs, int) {
	obs := []cpuobs{}
	reg := regexp.MustCompile(`(Before|After)?:?( *?)(\[?)(\d+),? (\d+),? (\d+),? (\d+)(\]?)`)
	last := 0
	for i := 0; i < len(lines); i += 4 {
		if lines[i] == "" {
			last = i
			break
		}
		before := reg.FindAllStringSubmatch(lines[i], -1)
		middle := reg.FindAllStringSubmatch(lines[i+1], -1)
		after := reg.FindAllStringSubmatch(lines[i+2], -1)
		// lines[i+3] is blank

		bints := getInts(before[0], 4, 4)
		ins := getInts(middle[0], 4, 4)
		aints := getInts(after[0], 4, 4)

		opr := op{opcode: ins[0], in: inputs{a: ins[1], b: ins[2]}, out: ins[3]}
		ob := cpuobs{sid: i, before: bints, after: aints, operation: opr}
		obs = append(obs, ob)

	}
	return obs, last
}

func getProgram(lines []string, start int) []op {
	ops := []op{}
	for i := start; i < len(lines); i++ {
		sp := strings.Fields(lines[i])
		if len(sp) != 4 {
			continue
		}

		opcode, _ := strconv.Atoi(sp[0])
		a, _ := strconv.Atoi(sp[1])
		b, _ := strconv.Atoi(sp[2])
		out, _ := strconv.Atoi(sp[3])

		opr := op{opcode: opcode, in: inputs{a: a, b: b}, out: out}
		ops = append(ops, opr)
	}
	return ops
}

func getInts(str []string, startpos, count int) []int {
	ints := []int{}

	for i := startpos; i < startpos+count; i++ {
		num, _ := strconv.Atoi(str[i])
		ints = append(ints, num)
	}
	return ints
}