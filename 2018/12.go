package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type rule struct {
	in  []bool
	out bool
}

var input = "12.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}

	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	pad := 150
	initial := parseInit(lines[0], pad)
	rules := parseRules(lines[2:])

	//part 1
	_, result := sim(initial, rules, 20, pad)
	sum := sumIndices(result, pad)
	fmt.Println("part 1:", sum)

	//part 2 . make more room in initial
	initial = append(initial, make([]bool, 160)...)
	dupTurn, result2 := sim(initial, rules, 200, pad)
	idcs := indices(result2, pad)
	sum2 := sumIndices(result2, pad)

	// remaining turns: 50 billion minus the turn that the duplicates start
	turnsLeft := int64(50000000000) - int64(dupTurn)
	count := len(idcs)

	fmt.Println("part 2:", sum2, int64(sum2)+turnsLeft*int64(count))

	fmt.Println("Time", time.Since(startTime))
}

func indices(result []bool, pad int) []int {
	idcs := []int{}

	for i, v := range result {
		if v {
			idcs = append(idcs, i-pad)
		}
	}
	return idcs
}

func sumIndices(result []bool, pad int) int {
	sum := 0
	for _, i := range indices(result, pad) {
		sum += i
	}
	return sum
}

func first(val []bool) int {
	first := 0
	for i := 0; i < len(val); i++ {
		if val[i] && first == 0 {
			first = i
			break
		}
	}
	return first
}

func last(val []bool) int {
	last := 0
	for i := len(val) - 1; i >= 0; i-- {
		if val[i] && last == 0 {
			last = i
			break
		}
	}
	return last
}

func tostr(val []bool) string {
	str := ""
	fst, lst := first(val), last(val)
	for i := fst; i < lst+1; i++ {
		v := val[i]
		c := "."
		if v {
			c = "#"
		}
		str += c
	}
	return str
}

func sim(initial []bool, rules []rule, turns, pad int) (int, []bool) {
	m := make(map[string]bool)
	gen := initial
	result := 0
	m[tostr(gen)] = true
	for i := 1; i < turns; i++ {
		gen = processMatches(rules, gen)
		key := tostr(gen)
		fmt.Println(indices(gen, pad))
		if _, ok := m[key]; ok {
			result = i
			break
		}
		m[key] = true
	}
	return result, gen
}

func processMatches(rules []rule, val []bool) []bool {
	result := make([]bool, len(val))
	for i := 0; i < len(val)-5; i++ {
		chunk := val[i : i+5]
		for _, r := range rules {
			m := match(r, chunk)
			if m {
				result[i+2] = r.out
			}
		}
	}
	return result
}

func match(r rule, val []bool) bool {
	for i, b := range r.in {
		if val[i] != b {
			return false
		}
	}
	return true
}

func parseInit(init string, pad int) []bool {
	s := len("initial state: ")
	grid := parse(string(init[s:]))
	// pad out grid on both sides
	for i := 0; i < pad; i++ {
		grid = append([]bool{false}, grid...)
		grid = append(grid, false)
	}
	return grid
}

func parseRules(r []string) []rule {
	rules := []rule{}
	for _, s := range r {
		parts := strings.Split(s, " => ")
		rl := rule{in: parse(parts[0]), out: getVal(rune(parts[1][0]))}

		rules = append(rules, rl)
	}
	return rules
}

func parse(s string) []bool {
	b := make([]bool, len(s))
	for i, c := range s {
		b[i] = getVal(c)
	}
	return b
}

func getVal(r rune) bool {
	return r == '#'
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
