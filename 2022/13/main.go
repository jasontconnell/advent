package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type subpacket struct {
	val        *int
	subpackets []*subpacket
	parent     *subpacket
}

func (p *subpacket) String() string {
	s := ""
	if p.val != nil {
		s += strconv.Itoa(*p.val)
	} else if p.parent != nil || len(p.subpackets) > 0 {
		s += "["
	}
	for _, sp := range p.subpackets {
		s += sp.String() + ","
	}
	s = strings.TrimRight(s, ",")
	if p.val == nil && (p.parent != nil || len(p.subpackets) > 0) {
		s += "]"
	}
	return s
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 13 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	packets := parseInput(in)
	return orderedPairs(packets)
}

func part2(in input) output {
	packets := parseInput(in)
	dividers := []int{2, 6}
	divider1 := &subpacket{subpackets: []*subpacket{{val: &dividers[0]}}}
	divider2 := &subpacket{subpackets: []*subpacket{{val: &dividers[1]}}}
	packets = append(packets, divider1, divider2)
	packets = orderPackets(packets)

	res := 1
	for idx, p := range packets {
		if p == divider1 || p == divider2 {
			res *= (idx + 1)
		}
	}
	return res
}

func orderPackets(p []*subpacket) []*subpacket {
	sort.Slice(p, func(i, j int) bool {
		return compareSubpacket(p[i], p[j]) < 0
	})
	return p
}

func orderedPairs(p []*subpacket) int {
	rightOrder := 0
	for i := 0; i < len(p); i += 2 {
		pair := (i + 3) / 2
		left, right := p[i], p[i+1]

		res := compareSubpacket(left, right)

		if res < 0 {
			rightOrder += pair
		}
	}
	return rightOrder
}

func compareSubpacket(left, right *subpacket) int {
	if left.val != nil && right.val != nil {
		return *left.val - *right.val
	}

	// left is val, right is list
	if left.val != nil && right.val == nil {
		lsub := &subpacket{val: left.val}
		left.val = nil
		left.subpackets = append(left.subpackets, lsub)
		return compareSubpacket(left, right)
	}

	// left is list, right is val
	if left.val == nil && right.val != nil {
		rsub := &subpacket{val: right.val}
		right.val = nil
		right.subpackets = append(right.subpackets, rsub)
		return compareSubpacket(left, right)
	}

	if len(left.subpackets) > 0 || len(right.subpackets) > 0 {
		result := 0
		for i := 0; i < len(left.subpackets); i++ {
			if i >= len(right.subpackets) {
				result = 1
				break
			}
			subres := compareSubpacket(left.subpackets[i], right.subpackets[i])
			if subres != 0 {
				result = subres
				break
			}
		}
		if result == 0 && len(left.subpackets) < len(right.subpackets) {
			result = -1
		} else if result == 0 && len(left.subpackets) > len(right.subpackets) {
			result = 1
		}
		return result
	}

	if len(left.subpackets) == 0 && len(right.subpackets) == 0 {
		return 0
	}

	return 0
}

func parseInput(in input) []*subpacket {
	packets := []*subpacket{}
	for _, line := range in {
		if line == "" {
			continue
		}

		var csub *subpacket = &subpacket{}

		cnum := ""
		for _, c := range line {
			switch c {
			case '[':
				nsub := &subpacket{parent: csub}
				csub.subpackets = append(csub.subpackets, nsub)
				csub = nsub
			case ']':
				if cnum != "" {
					n, _ := strconv.Atoi(cnum)
					nsub := &subpacket{val: &n, parent: csub}
					csub.subpackets = append(csub.subpackets, nsub)
					cnum = ""
				}
				csub = csub.parent
			case ',':
				if cnum != "" {
					n, _ := strconv.Atoi(cnum)
					nsub := &subpacket{val: &n, parent: csub}
					csub.subpackets = append(csub.subpackets, nsub)
					cnum = ""
				}
			default:
				cnum += string(c)
				// n, _ := strconv.Atoi(string(c))
				// nsub := &subpacket{val: &n, parent: csub}
				// csub.subpackets = append(csub.subpackets, nsub)
			}
		}

		packets = append(packets, csub)
	}
	return packets
}
