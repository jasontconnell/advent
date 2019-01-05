package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	//"regexp"
	"strconv"
	"strings"
	//"math"
)

var input = "23.txt"

type arguments struct {
	in, out   int
	inr, outr bool // whether in and out refer to registers
}

type instruction struct {
	opcode string
	f      func(r []int, args arguments) ([]int, int)
}

type op struct {
	inst instruction
	args arguments
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	imap := getInstructions()
	prog := []op{}
	for scanner.Scan() {
		var txt = scanner.Text()
		prog = append(prog, getOp(txt, imap))
	}

	p1 := run(prog, imap, false)
	p2 := run(prog, imap, true)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func run(prog []op, imap map[string]instruction, p2 bool) int {
	line := 0
	done := false
	r := []int{0, 0, 0, 0, 0, 0, 0, 0}
	if p2 {
		r[0] = 1
	}
	muls := 0

	for !done {
		cur := prog[line]
		upd, inc := cur.inst.f(r, cur.args)
		r = upd
		line += inc
		if cur.inst.opcode == "mul" {
			muls++
		}
		done = line < 0 || line > len(prog)-1
	}

	return muls
}

func getInstructions() map[string]instruction {
	m := make(map[string]instruction)
	m["set"] = instruction{
		opcode: "set",
		f: func(r []int, args arguments) ([]int, int) {
			v := args.out
			if args.outr {
				v = r[args.out]
			}

			r[args.in] = v
			return r, 1
		},
	}

	m["sub"] = instruction{
		opcode: "sub",
		f: func(r []int, args arguments) ([]int, int) {
			v := args.out
			if args.outr {
				v = r[args.out]
			}

			c := r[args.in] - v
			r[args.in] = c
			return r, 1
		},
	}

	m["mul"] = instruction{
		opcode: "mul",
		f: func(r []int, args arguments) ([]int, int) {
			v := args.out
			if args.outr {
				v = r[args.out]
			}

			c := r[args.in] * v
			r[args.in] = c
			return r, 1
		},
	}

	m["jnz"] = instruction{
		opcode: "jnz",
		f: func(r []int, args arguments) ([]int, int) {
			v := args.out
			if args.outr {
				v = r[args.out]
			}

			t := args.in
			if args.inr {
				t = r[args.in]
			}

			j := 1
			if t != 0 {
				j = v
			}

			return r, j
		},
	}

	return m
}

func getOp(line string, imap map[string]instruction) op {
	flds := strings.Fields(line)
	opcode := flds[0]

	args := arguments{}
	in, err := strconv.Atoi(flds[1])
	if err != nil {
		args.inr = true
		args.in = int(byte(flds[1][0]) - byte('a'))
	} else {
		args.in = in
	}

	out, err := strconv.Atoi(flds[2])
	if err != nil {
		args.outr = true
		args.out = int(byte(flds[2][0]) - byte('a'))
	} else {
		args.out = out
	}

	inst, _ := imap[opcode]
	return op{inst: inst, args: args}
}
