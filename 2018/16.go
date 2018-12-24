package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "16.txt"

type inputs struct {
	a, b int
}

type cpuobs struct { // cpu observation
	sid         int
	before      []int
	instruction int
	in          inputs
	output      int
	after       []int
}

func (obs cpuobs) String() string {
	s := fmt.Sprintf("Before: %v\nInstruction: %d Inputs: %d %d Output: %d\nAfter: %v\n\n", obs.before, obs.instruction, obs.in.a, obs.in.b, obs.output, obs.after)
	return s
}

type instruction struct {
	id   int
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

	observations := getCPUResults(lines)
	instmap := getInstructions()

	p1 := behaviorCounts(observations, instmap, 3)

	fmt.Println("Part 1:", p1)

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

func behaviorCounts(obs []cpuobs, imap map[string]instruction, c int) int {
	srmap := make(map[int]int)
	for _, o := range obs {
		for k, inst := range imap {
			rc := inst.f(o.before, o.in, o.output)
			if eq(rc, o.after) {
				imap[k].ops[o.instruction] = true
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

func getCPUResults(lines []string) []cpuobs {
	obs := []cpuobs{}
	reg := regexp.MustCompile(`(Before|After)?:?( *?)(\[?)(\d+),? (\d+),? (\d+),? (\d+)(\]?)`)
	for i := 0; i < len(lines); i += 4 {
		if lines[i] == "" {
			break
		}
		before := reg.FindAllStringSubmatch(lines[i], -1)
		middle := reg.FindAllStringSubmatch(lines[i+1], -1)
		after := reg.FindAllStringSubmatch(lines[i+2], -1)
		// lines[i+3] is blank

		bints := getInts(before[0], 4, 4)
		ins := getInts(middle[0], 4, 4)
		aints := getInts(after[0], 4, 4)

		ob := cpuobs{sid: i, before: bints, after: aints, instruction: ins[0], in: inputs{a: ins[1], b: ins[2]}, output: ins[3]}
		obs = append(obs, ob)
	}
	return obs
}

func getInts(str []string, startpos, count int) []int {
	ints := []int{}

	for i := startpos; i < startpos+count; i++ {
		num, _ := strconv.Atoi(str[i])
		ints = append(ints, num)
	}
	return ints
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
