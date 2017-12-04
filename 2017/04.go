package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	//"strconv"
	"sort"
	"strings"
	//"math"
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

	vc := 0
	vc2 := 0
	for scanner.Scan() {
		var txt = scanner.Text()
		if isValid(txt) {
			vc++
		}

		if isValidPart2(txt) {
			vc2++
		}
	}

	fmt.Println("valid count", vc)
	fmt.Println("valid count part 2", vc2)
	fmt.Println("Time", time.Since(startTime))
}

func isValid(line string) bool {
	sp := strings.Split(line, " ")
	sort.Strings(sp)

	valid := true
	for i, s := range sp {
		if i > 0 && s == sp[i-1] {
			valid = false
			break
		}
	}
	return valid
}

func isValidPart2(line string) bool {
	sp := strings.Split(line, " ")
	strs := []string{}

	for i := 0; i < len(sp); i++ {
		runes := make([]int, len(sp[i]))
		for _, c := range sp[i] {
			if c > 96 && c < 123 || c > 64 && c < 91 {
				runes = append(runes, int(c))
			}
		}

		sort.Ints(runes)
		str := ""
		for _, c := range runes {
			str += string(c)
		}

		tstr := strings.TrimSpace(str)
		strs = append(strs, tstr)
	}

	sort.Strings(strs)

	valid := true
	for i, _ := range strs {
		if i > 0 && strs[i] == strs[i-1] {
			valid = false
			break
		}
	}
	return valid
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	eq := true
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			eq = false
			break
		}
	}

	return eq
}
