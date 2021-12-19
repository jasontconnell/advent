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

func (reg Register) String() string {
	return reg.Name + ": " + strconv.Itoa(reg.Value)
}

type Instruction struct {
	Raw     string
	Command string

	Arg1    int
	Arg1Reg *Register

	Arg2    int
	Arg2Reg *Register

	Toggled bool
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
	fmt.Fprintln(w, "--2016 day 25 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	insts, regs := parseInput(in)
	return solve(insts, regs)
}

func part2(in input) string {
	return "no part 2 in day 25 :)"
}

func solve(instructions []Instruction, registers map[string]*Register) int {
	i := 1
	ans := 0
	solved := false
	for !solved {
		zeroregs(registers)
		registers["a"].Value = i
		if run(instructions, registers) {
			ans = i
			solved = true
		}
		i++
	}
	return ans
}

func parseInput(in input) ([]Instruction, map[string]*Register) {
	reg := regexp.MustCompile("^(tgl|dec|inc|cpy|jnz|out) (\\-?[a-z0-9\\-]+) ?(\\-?[a-z0-9]+)?")
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
				instruction.Arg1Reg = registers[groups[2]]
				instruction.Arg1 = -1
				break
			case "inc":
				instruction.Arg1Reg = registers[groups[2]]
				instruction.Arg1 = 1
				break
			case "cpy", "jnz":
				if v, err := strconv.Atoi(groups[2]); err == nil {
					instruction.Arg1 = v
				} else {
					instruction.Arg1Reg = registers[groups[2]]
				}

				if v, err := strconv.Atoi(groups[3]); err != nil {
					instruction.Arg2Reg = registers[groups[3]]
				} else {
					instruction.Arg2 = v
				}
				break
			case "tgl", "out":
				instruction.Arg1Reg = registers[groups[2]]
				break
			}
			instructions = append(instructions, instruction)
		}
	}

	return instructions, registers
}

func zeroregs(registers map[string]*Register) {
	for k, _ := range registers {
		registers[k].Value = 0
	}
}

func run(instructions []Instruction, registers map[string]*Register) bool {
	i := 0
	iterations := 0
	lastsig := -1
	result := false
	repeat := 0

	for i < len(instructions) && repeat < 1000 {
		inst := instructions[i]
		switch inst.Command {
		case "inc", "dec":
			inst.Arg1Reg.Value += inst.Arg1
			i++
			break
		case "jnz":
			zero := inst.Arg1 == 0
			if inst.Arg1Reg != nil {
				zero = inst.Arg1Reg.Value == 0
			}

			jmp := inst.Arg2
			if inst.Arg2Reg != nil {
				jmp = inst.Arg2Reg.Value
			}

			if !zero {
				i += jmp
			} else {
				i++
			}
			break
		case "cpy":
			v := inst.Arg1
			if inst.Arg1Reg != nil {
				v = inst.Arg1Reg.Value
			}

			if inst.Arg2Reg != nil { // should only be able to copy to a register
				inst.Arg2Reg.Value = v
			}

			i++
			break
		case "tgl":
			v := inst.Arg1
			if inst.Arg1Reg != nil {
				v = i + inst.Arg1Reg.Value
			}
			if v != -1 && v < len(instructions) {
				switch instructions[v].Command {
				case "inc":
					instructions[v].Command = "dec"
					instructions[v].Arg1 = -1
				case "dec", "tgl":
					instructions[v].Command = "inc"
					instructions[v].Arg1 = 1
				case "jnz":
					instructions[v].Command = "cpy"
				case "cpy":
					instructions[v].Command = "jnz"
				}

				instructions[v].Toggled = true
			}

			i++
			break
		case "out":
			if inst.Arg1Reg != nil {
				sig := inst.Arg1Reg.Value

				if sig == 0 || sig == 1 {
					if lastsig == -1 {
						lastsig = sig
						repeat++
					} else if lastsig == sig {
						i = 10000 // not alternating, break loop
					} else {
						repeat++
						lastsig = sig

						if repeat > 998 {
							result = true
						}
					}
				} else {
					i = 10000 // not a 0 or 1, break loop
				}
			} else {
				i = 10000
			}
			i++
			break
		}

		iterations++
	}

	return result
}
