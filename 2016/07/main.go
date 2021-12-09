package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type ip struct {
	raw          string
	outer, inner []string
}

func (addr ip) String() string {
	return fmt.Sprintf("ip: %s outer: %s inner: %s", addr.raw, addr.outer, addr.inner)
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	ips := parseInput(in)
	s := 0
	for _, addr := range ips {
		if isTLS(addr) {
			s++
		}
	}
	return s
}

func part2(in input) output {
	ips := parseInput(in)
	s := 0
	for _, addr := range ips {
		if isSSL(addr) {
			s++
		}
	}
	return s
}

func parseInput(in input) []ip {
	ips := []ip{}
	breg := regexp.MustCompile(`\[.*?\]`)
	for _, line := range in {
		addr := ip{raw: line}
		noInner := breg.ReplaceAllString(line, " ")
		addr.outer = strings.Split(noInner, " ")

		m := breg.FindAllString(line, -1)
		for i := 0; i < len(m); i++ {
			m[i] = strings.ReplaceAll(strings.ReplaceAll(m[i], "]", ""), "[", "")
		}
		addr.inner = m

		ips = append(ips, addr)
	}
	return ips
}

func isSSL(addr ip) bool {
	innerBab := []string{}
	for _, s := range addr.inner {
		innerBab = append(innerBab, abaMatch(s)...)
	}

	outerAba := false
	for _, out := range addr.outer {
		for _, s := range innerBab {
			aba := string(s[1]) + string(s[0]) + string(s[1])
			if strings.Contains(out, aba) {
				outerAba = true
			}
		}

		if outerAba {
			break
		}
	}
	return outerAba
}

func isTLS(addr ip) bool {
	out, in := false, false
	for _, s := range addr.outer {
		out = out || hasAbba(s)
	}

	for _, s := range addr.inner {
		in = in || hasAbba(s)
	}

	return out && !in
}

func hasAbba(str string) bool {
	abba := false
	for i := 0; i < len(str)-3 && !abba; i++ {
		abba = abba || (str[i] == str[i+3] && str[i+1] == str[i+2] && str[i] != str[i+1])
	}

	return abba
}

func abaMatch(str string) []string {
	ret := []string{}
	for i := 0; i < len(str)-2; i++ {
		aba := (str[i] == str[i+2] && str[i] != str[i+1])
		if aba {
			ret = append(ret, string(str[i:i+2]))
		}
	}

	return ret
}
