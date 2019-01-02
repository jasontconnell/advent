package main

import (
	"fmt"
	"strings"
	"time"
)

var input = "vzbxkghb"
var start, end byte = byte(97), byte(122)
var illegals = []rune{'i', 'o', 'l'}
var illegalBytes = []byte{byte('i'), byte('o'), byte('l')}

func main() {
	startTime := time.Now()
	valid := false
	pwd := input
	iterations := 0
	passwords := 0
	for passwords < 10000 {
		pwd = increment(pwd)
		valid = hasStraight(pwd) && twoPairs(pwd)
		iterations++
		if valid {
			fmt.Println(passwords, pwd, "after", iterations, "iterations")
			passwords++
		}
	}

	fmt.Println(time.Since(startTime))
}

func hasStraight(pwd string) bool {
	val := false
	for i := 0; i < len(pwd)-2; i++ {
		if pwd[i]+1 == pwd[i+1] && pwd[i]+2 == pwd[i+2] {
			val = true
			break
		}
	}
	return val
}

func noIllegals(pwd string) bool {
	pass := true
	for i := 0; i < len(illegals); i++ {
		pass = pass && strings.IndexByte(pwd, byte(illegals[i])) == -1
	}
	return pass
}

func twoPairs(pwd string) bool {
	pairs := 0
	for i := 0; i < len(pwd)-1; i++ {
		if pwd[i] == pwd[i+1] {
			pairs++
			i++
		}
	}
	return pairs == 2
}

func increment(pw string) string {
	return inc(pw, len(pw)-1)
}

func inc(pw string, ch int) string {
	cp := pw
	b := cp[ch] + 1

	for _, illegal := range illegalBytes {
		if b == illegal {
			b = b + 1
		}
	}

	if loop := b > end; loop {
		b = start
		cp = inc(cp, ch-1)
	}
	cp = cp[0:ch] + string(b) + cp[ch+1:]
	return cp
}
