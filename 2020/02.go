package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

var input = "02.txt"

type pw struct {
	Min, Max int
	Rune     rune
	Value    string
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pwds := []pw{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var txt = scanner.Text()
		pwds = append(pwds, getPassword(txt))
	}

	valid := getValid(pwds)
	fmt.Println("Valid count: ", valid)

	valid2 := getValidPart2(pwds)
	fmt.Println("Valid count (p2): ", valid2)

	fmt.Println("Time", time.Since(startTime))
}

func getValid(pwds []pw) int {
	validCount := 0
	for _, pwd := range pwds {
		cnt := 0
		for _, x := range pwd.Value {
			if x == pwd.Rune {
				cnt++
			}
		}

		valid := cnt >= pwd.Min && cnt <= pwd.Max
		if valid {
			validCount++
		}
	}

	return validCount
}

func getValidPart2(pwds []pw) int {
	validCount := 0
	for _, pwd := range pwds {
		pos1, pos2 := rune(pwd.Value[pwd.Min-1]), rune(pwd.Value[pwd.Max-1])
		if pos1 == pos2 {
			continue
		}

		valid := pos1 == pwd.Rune || pos2 == pwd.Rune
		if valid {
			validCount++
		}
	}
	return validCount
}

var reg *regexp.Regexp = regexp.MustCompile("([0-9]+)-([0-9]+) ([a-z]): (.*)")

func getPassword(line string) pw {
	pwd := pw{}
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		min, _ := strconv.Atoi(groups[1])
		max, _ := strconv.Atoi(groups[2])

		ch := rune(groups[3][0])
		word := groups[4]

		pwd.Min = min
		pwd.Max = max
		pwd.Rune = ch
		pwd.Value = word
	}

	return pwd
}

/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
