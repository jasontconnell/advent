package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "18.txt"

type operation int

const (
	num operation = iota
	add
	mult
	group
)

type expr struct {
	parent *expr
	exprs  []*expr
	sym    string
	val    int
	op     operation
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

	probs := parseProblems(lines)
	p1 := solveAll(probs)

	grouped := changeAddPrecedence(probs)
	p2 := solveAll(grouped)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func changeAddPrecedence(ps []expr) []expr {
	grouped := []expr{}

	for _, ex := range ps {
		gex := groupAdd(ex.exprs)
		p := expr{exprs: gex}
		grouped = append(grouped, p)
	}

	return grouped
}

func groupAdd(exprs []*expr) []*expr {
	grouped := []*expr{}

	done := false
	i := 0
	for !done {
		ex := exprs[i]

		switch ex.op {
		case add:
			left := grouped[len(grouped)-1]
			right := exprs[i+1]

			if right.op == group {
				rexps := groupAdd(right.exprs)
				right = &expr{op: group, exprs: rexps}
			}

			g := []*expr{
				left,
				ex,
				right,
			}

			grouped = grouped[:len(grouped)-1]
			ng := &expr{op: group, exprs: g}
			grouped = append(grouped, ng)
			i += 2
		case mult:
			grouped = append(grouped, ex)
			i++
		case group:
			g := groupAdd(ex.exprs)

			ng := &expr{op: group, exprs: g}
			grouped = append(grouped, ng)
			i++
		case num:
			grouped = append(grouped, ex)
			i++
		}

		done = i >= len(exprs)
	}
	return grouped
}

func solveAll(ps []expr) int {
	sum := 0
	for _, p := range ps {
		sum += solveExpressions(p.exprs)
	}
	return sum
}

func solveExpressions(exprs []*expr) int {
	var i, ans int
	done := false
	for !done {
		nval, nidx := exprVal(i, ans, exprs)

		i = nidx
		ans = nval

		done = nidx == -1
	}
	return ans
}

func exprVal(i int, curval int, exprs []*expr) (nval, nidx int) {
	if i >= len(exprs) {
		return curval, -1
	}
	ex := exprs[i]

	switch ex.op {
	case num:
		nval = ex.val
		nidx = i + 1
	case add:
		nval, nidx = exprVal(i+1, curval, exprs)
		nval += curval
	case mult:
		nval, nidx = exprVal(i+1, curval, exprs)
		nval *= curval
	case group:
		nval = solveExpressions(ex.exprs)
		nidx = i + 1
	}

	return nval, nidx
}

func printExprs(exs []*expr) string {
	val := ""
	for _, ex := range exs {
		val += printExpr(*ex)
	}
	return val
}

func printExpr(ex expr) string {
	val := ""
	if ex.sym != "" {
		val += fmt.Sprint(ex.sym)
	} else if len(ex.exprs) > 0 {
		val += "( "
		for _, ec := range ex.exprs {
			val += printExpr(*ec)
		}
		val += " )"
	} else {
		val += fmt.Sprint(ex.val)
	}
	return val
}

func parseProblems(lines []string) []expr {
	probs := []expr{}
	for _, line := range lines {
		p := parseProblem(line)
		probs = append(probs, p)
	}

	return probs
}

func parseProblem(line string) expr {
	innum := false
	curnum := ""
	level := 0

	p := &expr{}
	curgroup := p

	for i, ch := range line {
		switch ch {
		case ' ':
			if innum {
				n, _ := strconv.Atoi(curnum)
				curgroup.exprs = append(curgroup.exprs, &expr{val: n, op: num})
				curnum = ""
				innum = false
			}
			break
		case '+':
			curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch), op: add})
		case '*':
			curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch), op: mult})
		case '(':
			parent := curgroup
			curgroup = &expr{parent: parent, op: group}

			parent.exprs = append(parent.exprs, curgroup)
			// curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch)})
			level++
		case ')':
			if innum {
				n, _ := strconv.Atoi(curnum)
				curgroup.exprs = append(curgroup.exprs, &expr{val: n, op: num})
				curnum = ""
				innum = false
			}
			// curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch)})
			curgroup = curgroup.parent
		default:
			if !innum {
				innum = true
			}

			curnum += string(ch)
		}

		if i == len(line)-1 && innum {
			n, _ := strconv.Atoi(curnum)
			curgroup.exprs = append(curgroup.exprs, &expr{val: n, op: num})
			curnum = ""
			innum = false
		}
	}

	return *p
}
