package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type instruction struct {
	cmd  string
	args []operand
}

type operand struct {
	reg string
	val int
}

type register struct {
	name  string
	value int
}

type program struct {
	registers   map[string]*register
	instr       []instruction
	queue       []int
	sends       int
	instruction int
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
	fmt.Fprintln(w, "--2017 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	inst, reg := parseInput(in)
	return runProgram(inst, reg)
}

func part2(in input) output {
	inst0, reg0 := parseInput(in)
	inst1, reg1 := parseInput(in)

	p0 := &program{instr: inst0, registers: reg0}
	p1 := &program{instr: inst1, registers: reg1}
	runPrograms(p0, p1)
	return p1.sends
}

func runPrograms(p1, p2 *program) {
	var p1step, p2step int
	p1wait, p2wait := false, false

	for !p1wait && !p2wait {
		p1step = execInstruction(p1.instr[p1.instruction], p1.registers, func(opr operand) {
			p2.queue = append(p2.queue, operandValue(opr, p1.registers))
			p1.sends++
		}, func(opr operand) {
			if len(p1.queue) == 0 {
				p1wait = true
				return
			}
			p1.registers[opr.reg].value = p1.queue[0]
			p1.queue = p1.queue[1:]
			p1wait = false
		})

		p2step = execInstruction(p2.instr[p2.instruction], p2.registers, func(opr operand) {
			p1.queue = append(p1.queue, operandValue(opr, p2.registers))
			p2.sends++
		}, func(opr operand) {
			if len(p2.queue) == 0 {
				p2wait = true
				return
			}
			p2.registers[opr.reg].value = p2.queue[0]
			p2.queue = p2.queue[1:]
			p2wait = false
		})

		if !p1wait {
			p1.instruction += p1step
		}

		if !p2wait {
			p2.instruction += p2step
		}
	}
}

func runProgram(instrs []instruction, regs map[string]*register) int {
	cur := 0
	snd := 0
	terminate := false
	for !terminate && cur < len(instrs) {
		step := execInstruction(instrs[cur], regs, func(opr operand) {
			snd = operandValue(opr, regs)
		}, func(opr operand) {
			i := operandValue(opr, regs)
			terminate = i > 0
		})
		cur += step
	}
	return snd
}

func execInstruction(inst instruction, regs map[string]*register, snd func(opr operand), rcv func(opr operand)) int {
	incr := 1
	switch inst.cmd {
	case "snd":
		snd(inst.args[0])
	case "rcv":
		rcv(inst.args[0])
	case "set":
		l := inst.args[0]
		regs[l.reg].value = operandValue(inst.args[1], regs)
	case "add":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg].value = operandValue(l, regs) + operandValue(r, regs)
	case "mul":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg].value = operandValue(l, regs) * operandValue(r, regs)
	case "mod":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg].value = operandValue(l, regs) % operandValue(r, regs)
	case "jgz":
		l := inst.args[0]
		r := inst.args[1]
		x := operandValue(l, regs)
		jmp := operandValue(r, regs)

		if x > 0 {
			incr = jmp
		}
	}
	return incr
}

func operandValue(opr operand, registers map[string]*register) int {
	if opr.reg != "" {
		return registers[opr.reg].value
	}
	return opr.val
}

func parseInput(in input) ([]instruction, map[string]*register) {
	instr := []instruction{}
	regs := make(map[string]*register)
	for _, line := range in {
		sp := strings.Split(line, " ")
		inst := instruction{cmd: sp[0]}
		inst.args = append(inst.args, parseOperand(sp[1], regs))
		if len(sp) == 3 {
			inst.args = append(inst.args, parseOperand(sp[2], regs))
		}

		instr = append(instr, inst)
	}

	return instr, regs
}

func parseOperand(str string, registers map[string]*register) operand {
	val, err := strconv.Atoi(str)
	var reg *register
	opr := operand{}
	if err != nil {
		opr.reg = str
		if _, ok := registers[str]; ok {
			reg = registers[str]
		} else {
			reg = &register{name: str}
			registers[str] = reg
		}
	} else {
		opr.val = val
	}

	return opr
}
