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
	} else {
		s += "["
	}
	for _, p := range p.subpackets {
		s += p.String() + ","
	}
	s = strings.TrimRight(s, ",")
	if p.val == nil {
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
	return getInOrder(packets)
}

func part2(in input) output {
	return 0
}

func getInOrder(p []*subpacket) int {
	rightOrder := 0
	for i := 0; i < len(p); i += 2 {
		left, right := p[i], p[i+1]

		res := compareSubpacket(left, right)

		restxt := "IN THE RIGHT ORDER"
		if res != -1 {
			restxt = "NOT " + restxt
		}

		pair := (i + 3) / 2
		fmt.Println("--------- PAIR ", pair, "--------------")
		fmt.Println("compare", left, right)
		fmt.Println(" ===   ", restxt)
		fmt.Println()

		if res == -1 {
			rightOrder += pair
		}
	}
	return rightOrder
}

func compareSubpacket(left, right *subpacket) int {
	if left.val != nil && right.val != nil {
		fmt.Println("compare vals", *left.val, *right.val)
		if *left.val < *right.val {
			return -1
		}
		if *left.val == *right.val {
			return 0
		}
		if *left.val > *right.val {
			return 1
		}
	}

	// left is val, right is list
	if left.val != nil && right.val == nil && len(right.subpackets) > 0 {
		fmt.Println("left val, no right val", left, right)
		lsub := &subpacket{val: left.val}
		left.val = nil
		left.subpackets = append(left.subpackets, lsub)
		return compareSubpacket(left, right)
	}

	// left is list, right is val
	if left.val == nil && right.val != nil && len(left.subpackets) > 0 {
		fmt.Println("no left val, right val", left, right)
		rsub := &subpacket{val: right.val}
		right.val = nil
		right.subpackets = append(right.subpackets, rsub)
		return compareSubpacket(left, right)
	}

	if len(left.subpackets) > 0 && len(right.subpackets) > 0 {
		fmt.Println("two lists", left, right)
		result := 0
		for i := 0; i < len(left.subpackets); i++ {
			if i >= len(right.subpackets) {
				result = 1
				break
			}
			subres := compareSubpacket(left.subpackets[i], right.subpackets[i])
			fmt.Println("subresult", left.subpackets[i], right.subpackets[i], subres)
			if subres != 0 {
				result = subres
				break
			}
		}
		if result <= 0 && len(left.subpackets) < len(right.subpackets) {
			result = -1
		}
		return result
	}

	if len(left.subpackets) == 0 && len(right.subpackets) > 0 {
		return -1
	}

	if len(left.subpackets) > 0 && len(right.subpackets) == 0 {
		return 1
	}

	return 1
}

func parseInput(in input) []*subpacket {
	packets := []*subpacket{}
	for _, line := range in {
		if line == "" {
			continue
		}

		csub := &subpacket{parent: nil}

		for _, c := range line {
			switch c {
			case '[':
				nsub := &subpacket{parent: csub}
				csub.subpackets = append(csub.subpackets, nsub)
				csub = nsub
			case ']':
				csub = csub.parent
			case ',':
			default:
				n, _ := strconv.Atoi(string(c))
				nsub := &subpacket{val: &n, parent: csub}
				csub.subpackets = append(csub.subpackets, nsub)
			}
		}

		packets = append(packets, csub)
	}
	return packets
}
