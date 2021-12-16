package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type packet struct {
	id      int
	version int
	packets []*packet

	literal *int
	value   int
}

func (p *packet) String() string {
	lit := "<nil>"
	if p.literal != nil {
		lit = fmt.Sprintf("%d", *p.literal)
	}
	s := fmt.Sprintf("id: %d version: %d literal: %s", p.id, p.version, lit)
	s += fmt.Sprintf("\n    value: %d ", p.value)
	if len(p.packets) > 0 {
		s += "\n      subpackets: "
		s += "\n         " + fmt.Sprintf("%v", p.packets)
	}
	return s
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
	fmt.Fprintln(w, "--2021 day 16 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	p := decode(in)
	return sumVersions(p)
}

func part2(in input) output {
	p := decode(in)
	return eval(p)
}

func sumVersions(p *packet) int {
	sum := p.version

	for _, sub := range p.packets {
		sum += sumVersions(sub)
	}
	return sum
}

func eval(p *packet) int {

	switch p.id {
	case 0:
		p.value = sum(p.packets)
	case 1:
		p.value = prod(p.packets)
	case 2:
		p.value = min(p.packets)
	case 3:
		p.value = max(p.packets)
	case 4:
		p.value = *p.literal
	case 5:
		p.value = gt(p.packets)
	case 6:
		p.value = lt(p.packets)
	case 7:
		p.value = eq(p.packets)
	}
	return p.value
}

func eq(s []*packet) int {
	v := 0
	v1, v2 := eval(s[0]), eval(s[1])
	if v1 == v2 {
		v = 1
	}
	return v
}

func lt(s []*packet) int {
	v := 0
	v1, v2 := eval(s[0]), eval(s[1])
	if v1 < v2 {
		v = 1
	}
	return v
}
func gt(s []*packet) int {
	v := 0
	v1, v2 := eval(s[0]), eval(s[1])
	if v1 > v2 {
		v = 1
	}
	return v
}

func max(s []*packet) int {
	m := math.MinInt64
	for _, sub := range s {
		n := eval(sub)
		if n > m {
			m = n
		}
	}
	return m
}

func min(s []*packet) int {
	m := math.MaxInt64
	for _, sub := range s {
		n := eval(sub)
		if n < m {
			m = n
		}
	}
	return m
}

func prod(s []*packet) int {
	v := 1
	for _, sub := range s {
		v *= eval(sub)
	}
	return v
}

func sum(s []*packet) int {
	v := 0
	for _, sub := range s {
		v += eval(sub)
	}
	return v
}

func decode(in input) *packet {
	bin := toBinary(in)
	p := &packet{}
	parsePacket(bin, p)
	return p
}

func parsePacket(bin string, p *packet) int {
	version := bin[:3]
	id := bin[3:6]
	pos := 6
	p.id = parseBinary(id)
	p.version = parseBinary(version)

	if p.id == 4 {
		num, np := parseLiteral(bin[pos:])
		p.literal = &num
		pos += np
	} else {
		lenval := bin[pos]
		pos++
		sp := 11
		packetmode := true
		if lenval == '0' {
			sp = 15
			packetmode = false
		}
		spstr := bin[pos : pos+sp]
		spval := parseBinary(spstr)
		pos += sp
		if !packetmode {
			sub := parsePackets(bin[pos : pos+spval])
			p.packets = append(p.packets, sub...)
			pos += spval
		} else {
			sub, np := parseSubpackets(bin[pos:], spval)
			p.packets = append(p.packets, sub...)
			pos += np
		}
	}
	return pos
}

func parseSubpackets(bin string, n int) ([]*packet, int) {
	subpackets := []*packet{}
	pos := 0
	for len(subpackets) < n {
		p := &packet{}
		np := parsePacket(bin[pos:], p)
		subpackets = append(subpackets, p)
		pos += np
	}
	return subpackets, pos
}

func parsePackets(bin string) []*packet {
	subpackets := []*packet{}
	pos := 0
	for pos < len(bin) {
		p := &packet{}
		np := parsePacket(bin[pos:], p)
		pos += np
		subpackets = append(subpackets, p)
	}
	return subpackets
}

func parseBinary(b string) int {
	n, _ := strconv.ParseInt(b, 2, 64)
	return int(n)
}

func parseLiteral(b string) (int, int) {
	pos := 0
	last := false
	strval := ""
	for !last {
		np := b[pos : pos+5]
		strval += np[1:]
		last = np[0] == '0'
		pos += 5
	}
	num := parseBinary(strval)
	return num, pos
}

func toBinary(in input) string {
	lookup := map[byte]string{
		'0': "0000", '1': "0001", '2': "0010", '3': "0011",
		'4': "0100", '5': "0101", '6': "0110", '7': "0111",
		'8': "1000", '9': "1001", 'A': "1010", 'B': "1011",
		'C': "1100", 'D': "1101", 'E': "1110", 'F': "1111",
	}
	b := ""
	for _, c := range in {
		b += lookup[byte(c)]
	}
	return b
}
