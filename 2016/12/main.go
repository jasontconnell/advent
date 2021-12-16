package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Register struct {
	Name  string
	Value int
}

type Instruction struct {
	Raw            string
	Command        string
	TargetRegister *Register
	TargetValue    int

	Value         int
	ValueRegister *Register
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 12 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	regs, inst := parseInput(in)
	run(inst)

	return regs["a"].Value
}

func part2(in input) output {
	regs, inst := parseInput(in)
	regs["c"].Value = 1
	run(inst)
	return regs["a"].Value
}

func run(instructions []Instruction) {
	i := 0
	iterations := 0

	for i < len(instructions) {
		inst := instructions[i]

		switch inst.Command {
		case "inc", "dec":
			inst.TargetRegister.Value += inst.Value
			i++
			break
		case "jnz":
			if inst.TargetRegister != nil && inst.TargetRegister.Value != 0 {
				i += inst.Value
			} else if inst.TargetRegister == nil && inst.TargetValue != 0 {
				i += inst.Value
			} else {
				i++
			}
			break
		case "cpy":
			val := inst.Value
			if inst.ValueRegister != nil {
				val = inst.ValueRegister.Value
			}
			inst.TargetRegister.Value = val

			i++
			break
		}

		iterations++
	}
}

func parseInput(in input) (map[string]*Register, []Instruction) {
	reg := regexp.MustCompile("^(dec|inc|cpy|jnz) (\\-?[a-z0-9\\-]+) ?(\\-?[a-z0-9]+)?")
	instructions := []Instruction{}
	registers := make(map[string]*Register)

	registers["a"] = &Register{Name: "a", Value: 0}
	registers["b"] = &Register{Name: "b", Value: 0}
	registers["c"] = &Register{Name: "c", Value: 0}
	registers["d"] = &Register{Name: "d", Value: 0}

	for _, line := range in {
		if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
			instruction := Instruction{Command: groups[1], Raw: groups[0]}

			switch groups[1] {
			case "dec":
				instruction.TargetRegister = registers[groups[2]]
				instruction.Value = -1
				break
			case "inc":
				instruction.TargetRegister = registers[groups[2]]
				instruction.TargetRegister = registers[groups[2]]
				instruction.Value = 1
				break
			case "cpy":
				if v, err := strconv.Atoi(groups[2]); err == nil {
					instruction.Value = v
				} else {
					instruction.ValueRegister = registers[groups[2]]
				}

				instruction.TargetRegister = registers[groups[3]]
				break
			case "jnz":
				if v, err := strconv.Atoi(groups[2]); err != nil {
					instruction.TargetRegister = registers[groups[2]]
				} else {
					instruction.TargetValue = v
				}

				if v, err := strconv.Atoi(groups[3]); err == nil {
					instruction.Value = v
				}

				break
			}

			instructions = append(instructions, instruction)

		}
	}

	return registers, instructions

}
