package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2025 day 05 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	ranges, nums := parseInput(in)
	return getInRanges(ranges, nums)
}

func part2(in input) output {
	ranges, _ := parseInput(in)
	ranges = collapse(ranges)
	return getUniqueRanges(ranges)
}

func getInRanges(ranges []xy, nums []int) int {
	inrange := 0
	for _, i := range nums {
		for _, r := range ranges {
			if i >= r.x && i <= r.y {
				inrange++
				break
			}
		}
	}
	return inrange
}

func collapse(ranges []xy) []xy {
	for i := len(ranges) - 1; i >= 0; i-- {
		r := ranges[i]

		for j := len(ranges) - 1; j >= 0; j-- {
			if i == j {
				continue
			}

			modified := false
			if r.x >= ranges[j].x && r.y <= ranges[j].y {
				// completely encapsulated by another one, unnecessary
				ranges = append(ranges[:i], ranges[i+1:]...)
				modified = true
			} else if r.x >= ranges[j].x && r.x <= ranges[j].y {
				if r.y > ranges[j].y {
					ranges[j].y = r.y
					ranges = append(ranges[:i], ranges[i+1:]...)
					modified = true
				} else if r.y < ranges[j].y {
					ranges = append(ranges[:i], ranges[i+1:]...)
					modified = true
				}
			} else if r.y >= ranges[j].x && r.y <= ranges[j].y {
				if r.x < ranges[j].x {
					ranges[j].x = r.x
					ranges = append(ranges[:i], ranges[i+1:]...)
					modified = true
				} else if r.x > ranges[j].x {
					ranges = append(ranges[:i], ranges[i+1:]...)
				}
			}
			if modified {
				break
			}
		}
	}
	return ranges
}

func getUniqueRanges(ranges []xy) int {
	total := 0
	for _, r := range ranges {
		total += (r.y - r.x) + 1
	}
	return total
}

func parseInput(in input) ([]xy, []int) {
	ranges := []xy{}
	nums := []int{}

	nstart := 0
	for i, line := range in {
		if line == "" {
			nstart = i
			break
		}

		sp := strings.Split(line, "-")
		x, _ := strconv.Atoi(sp[0])
		y, _ := strconv.Atoi(sp[1])

		ranges = append(ranges, xy{x, y})
	}

	for i := nstart; i < len(in); i++ {
		x, _ := strconv.Atoi(in[i])
		nums = append(nums, x)
	}

	return ranges, nums
}
