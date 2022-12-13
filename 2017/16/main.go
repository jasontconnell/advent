package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = string

type Move struct {
	Type      rune
	Spin      int
	Exchange1 int
	Exchange2 int
	Partner1  rune
	Partner2  rune
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	programs := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
	moves := getMoves(in)
	res := dance(moves, programs)

	return string(res)
}

func part2(in input) output {
	programs := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
	moves := getMoves(in)
	res := danceForevs(moves, programs, 1000000000)
	return string(res)
}

func danceForevs(moves []Move, p []rune, iterations int) []rune {
	first := string(p)
	for i := 1; i <= iterations; i++ {
		p = dance(moves, p)

		if i != 0 && string(p) == first {
			i = i * int(iterations/i)
		}

	}

	return p
}

func dance(moves []Move, p []rune) []rune {
	for _, mv := range moves {
		switch mv.Type {
		case 's':
			p = spin(mv.Spin, p)
		case 'x':
			p = exchange(mv.Exchange1, mv.Exchange2, p)
		case 'p':
			p = partner(mv.Partner1, mv.Partner2, p)
		}
	}

	return p
}

func spin(size int, p []rune) []rune {
	s1 := p[len(p)-size:]
	//fmt.Println(size, len(s1), string(s1))
	return append(s1, p[0:len(p)-size]...)
}

func exchange(x1, x2 int, p []rune) []rune {
	p[x2], p[x1] = p[x1], p[x2]
	return p
}

func partner(r1, r2 rune, p []rune) []rune {
	var i1 int
	var i2 int

	for i := 0; i < len(p); i++ {
		if p[i] == r1 {
			i1 = i
		}
		if p[i] == r2 {
			i2 = i
		}
	}

	return exchange(i1, i2, p)
}

func getMoves(line string) []Move {
	moves := []Move{}
	mvs := strings.Split(line, ",")

	for _, mv := range mvs {
		moves = append(moves, getMove(mv))
	}

	return moves
}

func getMove(mv string) Move {
	m := Move{Type: rune(mv[0])}

	if m.Type == 's' {
		n, err := strconv.Atoi(string(mv[1:]))
		if err != nil {
			panic(err)
		}
		m.Spin = n
	} else if m.Type == 'x' {
		args := strings.Split(string(mv[1:]), "/")

		x1, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		x2, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		m.Exchange1 = x1
		m.Exchange2 = x2
	} else if m.Type == 'p' {
		m.Partner1 = rune(mv[1])
		m.Partner2 = rune(mv[3])
	}

	return m
}
