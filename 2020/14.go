package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var input = "14.txt"

type instruction struct {
	loc   int
	value int
}

type mask struct {
	overrides    []int
	instructions []instruction
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	masks := getMasks(lines)
	p1 := compute(masks)

	fmt.Println("Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
}

func getBits(num, length int) []int {
	bits := []int{}
	x := num
	pos := 1
	for x > 0 {
		div := int(math.Pow(2, float64(pos)))
		bits = append(bits, x%div)
		x = x / div
	}

	for len(bits) < length {
		bits = append(bits, 0)
	}

	return bits
}

func getNum(bits []int) int {
	x := 0
	for idx, v := range bits {
		x += (v * int(math.Pow(2, float64(idx))))
	}
	return x
}

func compute(masks []*mask) int {
	mem := make(map[int]int)

	for _, m := range masks {
		for _, inst := range m.instructions {
			bits := getBits(inst.value, 36)

			for i, x := range m.overrides {
				if x != -1 {
					bits[i] = x
				}
			}

			mem[inst.loc] = getNum(bits)
		}
	}

	result := 0
	for _, x := range mem {
		result += x
	}
	return result
}

var maskreg *regexp.Regexp = regexp.MustCompile("^mask = ([X01]+)$")
var memreg *regexp.Regexp = regexp.MustCompile(`^mem\[([0-9]+)] = ([0-9]+)$`)

func getMasks(lines []string) []*mask {
	masks := []*mask{}
	cur := &mask{}
	for _, line := range lines {
		mg := maskreg.FindStringSubmatch(line)

		if len(mg) == 2 {
			cur = &mask{}
			masks = append(masks, cur)

			for _, c := range mg[1] {
				v := 1
				switch c {
				case 'X':
					v = -1
				case '0':
					v = 0
				}
				cur.overrides = append([]int{v}, cur.overrides...)
			}
			continue
		}

		memg := memreg.FindStringSubmatch(line)
		if len(memg) == 3 {
			loc, _ := strconv.Atoi(memg[1])
			val, _ := strconv.Atoi(memg[2])

			cur.instructions = append(cur.instructions, instruction{loc, val})
		}
	}

	return masks
}
