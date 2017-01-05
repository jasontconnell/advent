package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "regexp"
    "strconv"
)
var input = "23.txt"

type Register struct {
    Name string
    Value int
}

func (reg Register) String() string {
    return reg.Name + ": " + strconv.Itoa(reg.Value)
}

type Instruction struct {
    Raw string
    Command string

    Arg1 int
    Arg1Reg *Register

    Arg2 int
    Arg2Reg *Register

    Toggled bool
}

func (inst Instruction) String() string {
    a1 := strconv.Itoa(inst.Arg1)
    if inst.Arg1Reg != nil {
        a1 = inst.Arg1Reg.Name
        a1 = a1 + "(" + strconv.Itoa(inst.Arg1Reg.Value) + ")"
    }

    a2 := strconv.Itoa(inst.Arg2)
    if inst.Arg2Reg != nil {
        a2 = inst.Arg2Reg.Name
        a2 = a2 + "(" + strconv.Itoa(inst.Arg2Reg.Value) + ")"
    }
    return inst.Command + " " + a1 + " " + a2
}

func main() {
    startTime := time.Now()
    if f, err := os.Open(input); err == nil {
        scanner := bufio.NewScanner(f)

        reg := regexp.MustCompile("^(tgl|dec|inc|cpy|jnz) (\\-?[a-z0-9\\-]+) ?(\\-?[a-z0-9]+)?")
        instructions := []Instruction{}
        registers := make(map[string]*Register)

        registers["a"] = &Register{ Name: "a", Value: 12 } // part 1 is 7
        registers["b"] = &Register{ Name: "b", Value: 0 }
        registers["c"] = &Register{ Name: "c", Value: 0 }
        registers["d"] = &Register{ Name: "d", Value: 0 }

        for scanner.Scan() {
            var txt = scanner.Text()
            if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                instruction := Instruction{ Command: groups[1], Raw: groups[0] }

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
                    if v,err := strconv.Atoi(groups[2]); err == nil {
                        instruction.Arg1 = v
                    } else {
                        instruction.Arg1Reg = registers[groups[2]]
                    }

                    if v,err := strconv.Atoi(groups[3]); err != nil {
                        instruction.Arg2Reg = registers[groups[3]]
                    } else {
                        instruction.Arg2 = v
                    }
                    break
                case "tgl":
                        instruction.Arg1Reg = registers[groups[2]]
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
                        case "jnz": instructions[v].Command = "cpy"
                        case "cpy": instructions[v].Command = "jnz"
                    }

                    instructions[v].Toggled = true
                }


                i++
                break
            }
            
            iterations++
        }

        for _, r := range registers {
            fmt.Println(r)
        }
    }

    fmt.Println("Time", time.Since(startTime))
}
