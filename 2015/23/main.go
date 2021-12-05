package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string
type output int

var reg *regexp.Regexp = regexp.MustCompile(`([a-z]+) (\-?\+?[a-z0-9]+),? ?(\-?\+?[a-z0-9]+)?$`)

type Instruction struct {
	Text     string
	Register string
	Value    int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	inst, regs := getInput(in)
	Execute(inst, regs)

	return output(regs["b"])
}

func part2(in input) output {
	inst, regs := getInput(in)

	regs["a"] = 1
	Execute(inst, regs)

	return output(regs["b"])
}

func getInput(in input) ([]Instruction, map[string]int) {
	instructionSet := []Instruction{}
	registers := make(map[string]int)
	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			instruction := Instruction{Text: groups[1]}
			if groups[2] == "a" || groups[2] == "b" {
				instruction.Register = groups[2]

				if _, exists := registers[instruction.Register]; !exists {
					registers[instruction.Register] = 0
				}
			} else {
				val, _ := strconv.Atoi(groups[2])
				instruction.Value = val
			}

			if groups[3] != "" {
				val, _ := strconv.Atoi(groups[3])
				instruction.Value = val
			}

			instructionSet = append(instructionSet, instruction)
		}
	}
	return instructionSet, registers
}

func Execute(list []Instruction, registers map[string]int) {
	done := false
	current := 0

	for !done {
		inst := list[current]

		switch inst.Text {
		case "hlf":
			registers[inst.Register] = registers[inst.Register] / 2
			current++
		case "tpl":
			registers[inst.Register] = registers[inst.Register] * 3
			current++
		case "inc":
			registers[inst.Register] = registers[inst.Register] + 1
			current++
		case "jmp":
			current = current + inst.Value
		case "jie":
			if registers[inst.Register]%2 == 0 {
				current = current + inst.Value
			} else {
				current++
			}
		case "jio":
			if registers[inst.Register] == 1 {
				current = current + inst.Value
			} else {
				current++
			}
		}

		done = current > len(list)-1

	}
	return
}
