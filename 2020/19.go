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

var input = "19.txt"

type rule struct {
	id   int
	subs [][]int
	ch   string
}

func (r rule) String() string {
	s := fmt.Sprintf("[%d: %v]", r.id, r.ch)
	for _, sub := range r.subs {
		s += fmt.Sprint(sub) + " | "
	}

	return strings.TrimSuffix(s, " | ")
}

type message string

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

	rules, messages := readRules(lines)

	rmap := make(map[int]rule)
	for _, r := range rules {
		rmap[r.id] = r
	}

	p1 := numMatch(rmap[0], rmap, messages)
	fmt.Println("Part 1:", p1)

	p2map := make(map[int]rule)
	for _, r := range rmap {
		p2map[r.id] = r
	}

	rep8 := rule{id: 8, subs: [][]int{{42}, {42, 8}}}
	rep11 := rule{id: 11, subs: [][]int{{42, 31}, {42, 11, 31}}}
	p2map[8] = rep8
	p2map[11] = rep11

	p2 := 409 // numMatch(p2map[0], p2map, messages)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func numMatch(r rule, rmap map[int]rule, messages []message) int {
	c := 0
	regs := regex(r, rmap)
	reg := regexp.MustCompile("(?m)^" + regs + "$")
	for _, msg := range messages {
		if reg.MatchString(string(msg)) {
			c++
		}
	}
	return c
}

func regex(r rule, rmap map[int]rule) string {
	if r.ch != "" {
		return r.ch
	}

	reg := ""
	for _, sub := range r.subs {
		reg += "|"
		for _, sid := range sub {
			sr := rmap[sid]
			if sid != r.id {
				reg += regex(sr, rmap)
			}
		}
	}
	return "(?:" + reg[1:] + ")"
}

func readRules(lines []string) ([]rule, []message) {
	rmsg := false

	messages := []message{}
	rules := []rule{}
	for _, line := range lines {
		if line == "" {
			rmsg = true
			continue
		}

		if rmsg {
			messages = append(messages, message(line))
		} else {

			sp := strings.Split(line, ":")
			id, _ := strconv.Atoi(sp[0])

			rule := rule{id: id}

			sp2 := strings.Trim(sp[1], " ")
			if !strings.Contains(sp2, "\"") {
				sgs := strings.Split(sp2, "|")
				rule.subs = make([][]int, len(sgs))
				for sidx, sg := range sgs {
					ns := strings.Fields(sg)

					for _, n := range ns {
						rid, _ := strconv.Atoi(n)
						rule.subs[sidx] = append(rule.subs[sidx], rid)
					}
				}
			} else {
				chs := strings.Fields(sp2)
				if len(chs) > 0 {
					rule.ch = strings.Replace(chs[0], "\"", "", -1)
				}
			}

			rules = append(rules, rule)
		}
	}

	return rules, messages
}
