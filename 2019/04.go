package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "04.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	var lbound, ubound int

	for scanner.Scan() {
		var txt = scanner.Text()
		parts := strings.Split(txt, "-")

		if len(parts) != 2 {
			panic("invalid input")
		}

		var err error

		lbound, err = strconv.Atoi(parts[0])
		ubound, err = strconv.Atoi(parts[1])

		if err != nil {
			panic("error parsing" + txt)
		}
	}

	p1 := part1(lbound, ubound)
	fmt.Println("Part 1: ", p1)

	p2 := part2(lbound, ubound)
	fmt.Println("Part 2: ", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(lbound, ubound int) int {
	total := 0
	for i := lbound; i < ubound; i++ {
		v := isValid(i)

		if v {
			total++
		}
	}
	return total
}

func part2(lbound, ubound int) int {
	total := 0
	for i := lbound; i < ubound; i++ {
		v := isValid2(i)

		if v {
			total++
		}
	}
	return total
}

func increasing(d []int) bool {
	increasing := true
	var l int
	for _, i := range d {
		increasing = increasing && i >= l
		l = i
	}
	return increasing
}

func isValid(val int) bool {
	d := digits(val)
	r := runes(d)

	inc := increasing(d)
	if !inc {
		return false
	}

	var lr rune
	hasdbl := false
	for _, c := range r {
		if lr == c {
			hasdbl = true
			break
		}
		lr = c
	}

	return hasdbl
}

func isValid2(val int) bool {
	d := digits(val)
	r := runes(d)

	inc := increasing(d)
	if !inc {
		return false
	}

	var lr rune
	count := 1
	validrpt := true
	rpts := false

	for n, c := range r {
		if lr == c {
			count++
		} else if count > 1 {
			rpts = count == 2
			if rpts {
				validrpt = true
				break
			}
			count = 1
		}
		lr = c

		if n == len(r)-1 {
			rpts = count == 2
		}
	}
	return rpts && validrpt
}

func digits(val int) []int {
	v := []int{}
	c := val
	div := 10
	done := false

	var x int
	for !done {
		x = c % div
		v = append(v, x)
		c = c / div

		done = c == 0
	}

	// reverse slice
	for i := 0; i < len(v)/2; i++ {
		v[i], v[len(v)-i-1] = v[len(v)-i-1], v[i]
	}

	return v
}

func runes(ds []int) []rune {
	r := []rune{}
	for _, d := range ds {
		r = append(r, rune(strconv.Itoa(d)[0]))
	}

	return r
}
