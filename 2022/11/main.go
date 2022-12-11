package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type monkey struct {
	id        int
	items     []int
	op        string
	opparam   int
	testmod   int
	truedest  int
	falsedest int
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
	fmt.Fprintln(w, "--2022 day 11 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	parseInput(in)
	return 0
}

func part2(in input) output {
	return 0
}

func parseInput(in input) []monkey {
	reg := regexp.MustCompile("Monkey ([0-9]+)")
	sreg := regexp.MustCompile("  Starting items: (.*)")
	opreg := regexp.MustCompile("  Operation: new = old (\\+|\\*) ([0-9]+)")
	tstreg := regexp.MustCompile("  Test: divisible by ([0-9]+)")
	casereg := regexp.MustCompile("    If (true|false): throw to monkey ([0-9]+)")

	monkeys := []monkey{}
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) == 2 {
			id, _ := strconv.Atoi(m[1])
			mk := monkey{id: id}
			monkeys = append(monkeys, mk)
			continue
		}

		idx := len(monkeys) - 1
		stm := sreg.FindStringSubmatch(line)
		if len(stm) == 2 {
			items := strings.Split(stm[1], ",")
			for _, s := range items {
				sn, _ := strconv.Atoi(strings.Trim(s, " "))
				monkeys[idx].items = append(monkeys[idx].items, sn)
			}
			continue
		}

		om := opreg.FindStringSubmatch(line)
		if len(om) == 3 {
			op := om[1]
			opparam, _ := strconv.Atoi(om[2])
			monkeys[idx].op = op
			monkeys[idx].opparam = opparam
			continue
		}

		tstm := tstreg.FindStringSubmatch(line)
		if len(tstm) == 2 {
			t, _ := strconv.Atoi(tstm[1])
			monkeys[idx].testmod = t
			continue
		}

		cm := casereg.FindStringSubmatch(line)
		if len(cm) == 3 {
			dest, _ := strconv.Atoi(cm[2])
			if cm[1] == "true" {
				monkeys[idx].truedest = dest
			} else {
				monkeys[idx].falsedest = dest
			}
		}
	}
	return monkeys
}
