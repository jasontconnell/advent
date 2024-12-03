package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type mul struct {
	a, b int
	cond bool
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 03 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	muls := parse(strings.Join(in, ""))
	x := 0
	for _, m := range muls {
		x += (m.a * m.b)
	}
	return x
}

func part2(in input) output {
	muls := parse(strings.Join(in, ""))
	x := 0
	for _, m := range muls {
		if m.cond {
			x += (m.a * m.b)
		}
	}
	return x
}

func parse(in string) []mul {
	m := []mul{}
	fname := ""
	cur := ""
	infunc := false
	inmul := false
	do := true
	var a, b int
	isb := false
	for i, c := range in {
		switch c {
		case 'm', 'u', 'l':
			if c == 'm' {
				fname = ""
			}
			fname += string(c)
			inmul = fname == "mul"
			infunc = false
		case 'd', 'o', 'n', '\'', 't':
			if c == 'd' {
				fname = ""
			}
			fname += string(c)
			peek := in[i+1 : i+3]
			if (fname == "do" || fname == "don't") && peek == "()" {
				do = fname == "do"
				if fname == "don't" {
					do = false
				}
				fname = ""
			}
			inmul = false
			infunc = false
			cur = ""
			a = 0
			b = 0
		case '(':
			if inmul {
				infunc = true
			}
			fname = ""
			cur = ""
			a = 0
			b = 0
		case ')':
			if infunc && isb {
				b, _ = strconv.Atoi(cur)
				cmul := mul{a, b, do}
				m = append(m, cmul)
			}
			cur = ""
			infunc = false
			isb = false
			inmul = false
			a = 0
			b = 0
		case ',':
			if inmul && infunc {
				a, _ = strconv.Atoi(cur)
				isb = true
			} else {
				inmul = false
				infunc = false
			}
			cur = ""
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			cur += string(c)
		default:
			inmul = false
			infunc = false
			cur = ""
			fname = ""
			a = 0
			b = 0
		}
	}
	return m
}
