package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "18_test.txt"

type operation int

const (
	num operation = iota
	add
	mult
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

	fmt.Println("Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
}

func solveAll(ps []expr) int {
	sum := 0
	for _, p := range ps {
		sum += solveProblem(p)
	}
	return sum
}

func solveProblem(p expr) int {
	ans := 0
	for _, ex := range p.exprs {
		switch ex.op {
		case add:
		case mult:
		}
	}
	return ans
}

func printExpr(ex expr, end string) {
	if ex.sym != "" {
		fmt.Print(" ", ex.sym, " ")
	} else if len(ex.exprs) > 0 {
		for _, ec := range ex.exprs {
			printExpr(*ec, "")
		}
	} else {
		fmt.Print(" ", ex.val, " ")
	}
	if end != "" {
		fmt.Print(end)
	}
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

	for _, ch := range line {
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
			curgroup = &expr{}
			curgroup.parent = parent

			parent.exprs = append(parent.exprs, curgroup)
			curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch)})
			level++
		case ')':
			if innum {
				n, _ := strconv.Atoi(curnum)
				curgroup.exprs = append(curgroup.exprs, &expr{val: n, op: num})
				curnum = ""
				innum = false
			}
			curgroup.exprs = append(curgroup.exprs, &expr{sym: string(ch)})
			curgroup = curgroup.parent
		default:
			if !innum {
				innum = true
			}

			curnum += string(ch)
		}
	}
	return *p
}
