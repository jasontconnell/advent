package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type rule struct {
	left, right int
}

func (r rule) String() string {
	return fmt.Sprintf("%d must come before %d", r.left, r.right)
}

type update []int

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 05 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	rules, updates := parse(in)
	rm := getRuleMap(rules)
	return verify(updates, rm)
}

func part2(in input) output {
	rules, updates := parse(in)
	rm := getRuleMap(rules)
	return fix(updates, rm)
}

func fix(updates []update, rules map[int][]int) int {
	fixed := []update{}
	for i := 0; i < len(updates); i++ {
		if !verifyUpdate(updates[i], rules) {
			fixed = append(fixed, fixUpdate(updates[i], rules))
		}
	}
	return verify(fixed, rules)
}

func fixUpdate(upd update, rules map[int][]int) update {
	nupdate := update{}
	for _, n := range upd {
		if len(nupdate) == 0 {
			nupdate = append(nupdate, n)
			continue
		}

		idx := getInsertIndex(nupdate, rules, n)
		if idx == -1 {
			nupdate = append([]int{n}, nupdate...)
		} else if idx >= len(nupdate) {
			nupdate = append(nupdate, n)
		} else {
			nupdate = append(nupdate[:idx], append([]int{n}, nupdate[idx:]...)...)
		}

	}
	return nupdate
}

func getInsertIndex(upd update, rules map[int][]int, n int) int {
	i := len(upd)
	if i == 0 {
		return 0
	}

	rn := rules[n]
	rnm := make(map[int]int)
	for _, r := range rn {
		rnm[r] = r
	}
	for idx, x := range upd {
		if _, ok := rnm[x]; ok {
			i = idx
			break
		}
	}

	return i
}

func verify(updates []update, rules map[int][]int) int {
	summids := 0
	for _, upd := range updates {
		if verifyUpdate(upd, rules) {
			mid := upd[len(upd)/2]
			summids += mid
		}
	}
	return summids
}

func verifyUpdate(upd update, rules map[int][]int) bool {
	correct := true
	for i, x := range upd {
		r := rules[x]
		if !verifyN(i, upd, r) {
			correct = false
			break
		}
	}
	return correct
}

func verifyN(idx int, upd update, after []int) bool {
	for _, x := range after {
		xidx := slices.Index(upd, x)
		if xidx != -1 && xidx < idx {
			return false
		}
	}
	return true
}

func getRuleMap(rules []rule) map[int][]int {
	m := make(map[int][]int)
	for _, r := range rules {
		m[r.left] = append(m[r.left], r.right)
	}
	return m
}

func parse(in []string) ([]rule, []update) {
	rules := []rule{}
	updates := []update{}
	var i int
	for {
		if in[i] == "" {
			break
		}
		sp := strings.Split(in[i], "|")
		if len(sp) != 2 {
			break
		}
		a, _ := strconv.Atoi(sp[0])
		b, _ := strconv.Atoi(sp[1])
		rules = append(rules, rule{a, b})
		i++
	}

	for {
		if i >= len(in) {
			break
		}

		sp := strings.Split(in[i], ",")
		upd := update{}
		for _, c := range sp {
			a, _ := strconv.Atoi(c)
			upd = append(upd, a)
		}
		if len(upd) > 0 {
			updates = append(updates, upd)
		}
		i++
	}
	return rules, updates
}
