package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var input = "16.txt"

type rangeRule struct {
	min, max int
}

type fieldRule struct {
	field  string
	ranges []rangeRule
}

type ticket []int

type ticketTranslation struct {
	fieldRules []fieldRule
	mine       ticket
	nearby     []ticket
}

type fieldValue struct {
	field string
	value int
}
type translatedTicket struct {
	fields []fieldValue
}

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

	tt := getTicketTranslation(lines)
	invalid := findInvalidTickets(tt)

	fmt.Println("Part 1:", invalid)

	found := findTicket(tt)

	p2 := 1
	for _, f := range found.fields {
		if strings.HasPrefix(f.field, "departure") {
			p2 *= f.value
		}
	}

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func findInvalidTickets(tt ticketTranslation) int {
	invalid := 0
	for _, t := range tt.nearby {
		for _, v := range t {
			invalid += getErrorRate(v, tt.fieldRules)
		}
	}
	return invalid
}

func getErrorRate(val int, fieldRules []fieldRule) int {
	rate := val
	insideAny := false
	for _, fr := range fieldRules {
		for _, rr := range fr.ranges {
			if val >= rr.min && val <= rr.max {
				insideAny = true
			}
		}
	}
	if insideAny {
		rate = 0
	}
	return rate
}

func findTicket(tt ticketTranslation) translatedTicket {
	vt := getValidTickets(tt)

	translated := translatedTicket{}
	determined := make(map[string]string)
	undetermined := make(map[int]int)
	for i, v := range tt.mine {
		undetermined[i] = v
	}

	for len(undetermined) > 0 {
		for i, v := range undetermined {
			n := findField(i, vt, tt.fieldRules, determined)
			determined[n] = n
			if n != "" {
				translated.fields = append(translated.fields, fieldValue{field: n, value: v})
				delete(undetermined, i)
			}
		}
	}

	return translated
}

func findField(index int, tickets []ticket, fieldRules []fieldRule, determined map[string]string) string {
	m := make(map[string]fieldRule)

	for _, fr := range fieldRules {
		_, ok := determined[fr.field]

		if !ok {
			m[fr.field] = fr
		}
	}

	for _, t := range tickets {
		for name, fr := range m {
			if !couldBeField(t[index], fr) {
				delete(m, name)
			}
		}
	}

	field := ""
	if len(m) == 1 {
		for k, _ := range m {
			field = k
		}
	}
	return field
}

func couldBeField(val int, rule fieldRule) bool {
	could := false
	for _, f := range rule.ranges {
		if val >= f.min && val <= f.max {
			could = true
			break
		}
	}
	return could
}

func getValidTickets(tt ticketTranslation) []ticket {
	valid := []ticket{}
	for _, t := range tt.nearby {
		if isTicketValid(tt, t) {
			valid = append(valid, t)
		}
	}
	return valid
}

func isTicketValid(tt ticketTranslation, t ticket) bool {
	valid := true
	for _, v := range t {
		rate := getErrorRate(v, tt.fieldRules)
		if rate != 0 {
			valid = false
		}
	}
	return valid
}

var frreg *regexp.Regexp = regexp.MustCompile("^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$")

func getTicketTranslation(lines []string) ticketTranslation {
	index := 0

	tt := ticketTranslation{}

	rules := []fieldRule{}
	for {
		line := lines[index]
		if line == "" {
			break
		}

		g := frreg.FindStringSubmatch(line)

		lrs, _ := strconv.Atoi(g[2])
		lre, _ := strconv.Atoi(g[3])

		hrs, _ := strconv.Atoi(g[4])
		hre, _ := strconv.Atoi(g[5])

		field := g[1]

		ranges := []rangeRule{
			{lrs, lre},
			{hrs, hre},
		}

		r := fieldRule{field: field, ranges: ranges}
		rules = append(rules, r)

		index++
	}

	for lines[index] != "your ticket:" {
		index++
	}

	index++

	mine := readTicket(lines[index])
	for lines[index] != "nearby tickets:" {
		index++
	}
	index++

	nearby := []ticket{}
	for index < len(lines) {
		nearby = append(nearby, readTicket(lines[index]))
		index++
	}

	tt.fieldRules = rules
	tt.mine = mine
	tt.nearby = nearby

	return tt
}

func readTicket(line string) ticket {
	t := ticket{}
	ss := strings.Split(line, ",")
	for _, s := range ss {
		n, _ := strconv.Atoi(s)
		t = append(t, n)
	}
	return t
}
