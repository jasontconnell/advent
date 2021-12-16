package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = string

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return fill(in, 272)
}

func part2(in input) output {
	return fill(in, 35651584)
}

func fill(in input, disclen int) string {
	cp := in
	for len(cp) < disclen {
		cp = dragonCurve(cp)
	}

	disccp := string(cp[:disclen])
	sum, even := checksum(disccp)

	for even {
		sum, even = checksum(sum)
	}
	return sum
}

func dragonCurve(str string) string {
	cp := reverse(str)
	cp = strings.Replace(cp, "0", "_", -1)
	cp = strings.Replace(cp, "1", "0", -1)
	cp = strings.Replace(cp, "_", "1", -1)

	return str + "0" + cp
}

func reverse(str string) string {
	n := len(str)
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = rune(str[i])
	}

	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes)
}

func checksum(str string) (string, bool) {
	sum := make([]rune, len(str)/2+2)
	n := 0
	for i := 0; i < len(str)-1; i += 2 {
		c := '1'
		if str[i] != str[i+1] {
			c = '0'
		}
		sum[n] = c
		n++
	}

	return string(sum[:n]), n%2 == 0
}
