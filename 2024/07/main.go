package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type problem struct {
	solution int
	operands []int
}

type op interface {
	do(a, b int) int
}

type add struct{}

func (op add) do(a, b int) int {
	return a + b
}

type mult struct{}

func (op mult) do(a, b int) int {
	return a * b
}

type concat struct{}

func (op concat) do(a, b int) int {
	ceil10 := 10
	for {
		if ceil10 > b {
			break
		}
		ceil10 *= 10
	}
	return a*ceil10 + b
}

func (op add) String() string {
	return "add"
}

func (op mult) String() string {
	return "mult"
}

func (op concat) String() string {
	return "concat"
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 07 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	list := parse(in)
	ops := []op{add{}, mult{}}

	return checkAll(list, ops, nil)
}

func part2(in input) output {
	list := parse(in)
	ops := []op{add{}, mult{}}

	return checkAll(list, ops, &concat{})
}

func checkAll(list []problem, validops []op, concatop op) int {
	total := 0
	for _, p := range list {
		if check(p, validops) {
			total += p.solution
		} else if concatop != nil && check(p, append(validops, concatop)) {
			total += p.solution
		}
	}
	return total
}

func check(p problem, validops []op) bool {
	ops := common.AllCombinations(validops, len(p.operands)-1)
	result := false
	for _, oplist := range ops {
		opnum := 0
		sol := p.operands[0]
		for i := 1; i < len(p.operands); i++ {
			op := oplist[opnum]
			sol = op.do(sol, p.operands[i])
			opnum++
		}
		if sol == p.solution {
			result = true
			break
		}
	}
	return result
}

func permOps(p problem, validops []op) [][]op {
	bits := len(p.operands) - 1
	if bits == 1 {
		res := [][]op{}
		for _, opp := range validops {
			res = append(res, []op{opp})
		}
		return res
	}
	end := int(math.Pow(2, float64(bits)))
	res := make([][]op, end)
	for i := 0; i < end; i++ {
		res[i] = []op{}
		for j := 0; j < bits; j++ {
			b := i >> j
			v := b & 1
			res[i] = append([]op{validops[v]}, res[i]...)
		}
	}

	return res
}

func parse(in []string) []problem {
	list := []problem{}
	for _, line := range in {
		sp := strings.Split(line, ":")
		sol, _ := strconv.Atoi(sp[0])

		p := problem{solution: sol}
		f := strings.Fields(sp[1])
		for _, op := range f {
			operand, _ := strconv.Atoi(op)
			p.operands = append(p.operands, operand)
		}
		list = append(list, p)
	}
	return list
}
