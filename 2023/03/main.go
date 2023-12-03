package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (p xy) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type num struct {
	number int
	start  xy
	end    xy
}

func (n num) String() string {
	return fmt.Sprintf("%d %v %v", n.number, n.start, n.end)
}

type symbol struct {
	sym   rune
	point xy
}

func (s symbol) String() string {
	return fmt.Sprintf("%c %v", s.sym, s.point)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 03 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	n, s := parseInput(in)
	pns := getPartNumbers(n, s)
	return sum(pns)
}

func part2(in input) output {
	return 0
}

func sum(nums []int) int {
	s := 0
	for _, n := range nums {
		s += n
	}
	return s
}

func getPartNumbers(nums []num, syms []symbol) []int {
	sm := make(map[xy]symbol)

	for _, s := range syms {
		sm[s.point] = s
	}

	partnums := []int{}
	for _, n := range nums {
		found := false
		for y := n.start.y - 1; y <= n.end.y+1 && !found; y++ {
			for x := n.start.x - 1; x <= n.end.x+1 && !found; x++ {
				pt := xy{x, y}

				if _, ok := sm[pt]; ok {
					partnums = append(partnums, n.number)
					found = true
				}
			}
		}
	}
	return partnums
}

func parseInput(in input) ([]num, []symbol) {
	nums := []num{}
	symbols := []symbol{}

	for y, line := range in {
		cn := ""
		cns := xy{-1, y}
		innum := false
		for x, c := range line + "." { // catch end of line numbers
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				cn += string(c)
				innum = true
				if cns.x < 0 {
					cns.x = x
				}
			default:
				if innum {
					n, _ := strconv.Atoi(cn)
					nums = append(nums, num{number: n, start: cns, end: xy{cns.x + len(cn) - 1, y}})
				}

				if c != '.' {
					symbols = append(symbols, symbol{sym: c, point: xy{x, y}})
				}
				cn = ""
				cns.x = -1
				innum = false
			}
		}
	}

	return nums, symbols
}
