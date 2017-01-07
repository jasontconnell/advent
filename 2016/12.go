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

var input = "12.txt"

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
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		reg := regexp.MustCompile("^(dec|inc|cpy|jnz) (\\-?[a-z0-9\\-]+) ?(\\-?[a-z0-9]+)?")
		instructions := []Instruction{}
		registers := make(map[string]*Register)

		registers["a"] = &Register{Name: "a", Value: 0}
		registers["b"] = &Register{Name: "b", Value: 0}
		registers["c"] = &Register{Name: "c", Value: 1} // 0 for part 1
		registers["d"] = &Register{Name: "d", Value: 0}

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
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

		for k, r := range registers {
			fmt.Println(k, r)
		}
	}

	fmt.Println("Time", time.Since(startTime))
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
