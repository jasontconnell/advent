package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var input = "08.txt"

type patch struct {
	op       string
	location int
}
type program struct {
	instructions []instr
	location     int
	acc          int
	patch        *patch
}

type instr struct {
	op  string
	val int
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

	prog := getProgram(lines)
	p1, _ := exec(prog)

	fmt.Println("Part 1:", p1)
	patches := getPatches(prog)
	p2 := p2exec(prog, patches)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getPatches(p program) []patch {
	patches := []patch{}
	for i, inst := range p.instructions {
		switch inst.op {
		case "nop", "jmp":
			pop := "jmp"
			if inst.op == "jmp" {
				pop = "nop"
			}
			pch := patch{op: pop, location: i}
			patches = append(patches, pch)
		}
	}
	return patches
}

func p2exec(p program, patches []patch) int {
	acc := 0
	for _, pch := range patches {
		p.acc = 0
		p.patch = &pch
		tmp, inf := exec(p)

		if !inf {
			acc = tmp
			break
		}
		p.patch = nil
	}
	return acc
}

func exec(p program) (int, bool) {
	visited := make(map[int]int)
	done := false
	infinite := false

	for !done {
		if p.location >= len(p.instructions) {
			done = true
			infinite = false
			break
		}

		ins := p.instructions[p.location]
		if _, ok := visited[p.location]; ok {
			infinite = true
			done = true
			break
		}
		visited[p.location] = p.acc

		nins := instr{op: ins.op, val: ins.val}

		if p.patch != nil && p.patch.location == p.location {
			nins.op = p.patch.op
		}

		switch nins.op {
		case "acc":
			p.acc += nins.val
			p.location++
		case "jmp":
			p.location += nins.val
		case "nop":
			p.location++
			break
		}
	}

	return p.acc, infinite
}

var opreg *regexp.Regexp = regexp.MustCompile("^(acc|nop|jmp) ([+|-][0-9]+)$")

func getProgram(lines []string) program {
	p := program{}
	for _, line := range lines {
		groups := opreg.FindStringSubmatch(line)
		if groups == nil || len(groups) != 3 {
			continue
		}

		op := groups[1]
		i, err := strconv.Atoi(groups[2])
		if err != nil {
			fmt.Println(err)
		}

		inst := instr{op: op, val: i}
		p.instructions = append(p.instructions, inst)
	}
	return p
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
