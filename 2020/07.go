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

var input = "07.txt"

type bagcount struct {
	bag   string
	count int
}

type bag struct {
	bag      string
	contains []bagcount
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	bags := []bag{}
	for scanner.Scan() {
		var txt = scanner.Text()
		bags = append(bags, getBagRule(txt))
	}

	bm := make(map[string]bag)
	for _, b := range bags {
		bm[b.bag] = b
	}

	p1 := 0
	for _, b := range bags {
		can := canContainBag(b, "shiny gold", bm)
		if can {
			p1++
		}
	}

	p2 := getContainedBagCount("shiny gold", bm)

	fmt.Println("part 1:", p1)
	fmt.Println("part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

var reg *regexp.Regexp = regexp.MustCompile("^([a-z ]*?) bags contain (.*)$")
var subreg *regexp.Regexp = regexp.MustCompile("([0-9]+) ([a-z ]*?) bag")

func getBagRule(line string) bag {
	groups := reg.FindStringSubmatch(strings.TrimSuffix(line, "."))
	r := bag{}
	if groups != nil && len(groups) > 1 {
		r.bag = groups[1]
		bags := strings.Split(groups[2], ",")
		for _, b := range bags {
			subgroups := subreg.FindStringSubmatch(strings.Trim(b, " "))
			if len(subgroups) == 3 {
				c, _ := strconv.Atoi(subgroups[1])
				bc := bagcount{bag: subgroups[2], count: c}
				r.contains = append(r.contains, bc)
			}
		}
	}

	return r
}

func canContainBag(r bag, b string, bm map[string]bag) bool {
	for _, bc := range r.contains {
		if bc.bag == b {
			return true
		}
		if bct, ok := bm[bc.bag]; ok {
			contains := canContainBag(bct, b, bm)
			if contains {
				return true
			}
		}

	}
	return false
}

func getContainedBagCount(bn string, bm map[string]bag) int {
	b, _ := bm[bn]

	i := 0
	for _, bc := range b.contains {
		i += bc.count + bc.count*getContainedBagCount(bc.bag, bm)
	}

	return i
}
