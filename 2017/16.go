package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	"strings"
	//"math"
)

var input = "16.txt"

type Move struct {
	Type      rune
	Spin      int
	Exchange1 int
	Exchange2 int
	Partner1  rune
	Partner2  rune
}

var programs []rune = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	moves := []Move{}

	for scanner.Scan() {
		var txt = scanner.Text()
		moves = getMoves(txt)
	}

	p2 := make([]rune, len(programs))
	copy(p2, programs)

	res := dance(moves, programs)
	res2 := danceForevs(moves, p2, 1000000000)

	fmt.Println("Dancing one time          :", string(res))
	fmt.Println("Dancing one billion times :", string(res2))

	fmt.Println("Time", time.Since(startTime))
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
