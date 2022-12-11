package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type monkey struct {
	id          int
	items       []int
	op          string
	opparam     int
	opparamself bool
	testmod     int
	truedest    int
	falsedest   int
	inspected   int
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
	monkeys := parseInput(in)
	for i := 0; i < 20; i++ {
		monkeys = runRound(monkeys, func(i int) int { return i / 3 })
	}
	insp := []int{}
	for _, m := range monkeys {
		insp = append(insp, m.inspected)
	}
	sort.Ints(insp)
	return insp[len(insp)-1] * insp[len(insp)-2]

}

func part2(in input) output {
	monkeys := parseInput(in)
	mod := 1
	for _, m := range monkeys {
		mod = lcm(mod, m.testmod)
	}
	for i := 0; i < 10000; i++ {
		monkeys = runRound(monkeys, func(i int) int { return i % mod })
	}
	insp := []int{}
	for _, m := range monkeys {
		insp = append(insp, m.inspected)
	}
	sort.Ints(insp)
	return insp[len(insp)-1] * insp[len(insp)-2]
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func runRound(monkeys []monkey, worrymetohd func(i int) int) []monkey {
	for i := 0; i < len(monkeys); i++ {
		// fmt.Println("monkey", i, "items", monkeys[i].items)
		monkeys[i].inspected += len(monkeys[i].items)
		for len(monkeys[i].items) > 0 {
			item := monkeys[i].items[0]
			monkeys[i].items = monkeys[i].items[1:]

			prm := monkeys[i].opparam
			if monkeys[i].opparamself {
				prm = item
			}

			if monkeys[i].op == "+" {
				item = item + prm
			} else {
				item = item * prm
			}

			item = worrymetohd(item)

			if item%monkeys[i].testmod == 0 {
				monkeys[monkeys[i].truedest].items = append(monkeys[monkeys[i].truedest].items, item)
				// fmt.Println("true", item, monkeys[i].testmod, "thrown to", monkeys[i].truedest, monkeys[monkeys[i].truedest].items)
			} else {
				monkeys[monkeys[i].falsedest].items = append(monkeys[monkeys[i].falsedest].items, item)
				// fmt.Println("false", item, monkeys[i].testmod, "thrown to", monkeys[i].falsedest, monkeys[monkeys[i].falsedest].items)
			}
		}
	}
	return monkeys
}

func parseInput(in input) []monkey {
	reg := regexp.MustCompile("Monkey ([0-9]+)")
	sreg := regexp.MustCompile("  Starting items: (.*)")
	opreg := regexp.MustCompile("  Operation: new = old (\\+|\\*) ([0-9]+|old)")
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
			if om[2] != "old" {
				opparam, _ := strconv.Atoi(om[2])
				monkeys[idx].opparam = opparam
			} else {
				monkeys[idx].opparamself = true
			}
			monkeys[idx].op = op
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
