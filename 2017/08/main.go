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

type Condition struct {
	Register string
	Comp     string
	Value    int
}

type Instruction struct {
	Register string
	Op       string
	Value    int
	Cond     Condition
}

type Register struct {
	Name  string
	Value int
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
	fmt.Fprintln(w, "--2017 day 08 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return 0
}

func part2(in input) output {
	return 0
}

func parseInput(in input) []*Instruction {
	insts := []*Instruction{}
	for _, line := range in {
		inst := getInstruction(line)
		insts = append(insts, inst)
	}
	return insts
}

func process(rmap map[string]*Register, insts []*Instruction) int {
	max := 0
	for _, v := range insts {
		c := evalInstruction(rmap, v)
		if c > max {
			max = c
		}
	}
	return max
}

func evalInstruction(rmap map[string]*Register, inst *Instruction) int {
	reg := rmap[inst.Register]
	res := evalCondition(rmap, inst.Cond)
	if res {
		switch inst.Op {
		case "inc":
			reg.Value += inst.Value
		case "dec":
			reg.Value -= inst.Value
		}
	}

	return reg.Value
}

func evalCondition(rmap map[string]*Register, cond Condition) bool {
	reg, ok := rmap[cond.Register]
	if !ok {
		fmt.Println("couldn't find register", cond.Register)
		return false
	}

	val := reg.Value
	res := false
	switch cond.Comp {
	case "==":
		res = val == cond.Value
	case "!=":
		res = val != cond.Value
	case "<=":
		res = val <= cond.Value
	case ">":
		res = val > cond.Value
	case "<":
		res = val < cond.Value
	case ">=":
		res = val >= cond.Value
	default:
		fmt.Println("not found", cond.Comp)
	}

	return res
}

func getMap(insts []*Instruction) map[string]*Register {
	rmap := make(map[string]*Register)
	for _, i := range insts {
		rmap[i.Register] = &Register{Name: i.Register, Value: 0}
	}
	return rmap
}

var reg = regexp.MustCompile("^([a-z]+) ([a-z]+) ([0-9\\-]+) if ([a-z]+) ([!=-><]+) ([0-9\\-]+)$")

func getInstruction(line string) *Instruction {
	var inst *Instruction
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		r := groups[1]
		op := groups[2]
		val, err := strconv.Atoi(groups[3])
		if err != nil {
			fmt.Println("parse", err, groups[3])
			return nil
		}

		r2 := groups[4]
		cmp := groups[5]
		val2, err := strconv.Atoi(groups[6])
		if err != nil {
			fmt.Println("parse", err, groups[6])
			return nil
		}

		inst = &Instruction{Register: r, Op: op, Value: val, Cond: Condition{Register: r2, Comp: cmp, Value: val2}}
	}

	return inst
}
