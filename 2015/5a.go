package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var input = "5.txt"

func main() {
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		goodCount := 0
		allCount := 0

		for scanner.Scan() {
			var txt = scanner.Text()

			if pairCheck(txt) && wrapCheck(txt) {
				goodCount++
			}

			allCount++
		}

		fmt.Println("Good:", goodCount, " Total:", allCount)
	}
}

func pairCheck(s string) bool {
	pairs := 0
	for i, _ := range s {
		if i > 0 {
			overlaps := false
			current := string(s[i-1]) + string(s[i])
			c := strings.Count(s[i:], current)
			if s[i-1] == s[i] && i < len(s)-1 {
				if s[i] == s[i+1] {
					overlaps = true
				}
			}

			if c > 0 && !overlaps {
				pairs++
			}
		}
	}
	return pairs > 0
}

func wrapCheck(s string) bool {
	wrap := false
	for i, _ := range s {
		if i > 1 && s[i-2] == s[i] {
			wrap = true
			break
		}
	}
	return wrap
}
