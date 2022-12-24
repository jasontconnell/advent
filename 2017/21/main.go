package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Rule struct {
	Match       [][]bool
	On          int
	Enhancement [][]bool
	Cols        int
}

func print(img [][]bool, sep, end string) {
	for i, r := range img {
		for _, c := range r {
			if c {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		if i < len(r)-1 {
			fmt.Print(sep)
		}
	}
	fmt.Print(end)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	start := [][]bool{{false, true, false}, {false, false, true}, {true, true, true}}
	rules := parseInput(in)
	rmap := getRuleMap(rules)
	return solve(start, rmap, 3)
}

func part2(in input) output {
	start := [][]bool{{false, true, false}, {false, false, true}, {true, true, true}}
	rules := parseInput(in)
	rmap := getRuleMap(rules)
	return 0
	return solve(start, rmap, 18)
}

func solve(start [][]bool, rmap map[int][]Rule, loops int) int {
	img := start
	for i := 0; i < loops; i++ {
		img = process(img, rmap)
		print(img, "\n", "")
		fmt.Println()
	}
	return getOn(img)
}

func process(img [][]bool, rmap map[int][]Rule) [][]bool {
	cimg := copyimg(img)
	sp := split(cimg)

	result := [][][]bool{}
	for i := 0; i < len(sp); i++ {
		nimg := sp[i]

		rule := getRule(nimg, rmap)
		print(nimg, "/", "")
		fmt.Print(" matches -> ")
		print(rule.Match, "/", "")
		fmt.Print(" -> ", "")
		print(rule.Enhancement, "/", "\n")
		enhanced := applyRule(nimg, rule)
		result = append(result, enhanced)
	}

	joined := join(result)
	fmt.Println("joined")
	print(joined, "/", "\n")

	return joined
}

func join(imgs [][][]bool) [][]bool {
	rsize := sqrt(len(imgs)) * len(imgs[0])
	result := make([][]bool, rsize)
	for i := 0; i < rsize; i++ {
		result[i] = make([]bool, rsize)
	}

	sq := sqrt(len(imgs))
	for i := 0; i < len(imgs); i++ {
		for y := 0; y < len(imgs[i]); y++ {
			for x := 0; x < len(imgs[i][y]); x++ {
				xb := i / sq
				yb := i % sq
				row := (y % len(imgs[i])) + y/len(imgs[i]) + (len(imgs[i]) * yb)
				col := (x % len(imgs[i][y])) + x/len(imgs[i][y]) + (len(imgs[i][y]) * xb)

				fmt.Println(row, col, y, x)
				result[row][col] = imgs[i][y][x]
			}
		}
	}
	return result
}

func applyRule(img [][]bool, rule Rule) [][]bool {
	cimg := copyimg(rule.Enhancement)
	return cimg
}

func split(cimg [][]bool) [][][]bool {
	size := int(math.Sqrt(float64(len(cimg) * len(cimg[0]))))
	if size == 2 || size == 3 {
		return [][][]bool{cimg}
	}

	var nsize int
	if size%2 == 0 {
		nsize = 2
	} else {
		nsize = 3
	}
	num := (len(cimg) * len(cimg[0])) / (nsize * nsize)

	fmt.Println("split")
	print(cimg, "/", "\n")

	ret := [][][]bool{}
	for i := 0; i < num; i++ {
		nimg := make([][]bool, nsize)
		for j := 0; j < nsize; j++ {
			nimg[j] = make([]bool, nsize)
		}
		ret = append(ret, nimg)
	}

	tmp := make([][][]bool, len(cimg)/nsize)
	idx := 0
	for i := 0; i < len(cimg)/nsize; i++ {
		tmp[i] = make([][]bool, nsize)
		tmp[i][0] = cimg[idx]
		tmp[i][1] = cimg[idx+1]
		if nsize == 3 {
			tmp[i][2] = cimg[idx+2]
		}
		idx += nsize
	}

	for i := 0; i < len(tmp); i++ {
		for y := 0; y < len(tmp[i]); y++ {
			for x := 0; x < len(tmp[i][y]); x++ {
				ridx := (i*nsize + i/nsize) + x/nsize + y/nsize
				ret[ridx][y][x%nsize] = tmp[i][y][x]
			}
		}
	}

	fmt.Println("after split")
	for _, r := range ret {
		print(r, "/", "\n")
	}

	return ret
}

func sqrt(s int) int {
	return int(math.Sqrt(float64(s)))
}

func getRule(img [][]bool, rmap map[int][]Rule) Rule {
	ln := len(img)

	potentials := rmap[ln]

	var r *Rule
	for _, pr := range potentials {
		m := isMatch(img, pr)
		if m {
			r = &pr
			break
		}
	}
	return *r
}

func isMatch(img [][]bool, rule Rule) bool {
	cimg := copyimg(img)
	match := true

	apply := []func([][]bool) [][]bool{}
	apply = append(apply, noop, rotate, rotate, rotate, flip, rotate, rotate, rotate)

	for _, fn := range apply {
		cimg = fn(cimg)
		match = true
		for y := 0; y < len(cimg) && match; y++ {
			for x := 0; x < len(cimg[y]) && match; x++ {
				match = match && cimg[y][x] == rule.Match[y][x]
			}
		}
		if match {
			break
		}
	}
	return match
}

func copyimg(img [][]bool) [][]bool {
	cp := make([][]bool, len(img))
	for y := 0; y < len(img); y++ {
		cp[y] = make([]bool, len(img[y]))
		copy(cp[y], img[y])
	}
	return cp
}

// noop
func noop(img [][]bool) [][]bool {
	return img
}

// rotate clockwise
func rotate(img [][]bool) [][]bool {
	cp := make([][]bool, len(img))
	for i := 0; i < len(cp); i++ {
		cp[i] = make([]bool, len(img[i]))
	}

	for y := 0; y < len(cp); y++ {
		for x := 0; x < len(cp[y]); x++ {
			cp[y][x] = img[len(img[y])-1-x][y]
		}
	}
	return cp
}

// flip over (left becomes right)
func flip(img [][]bool) [][]bool {
	cp := make([][]bool, len(img))
	for i := 0; i < len(cp); i++ {
		cp[i] = make([]bool, len(img[i]))
	}

	for y := 0; y < len(cp); y++ {
		for x := 0; x < len(cp[y]); x++ {
			cp[y][x] = img[y][len(img[y])-x-1]
		}
	}
	return cp
}

func getOn(img [][]bool) int {
	on := 0
	for _, r := range img {
		for _, c := range r {
			if c {
				on++
			}
		}
	}
	return on
}

func getRuleMap(rules []Rule) map[int][]Rule {
	m := make(map[int][]Rule)
	for _, r := range rules {
		m[r.Cols] = append(m[r.Cols], r)
	}
	return m
}

func parseInput(in input) []Rule {
	rules := []Rule{}

	for _, line := range in {
		sp := strings.Split(line, " => ")
		rstr := sp[0]
		enhstr := sp[1]

		match := parsePart(rstr)
		enh := parsePart(enhstr)

		rules = append(rules, Rule{Match: match, On: getOn(match), Enhancement: enh, Cols: len(match)})
	}
	return rules
}

func parsePart(str string) [][]bool {
	sp := strings.Split(str, "/")
	bb := [][]bool{}
	for _, s := range sp {
		r := []bool{}
		for _, c := range s {
			b := true
			if c == '.' {
				b = false
			}
			r = append(r, b)
		}
		bb = append(bb, r)
	}
	return bb
}
