package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type program struct {
	input        []int
	instructions []instruction
	registers    map[string]*register
	debug        bool
}

func (p *program) in(i ...int) {
	p.input = append(p.input, i...)
}

func (p *program) run() {
	for _, inst := range p.instructions {
		m := instmap[inst.instr]
		m(inst.p1, inst.p2, p)
	}
}

func (p *program) reset() {
	for k := range p.registers {
		p.registers[k].value = 0
	}
	p.input = []int{}
}

type instruction struct {
	instr  string
	p1, p2 string
}

type register struct {
	name  string
	value int
	// regptr string
}

func regorval(a string) (reg bool, val int) {
	val, err := strconv.Atoi(a)
	reg = err != nil
	return reg, val
}

type ifunc func(a, b string, prog *program)

var instmap map[string]ifunc

func init() {
	instmap = map[string]ifunc{
		"inp": func(a, b string, prog *program) {
			if prog.debug {
				fmt.Println("read input", prog.input[0])
				printregs(prog)
			}
			prog.registers[a].value = prog.input[0]
			prog.input = prog.input[1:]
		},
		"add": func(a, b string, prog *program) {
			rb, vb := regorval(b)

			v1 := vb
			if rb {
				v1 = prog.registers[b].value
			}

			v0 := prog.registers[a].value

			prog.registers[a].value = v0 + v1
		},
		"mul": func(a, b string, prog *program) {
			rb, vb := regorval(b)
			v1 := vb
			if rb {
				v1 = prog.registers[b].value
			}

			v0 := prog.registers[a].value

			prog.registers[a].value = v0 * v1
		},
		"div": func(a, b string, prog *program) {
			rb, vb := regorval(b)
			v1 := vb
			if rb {
				v1 = prog.registers[b].value
			}

			v0 := prog.registers[a].value

			prog.registers[a].value = v0 / v1
		},
		"mod": func(a, b string, prog *program) {
			rb, vb := regorval(b)
			v1 := vb
			if rb {
				v1 = prog.registers[b].value
			}

			v0 := prog.registers[a].value

			prog.registers[a].value = v0 % v1
		},
		"eql": func(a, b string, prog *program) {
			rb, vb := regorval(b)
			v1 := vb
			if rb {
				v1 = prog.registers[b].value
			}

			v0 := prog.registers[a].value

			val := 1
			if v0 != v1 {
				val = 0
			}
			prog.registers[a].value = val
		},
	}
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
	fmt.Fprintln(w, "--2021 day 24 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	prog := parseInput(in)
	prog.debug = false
	findDigits(prog)
	printregs(prog)
	return 0
}

func part2(in input) output {
	return 0
}

func pow10(p int) int {
	return int(math.Pow10(p))
}

func arrtoint(v []int) int {
	ret := 0
	for x := len(v) - 1; x >= 0; x-- {
		p := len(v) - x - 1
		ret += v[x] * pow10(p)
	}
	return ret
}

func findDigits(p *program) int {
	maxvalid := 0
	done := false
	itr := 0

	// steps where input doesn't affect outcome are just 9
	checked := map[int]bool{}
	pool := [14][]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 1
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 2
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 3
		{9},                         // step 4
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 5
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 6
		{9},                         // step 7
		{9},                         // step 8
		{9},                         // step 9
		{9},                         // step 10
		{9},                         // step 11
		{9},                         // step 12
		{9},                         // step 13
		{1, 2, 3, 4, 5, 6, 7, 8, 9}, // step 14
	}

	for !done {
		p.reset()
		digs := []int{}
		for i := 0; i < 14; i++ {
			v := rand.Int() % len(pool[i])
			digs = append(digs, pool[i][v])
		}

		asnum := arrtoint(digs)
		if asnum < maxvalid {
			continue
		}

		if _, ok := checked[asnum]; ok {
			continue
		}
		checked[asnum] = true

		p.in(digs...)
		p.run()

		if p.registers["z"].value == 0 {
			fmt.Println("found one")
			maxvalid = asnum
			done = true
		}

		itr++
		if itr%25000 == 0 {
			fmt.Println(itr, "checked", asnum, "z is", p.registers["z"].value, len(checked))
		}
	}
	return maxvalid
}

func printregs(prog *program) {
	keys := []string{}
	for k := range prog.registers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		reg := prog.registers[k]
		fmt.Println(reg.name, reg.value)
	}
}

func parseInput(in input) *program {
	insts := []instruction{}
	regs := map[string]*register{}
	for _, line := range in {
		flds := strings.Fields(line)
		var b string
		op, a := flds[0], flds[1]
		if len(flds) == 3 {
			b = flds[2]
		}

		if _, err := strconv.Atoi(a); err != nil {
			regs[a] = &register{name: a}
		}

		inst := instruction{instr: op, p1: a, p2: b}
		insts = append(insts, inst)
	}
	prog := &program{input: []int{}, instructions: insts, registers: regs}
	return prog
}
