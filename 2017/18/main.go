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
	reg *register
	val int
}

type register struct {
	name  string
	value int
}

type program struct {
	registers    map[string]*register
	instructions []instruction
	queue        []int
	parallel     *program
	deadlocked   bool
	sends        int
	instruction  int
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
	return 0
}

func runProgram(instrs []instruction, regs map[string]*register) int {
	cur := 0
	snd := 0
	for cur < len(instrs) {
		step := execInstruction(instrs[cur], regs, func(i int) {
			snd = i
		}, func(i int) {
			if i > 0 {
				cur = len(instrs)
			}
		})
		cur += step
	}
	return snd
}

func execInstruction(inst instruction, regs map[string]*register, snd func(i int), rcv func(i int)) int {
	incr := 1
	switch inst.cmd {
	case "snd":
		l := operandValue(inst.args[0], regs)
		snd(l)
	case "rcv":
		l := operandValue(inst.args[0], regs)
		rcv(l)
	case "set":
		l := inst.args[0].reg
		l.value = operandValue(inst.args[1], regs)
	case "add":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg.name].value = operandValue(l, regs) + operandValue(r, regs)
	case "mul":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg.name].value = operandValue(l, regs) * operandValue(r, regs)
	case "mod":
		l := inst.args[0]
		r := inst.args[1]
		regs[l.reg.name].value = operandValue(l, regs) % operandValue(r, regs)
	case "jgz":
		l := inst.args[0]
		r := inst.args[1]
		x := operandValue(l, regs)
		jmp := operandValue(r, regs)

		if x > 0 && jmp != 0 {
			incr = jmp
		}
	}
	return incr
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

func operandValue(opr operand, registers map[string]*register) int {
	if opr.reg != nil {
		return opr.reg.value
	}
	return opr.val
}

func parseOperand(str string, registers map[string]*register) operand {
	val, err := strconv.Atoi(str)
	var reg *register
	if err != nil {
		if _, ok := registers[str]; ok {
			reg = registers[str]
		} else {
			reg = &register{name: str}
			registers[str] = reg
		}
	}

	return operand{reg: reg, val: val}
}
