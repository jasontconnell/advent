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

type op string

const (
	nul op = ""
	add op = "+"
	sub op = "-"
	mul op = "*"
	div op = "/"
	eql op = "=="
)

type monkey struct {
	id string
	l  target
	op op
	r  *target
}

type target struct {
	id    string
	value int
}

func (m monkey) String() string {
	return fmt.Sprintf("%s %v %v %v", m.id, m.l, m.op, m.r)
}

func (t *target) String() string {
	return fmt.Sprintf("[%s %d]", t.id, t.value)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 21 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	monkeys := parseInput(in)
	val, _, _ := findRoot(monkeys, "root")
	return val
}

func part2(in input) output {
	monkeys := parseInput(in)
	return findHuman(monkeys, "root", "humn")
}

func findHuman(monkeys []monkey, find, human string) int {
	hidx := 0
	for i, mm := range monkeys {
		if mm.id == find {
			monkeys[i].op = eql
		}
		if mm.id == human {
			hidx = i
			monkeys[i].l.value = 3952288690728
		}
	}

	fval, test, success := findRoot(monkeys, find)
	fmt.Println(fval - test)
	var mult float64
	val := 1
	for !success {
		mult = (float64(fval) / float64(test))
		val = int(float64(val) * mult)

		monkeys[hidx].l.value = val
		fval, test, success = findRoot(monkeys, find)
	}
	return val
}

func findRoot(monkeys []monkey, find string) (int, int, bool) {
	yell := make(map[string]int)
	for _, mm := range monkeys {
		if mm.l.id == "" {
			yell[mm.id] = mm.l.value
		}
	}

	_, found := yell[find]
	var equal bool
	var result int
	for !found {
		for _, m := range monkeys {
			if _, ok := yell[m.id]; ok {
				continue
			}
			lv, rv := 0, 0
			h := 0
			if y, ok := yell[m.l.id]; ok {
				lv = y
				h++
			}

			if m.r != nil {
				if y, ok := yell[(*m.r).id]; ok {
					rv = y
					h++
				}
			}

			if h == 2 {
				switch m.op {
				case add:
					yell[m.id] = lv + rv
				case sub:
					yell[m.id] = lv - rv
				case mul:
					yell[m.id] = lv * rv
				case div:
					yell[m.id] = lv / rv
				case eql:
					result = rv
					yell[m.id] = lv
					equal = lv == rv
				}
			}
		}
		_, found = yell[find]
	}

	return yell[find], result, equal
}

func parseInput(in input) []monkey {
	list := []monkey{}
	for _, line := range in {
		sp := strings.Split(line, " ")
		if len(sp) == 0 {
			continue
		}

		mnk := monkey{id: sp[0][:len(sp[0])-1]}
		mnk.l = *getTarget(sp[1])
		if len(sp) > 2 {
			mnk.op = op(sp[2])
			mnk.r = getTarget(sp[3])
		}
		list = append(list, mnk)
	}
	return list
}

func getTarget(str string) *target {
	t := &target{}
	n, err := strconv.Atoi(str)
	if err == nil {
		t.value = n
	} else {
		t.id = str
	}
	return t
}
