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
	p2 := computeV2(masks)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
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

func computeV2(masks []*mask) int {
	mem := make(map[int]int)

	for _, m := range masks {
		for _, inst := range m.instructions {
			locs := getVariations(inst.loc, m.overrides)
			for _, loc := range locs {
				mem[loc] = inst.value
			}
		}
	}

	result := 0
	for _, x := range mem {
		result += x
	}
	return result
}

func addBit(bit int, arrs [][]int) {
	for i := 0; i < len(arrs); i++ {
		arrs[i] = append(arrs[i], bit)
	}
}

func getVariations(val int, ovr []int) []int {
	bits := getBits(val, 36)
	vmap := make(map[int][]int)

	for idx, v := range ovr {
		vmap[idx] = pvalues(v)
	}

	perms := getPermutations(bits, vmap)

	variations := []int{}
	for _, p := range perms {
		variations = append(variations, getNum(p))
	}

	return variations
}

func getPermutations(bits []int, vmap map[int][]int) [][]int {
	arrs := [][]int{}
	arrs = append(arrs, []int{})

	for i, b := range bits {
		if ov, ok := vmap[i]; ok {
			if len(ov) == 2 {
				olen := len(arrs)
				for x := 0; x < olen; x++ {
					c := make([]int, len(arrs[x]))
					copy(c, arrs[x])
					arrs[x] = append(arrs[x], ov[0])

					c = append(c, ov[1])
					arrs = append(arrs, c)
				}
			} else if len(ov) == 1 {
				for j := 0; j < len(arrs); j++ {
					nb := b
					if ov[0] == 1 {
						nb = 1
					}
					arrs[j] = append(arrs[j], nb)
				}
			}
		} else {
			panic("wrong")
		}
	}

	return arrs
}

func pvalues(b int) []int {
	if b == -1 {
		return []int{0, 1}
	}
	return []int{b}
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
