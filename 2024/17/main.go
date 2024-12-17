package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = string

type register struct {
	name string
	data int64
}

type program struct {
	ptr int
	ops []int64
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	regs, prog := parse(in)
	vals := runProgram(prog, regs)
	s := ""
	for i := 0; i < len(vals); i++ {
		s += fmt.Sprintf("%d", vals[i])
		if i != len(vals)-1 {
			s += ","
		}
	}
	return s
}

func part2(in input) int64 {
	regs, prog := parse(in)
	return findQuine(prog, regs)
}

type state struct {
	segs []int64
}

func findQuine(prog program, regs []register) int64 {
	queue := []state{}
	for i := 0; i < 8; i++ {
		queue = append(queue, state{[]int64{int64(i)}})
	}
	var final int64
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		var x int64
		for i := len(cur.segs) - 1; i >= 0; i-- {
			s := cur.segs[i] << (3 * i)
			x = x | s
		}

		regs[0].data = x
		vals := runProgram(prog, regs)
		vp := 0
		matched := true
		for p := len(prog.ops) - len(vals); p < len(prog.ops); p++ {
			if vals[vp] != prog.ops[p] {
				matched = false
				break
			}
			vp++
		}

		done := matched && len(prog.ops) == len(vals)
		if done {
			final = x
			break
		}

		if matched {
			for i := 0; i < 8; i++ {
				nseg := make([]int64, len(cur.segs))
				copy(nseg, cur.segs)
				nseg = append([]int64{int64(i)}, nseg...)
				queue = append(queue, state{nseg})
			}
		}
	}
	return final
}

func runProgram(prog program, regs []register) []int64 {
	output := []int64{}
	rm := make(map[string]register)
	for _, r := range regs {
		rm[r.name] = r
	}

	getVal := func(i int64) int64 {
		switch i {
		case 0, 1, 2, 3, 7:
			return i
		case 4:
			return rm["A"].data
		case 5:
			return rm["B"].data
		case 6:
			return rm["C"].data
		}
		return -1
	}

	setVal := func(r string, i int64) {
		reg := rm[r]
		reg.data = i
		rm[r] = reg
	}

	for {
		if prog.ptr >= len(prog.ops) {
			break
		}
		cur := prog.ops[prog.ptr]
		operand := prog.ops[prog.ptr+1]

		dojump := true
		jmp := 2
		switch cur {
		case 0:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("A", v)
			}
		case 1:
			n := rm["B"].data
			v := n ^ operand
			setVal("B", v)
		case 2:
			n := getVal(operand)
			n = n % 8
			setVal("B", n)
		case 3:
			jnz := rm["A"].data
			if jnz != 0 {
				dojump = false
				prog.ptr = int(operand)
			}
		case 4:
			b := rm["B"].data
			c := rm["C"].data
			v := b ^ c
			setVal("B", v)
		case 5:
			v := getVal(operand) % 8
			output = append(output, v)
		case 6:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("B", v)
			}
		case 7:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("C", v)
			}
		}
		if dojump {
			prog.ptr += jmp
		}

	}
	return output
}

func parse(in []string) ([]register, program) {
	rreg := regexp.MustCompile("^Register ([ABC]): ([0-9]+)$")
	preg := regexp.MustCompile("^Program: ([0-9,]+)$")

	regs := []register{}
	prog := program{}

	for _, line := range in {
		rm := rreg.FindStringSubmatch(line)
		if len(rm) != 0 {
			name := rm[1]
			val, _ := strconv.ParseInt(rm[2], 10, 64)
			reg := register{name: name, data: val}
			regs = append(regs, reg)
		}
		pm := preg.FindStringSubmatch(line)
		if len(pm) != 0 {
			ops := strings.Split(pm[1], ",")
			vals := []int64{}
			for _, s := range ops {
				v, _ := strconv.ParseInt(s, 10, 64)
				vals = append(vals, v)
			}
			prog.ptr = 0
			prog.ops = vals
		}
	}
	return regs, prog
}
