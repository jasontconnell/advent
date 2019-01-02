package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var input = "5.txt"

var bad []string = []string{"ab", "cd", "pq", "xy"}
var vowels = []rune{'a', 'e', 'i', 'o', 'u'}

func main() {
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		goodCount := 0

		for scanner.Scan() {
			var txt = scanner.Text()

			if !badCheck(txt, bad) && appearanceCheck(txt, vowels, 3) && doubleCheck(txt) {
				goodCount++
			}
		}

		fmt.Println("Good:", goodCount)
	}
}

func appearanceCheck(s string, runes []rune, count int) bool {
	total := 0
	for _, r := range runes {
		total += strings.Count(s, string(r))
	}
	return total >= count
}

func doubleCheck(s string) bool {
	double := false

	for i, _ := range s {
		if i > 0 && s[i-1] == s[i] {
			double = true
			break
		}
	}
	return double
}

func badCheck(s string, bad []string) bool {
	isBad := false
	for _, b := range bad {
		if !isBad {
			isBad = strings.Contains(s, b)
		}
	}
	return isBad
}
