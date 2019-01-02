package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
)

var input = "23.txt"

type Instruction struct {
	Text     string
	Register string
	Value    int
}

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		instructionSet := []Instruction{}
		registers := make(map[string]int)

		for scanner.Scan() {
			var txt = scanner.Text()

			reg := regexp.MustCompile(`([a-z]+) (\-?\+?[a-z0-9]+),? ?(\-?\+?[a-z0-9]+)?$`)
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

		Execute(instructionSet, registers)
		fmt.Println(registers)

		// part 2
		registers["a"] = 1
		registers["b"] = 0
		Execute(instructionSet, registers)
		fmt.Println(registers)

	}

	fmt.Println("Time", time.Since(startTime))
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
