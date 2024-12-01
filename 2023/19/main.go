package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type op int

const (
	none op = iota
	lt
	gt
)

func (o op) String() string {
	s := ">"
	if o == lt {
		s = "<"
	}
	return s
}

type workflow struct {
	name string
	rule *rule
}

type part struct {
	x, m, a, s int
}

func (p part) String() string {
	return fmt.Sprintf("x:%d m:%d a:%d s:%d", p.x, p.m, p.a, p.s)
}

type rule struct {
	jump   string
	lparam string
	rval   *int

	op     op
	destwf string

	accept bool
	reject bool

	passrule *rule
	failrule *rule
	parent   *rule
}

type result struct {
	accept, reject bool
	jump           string
}

type exresult struct {
	presults []presult
	accept   bool
	reject   bool
}

type presult struct {
	param    string
	min, max *int
}

func (r *rule) String() string {
	var ops string
	if r.op != none {
		ops = "<"
		if r.op == gt {
			ops = ">"
		}
	}

	left := r.lparam
	right := ""
	if r.rval != nil {
		right = fmt.Sprintf("%d", *r.rval)
	}

	if r.accept {
		return "A"
	}

	if r.reject {
		return "R"
	}

	if r.jump != "" {
		return r.jump
	}

	sub := ""
	if r.passrule != nil {
		sub += ":" + r.passrule.String()
		if r.failrule != nil {
			sub += "," + r.failrule.String()
		}
	}

	return fmt.Sprintf("%s%s%s%s", left, ops, right, sub)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	wflist, plist := parseInput(in)
	return evalWorkflow(wflist, plist, "in")
}

func part2(in input) int64 {
	wflist, _ := parseInput(in)
	return examine(wflist, "in", 1, 4000)
}

func examine(wflist []workflow, startwf string, min, max int) int64 {
	wfm := make(map[string]workflow)
	for _, w := range wflist {
		wfm[w.name] = w
	}
	start := wfm[startwf]
	mins := map[string]int{
		"x": max, "m": max, "a": max, "s": max,
	}
	maxes := map[string]int{
		"x": min, "m": min, "a": min, "s": min,
	}

	exresults := examineRule(start.rule, wfm)
	fmt.Println("in examine", exresults)
	for _, res := range exresults.presults {
		if res.min != nil && *res.min < mins[res.param] {
			mins[res.param] = *res.min
		} else if res.max != nil && *res.max > maxes[res.param] {
			maxes[res.param] = *res.max
		}
	}
	var total int64 = 1
	for k := range mins {
		if maxes[k] == min {
			maxes[k] = max
		}
		total *= int64(maxes[k] - mins[k])
	}
	fmt.Println(mins)
	fmt.Println(maxes)
	fmt.Println(167409079868000)
	fmt.Println(total)
	fmt.Println(4000 * 4000 * 4000 * 4000)
	return total
}

func abs(i int) int {
	if i < 0 {
		i = i * -1
	}
	return i
}

func examineRule(r *rule, wfm map[string]workflow) exresult {
	// min, max := 0, 0
	// if r.lparam != "" && r.rval != nil {
	// 	switch r.op {
	// 	case lt:
	// 		max = *r.rval
	// 	case gt:
	// 		min = *r.rval
	// 	}
	// }

	if r.accept {
		return exresult{accept: true}
	} else if r.reject {
		return exresult{reject: true}
	}

	presults := []presult{}
	if r.passrule != nil && r.failrule != nil {
		fmt.Println("recurse pass or fail")
		for _, porf := range []*rule{r.passrule, r.failrule} {
			res := examineRule(porf, wfm)
			if res.accept {
				if r.op == lt {
					pres := presult{param: r.lparam, max: r.rval}
					presults = append(presults, pres)
				} else {
					pres := presult{param: r.lparam, min: r.rval}
					presults = append(presults, pres)
				}
			}
			// else if res.reject {
			// 	if r.op == lt {
			// 		pres := presult{param: r.lparam, min: r.rval}
			// 		presults = append(presults, pres)
			// 	} else {
			// 		pres := presult{param: r.lparam, max: r.rval}
			// 		presults = append(presults, pres)
			// 	}
			// }
			presults = append(presults, res.presults...)
		}
	}

	// if r.failrule != nil {
	// 	fmt.Println("recurse fail")
	// 	examineRule(r.failrule, wfm)
	// }

	if r.jump != "" {
		wf := wfm[r.jump]
		res := examineRule(wf.rule, wfm)
		fmt.Println("from jump", res)
		presults = append(presults, res.presults...)
	}

	// fmt.Println(min, max)
	return exresult{presults: presults}
}

