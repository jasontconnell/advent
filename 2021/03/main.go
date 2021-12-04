package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) uint64 {
	gamma, epsilon := getGammaEpsilon(in)
	return gamma * epsilon
}

func part2(in input) uint64 {
	sograte := narrow(in, false)
	sco2scrub := narrow(in, true)

	ograting, _ := strconv.ParseUint(sograte, 2, len(sograte))
	co2scrub, _ := strconv.ParseUint(sco2scrub, 2, len(sco2scrub))

	return ograting * co2scrub
}

func getGammaEpsilon(in input) (uint64, uint64) {
	var l int = len(in[0])
	var zeros, ones = make([]int, l), make([]int, l)

	for i := 0; i < l; i++ {
		b := commonBit(in, i)
		switch b {
		case "0":
			zeros[i]++
		case "1":
			ones[i]++
		}
	}

	gs, es := "", ""
	for i := 0; i < l; i++ {
		z, o := zeros[i], ones[i]

		if z > o {
			gs += "0"
			es += "1"
		} else {
			gs += "1"
			es += "0"
		}
	}

	gamma, _ := strconv.ParseUint(gs, 2, len(gs))
	epsilon, _ := strconv.ParseUint(es, 2, len(es))

	return gamma, epsilon
}

func narrow(vals input, flip bool) string {
	cp := make(input, len(vals))
	copy(cp, vals)
	pos := 0
	for len(cp) > 1 {
		b := commonBit(cp, pos)

		if flip {
			if b == "0" {
				b = "1"
			} else {
				b = "0"
			}
		}

		for i := len(cp) - 1; i >= 0; i-- {
			if cp[i][pos] != b[0] {
				cp = append(cp[:i], cp[i+1:]...)
			}
		}

		pos++
	}

	return cp[0]
}

func commonBit(vals input, pos int) string {
	var z, o int
	for _, val := range vals {
		switch val[pos] {
		case '0':
			z++
		case '1':
			o++
		}
	}

	val := "1"
	if z > o {
		val = "0"
	}
	return val
}
