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

var input = "08.txt"

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

func (r Register) String() string {
	return fmt.Sprintf("%v = %v", r.Name, r.Value)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	insts := []*Instruction{}

	for scanner.Scan() {
		var txt = scanner.Text()
		inst := getInstruction(txt)
		if inst != nil {
			insts = append(insts, inst)
		}
	}

	rmap := getMap(insts)

	maxp2 := process(rmap, insts)

	max := -999999
	for _, reg := range rmap {
		if reg.Value > max {
			max = reg.Value
		}
	}

	fmt.Println("max after processing   (p1)     ", max)
	fmt.Println("max during processing  (p2)     ", maxp2)

	fmt.Println("Time", time.Since(startTime))
}

var reg = regexp.MustCompile("^([a-z]+) ([a-z]+) ([0-9\\-]+) if ([a-z]+) ([!=-><]+) ([0-9\\-]+)$")

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
