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

var partitrs []int = []int{2, 0}
var debug bool = false

func print(img [][]bool, sep, end string) {
	if !debug {
		return
	}

	if sep == "\n" {
		fmt.Println(len(img), "x", len(img[0]))
	}
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
	fmt.Fprintln(w, "--2017 day 21 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	start := [][]bool{{false, true, false}, {false, false, true}, {true, true, true}}
	rules := parseInput(in)
	rmap := getRuleMap(rules)
	return solve(start, rmap, partitrs[0])
}

func part2(in input) output {
	start := [][]bool{{false, true, false}, {false, false, true}, {true, true, true}}
	rules := parseInput(in)
	rmap := getRuleMap(rules)
	return solve(start, rmap, partitrs[1])
}

func solve(start [][]bool, rmap map[int][]Rule, loops int) int {
	img := start

	fmt.Println("start")
	print(img, "\n", "\n\n")
	for i := 0; i < loops; i++ {
		debug = true
		img = process(img, rmap)
		fmt.Println("-------- result --------")
		print(img, "\n", "\n\n")
	}
	return getOn(img)
}

func process(img [][]bool, rmap map[int][]Rule) [][]bool {
	cimg := copyimg(img)
	sp, didsplit := split(cimg)

	if didsplit {
		result := [][][][]bool{}
		for y := 0; y < len(sp); y++ {
			var rowresult [][][]bool
			for x := 0; x < len(sp[y]); x++ {
				nimg := sp[y][x]

				rule := getRule(nimg, rmap)
				enhanced := applyRule(nimg, rule)
				print(nimg, "/", "  ===>  ")
				print(rule.Match, "/", " -> ")
				print(rule.Enhancement, "/", "\n")
				rowresult = append(rowresult, enhanced)
			}
			result = append(result, rowresult)
		}

		fmt.Println("joining", len(result), len(result[0]))
		return join(result)
	} else {
		rule := getRule(img, rmap)
		enhanced := applyRule(img, rule)
		return enhanced
	}
}

func join(imgs [][][][]bool) [][]bool {
	if len(imgs) == 0 || len(imgs[0]) == 0 {
		log.Fatal("no records to join")
	}

	//how many
	joinw := len(imgs[0][0][0])
	joinh := len(imgs[0][0])

	fmt.Println("joinw", joinw, "joinh", joinh)

	rsize := sqrt(len(imgs)) * len(imgs[0])
	result := make([][]bool, rsize)
	for i := 0; i < rsize; i++ {
		result[i] = make([]bool, rsize)
	}

	fmt.Println(imgs)
	fmt.Println("len imgs", rsize, len(imgs))

	for ridx := 0; ridx < len(imgs); ridx++ {
		for cidx := 0; cidx < len(imgs[ridx]); cidx++ {
			fmt.Println("joining")
			print(imgs[ridx][cidx], "/", "\n")
			for y := 0; y < len(imgs[ridx][cidx]); y++ {
				dy := ridx*y + y
				for x := 0; x < len(imgs[ridx][cidx][y]); x++ {
					dx := cidx*x + x

					result[dy][dx] = imgs[ridx][cidx][y][x]
				}
			}
		}
	}

	return result
}

func applyRule(img [][]bool, rule Rule) [][]bool {
	cimg := copyimg(rule.Enhancement)
	return cimg
}

func split(cimg [][]bool) ([][][][]bool, bool) {
	if len(cimg) == 2 || len(cimg) == 3 {
		return nil, false
	}

	resultSize := 2
	if len(cimg)%3 == 0 {
		resultSize = 3
	}

	subgrids := len(cimg) / resultSize

	// it's a 4D array
	tmp := make([][][][]bool, resultSize)
	for ridx := 0; ridx < subgrids; ridx++ {
		tmp[ridx] = make([][][]bool, subgrids)
		for cidx := 0; cidx < subgrids; cidx++ {
			tmp[ridx][cidx] = make([][]bool, subgrids)
			for y := 0; y < subgrids; y++ {
				tmp[ridx][cidx][y] = make([]bool, subgrids)
			}
		}
	}

	for y := 0; y < len(cimg); y++ {
		cy := y % resultSize
		for x := 0; x < len(cimg[y]); x++ {
			cx := x % resultSize
			yidx := y / resultSize
			xidx := x / resultSize

			tmp[yidx][xidx][cy][cx] = cimg[y][x]
		}
	}

	return tmp, true
}

func sqrt(s int) int {
	return int(math.Sqrt(float64(s)))
}

func getRule(img [][]bool, rmap map[int][]Rule) Rule {
	ln := len(img)

	potentials := rmap[ln]

	fmt.Println("potentials", ln, len(potentials))
	print(img, "\n", "\n\n")
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