func evalWorkflow(wflist []workflow, plist []part, startwf string) int {
	wfm := make(map[string]workflow)
	for _, w := range wflist {
		wfm[w.name] = w
	}

	total := 0
	start := wfm[startwf]
	for _, p := range plist {
		cur := start
		done := false
		for !done {
			res := evalRule(p, cur.rule)
			if res.accept {
				total += p.x + p.m + p.a + p.s
			} else if res.jump != "" {
				cur = wfm[res.jump]
			}

			done = res.reject || res.accept
		}
	}
	return total
}

func evalRule(p part, r *rule) result {
	if r.accept {
		return result{accept: true}
	}
	if r.jump != "" {
		return result{jump: r.jump}
	}
	if r.reject {
		return result{reject: true}
	}
	cval := 0
	switch r.lparam {
	case "x":
		cval = p.x
	case "m":
		cval = p.m
	case "a":
		cval = p.a
	case "s":
		cval = p.s
	}

	if r.rval == nil {
		log.Fatal("comparing to nil value ", r.lparam)
	}

	var cres bool
	switch r.op {
	case lt:
		cres = cval < *r.rval
	case gt:
		cres = cval > *r.rval
	}

	if cres {
		return evalRule(p, r.passrule)
	} else {
		return evalRule(p, r.failrule)
	}
}

func parseInput(in input) ([]workflow, []part) {
	reg := regexp.MustCompile(`([a-z]+)\{(.*?)\}`)
	preg := regexp.MustCompile(`\{x=([\-0-9]+),m=([\-0-9]+),a=([\-0-9]+),s=([\-0-9]+)\}`)
	wfs := []workflow{}
	plist := []part{}
	parts := false
	for _, line := range in {
		if line == "" {
			parts = true
			continue
		}
		if !parts {
			m := reg.FindStringSubmatch(line)
			nm := m[1]

			r := parseRule(m[2])
			wf := workflow{name: nm, rule: r}
			wfs = append(wfs, wf)
		} else {
			xmas := preg.FindStringSubmatch(line)
			x, _ := strconv.Atoi(xmas[1])
			m, _ := strconv.Atoi(xmas[2])
			a, _ := strconv.Atoi(xmas[3])
			s, _ := strconv.Atoi(xmas[4])

			p := part{x: x, m: m, a: a, s: s}
			plist = append(plist, p)
		}
	}
	return wfs, plist
}

func parseRule(s string) *rule {
	cur := &rule{parent: nil}
	prule := cur
	ns := ""
	cs := ""
	instr := false
	for i, c := range s {
		switch c {
		case '<', '>':
			instr = false
			if len(cs) > 0 {
				cur.lparam = cs
				cs = ""
			}
			cur.op = lt
			if c == '>' {
				cur.op = gt
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			instr = false
			ns += string(c)
		case 'A':
			instr = false
			cur.accept = true
		case 'R':
			instr = false
			cur.reject = true
		case ':':
			instr = false
			if len(ns) > 0 {
				n, _ := strconv.Atoi(ns)
				ns = ""
				cur.rval = &n
			}
			if len(cs) > 0 {
				cur.jump = cs
				cs = ""
			}
			cur.passrule = &rule{parent: cur}
			cur = cur.passrule
		case ',':
			if instr {
				cur.jump = cs
				cs = ""
			}
			instr = false
			cur.parent.failrule = &rule{parent: cur.parent}
			cur = cur.parent.failrule
		default:
			instr = true
			cs += string(c)
		}

		if i == len(s)-1 && instr {
			cur.jump = cs
		}
	}

	return prule
}
