package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

var input = "21.txt"
var start []bool = []bool{false, true, false, false, false, true, true, true, true}

var reg *regexp.Regexp = regexp.MustCompile("^(.*) => (.*)$")

type Rule struct {
	Match       []bool
	On          int
	Enhancement []bool
	Cols        int
}

func (r Rule) String() string {
	return fmt.Sprintf("m len: %v  on: %v cols: %v  => %v", len(r.Match), r.On, r.Cols, r.Match)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	rules := []Rule{}

	for scanner.Scan() {
		var txt = scanner.Text()
		r := getRule(txt)
		if r != nil {
			rules = append(rules, *r)
		}
	}

	rmap := getRuleMap(rules)

	on := solve(start, rmap, 5)
	fmt.Println("on  ", on)

	on2 := solve(start, rmap, 18)
	fmt.Println("on2 ", on2)

	fmt.Println("Time", time.Since(startTime))
}

func solve(start []bool, rmap map[int][]Rule, loops int) int {
	img := start
	for i := 0; i < loops; i++ {
		img = process(img, rmap)
	}

	return getOn(img)
}

func process(img []bool, rmap map[int][]Rule) []bool {
	size := int(math.Sqrt(float64(len(img))))
	dosplit := size%2 == 0 && size > 2 || size%3 == 0 && size > 3

	var enh []bool

	if dosplit {
		bb := split(img)

		for i := 0; i < len(bb); i++ {
			bb[i] = process(bb[i], rmap)
		}

		enh = join(bb)

	} else {
		potentials := []Rule{}

		for _, r := range rmap[size] {
			on := getOn(img)
			if r.On == on {
				potentials = append(potentials, r)
			}
		}

		for _, p := range potentials {
			if matches(p, img) {
				enh = p.Enhancement
				break
			}
		}
	}

	if len(enh) == 0 {
		panic("couldn't find match")
	}

	return enh
}

func join(bb [][]bool) []bool {
	var img []bool
	size := int(math.Sqrt(float64(len(bb[0]))))

	count := len(bb) * len(bb[0])
	chunk := 0

	bc := make([]int, len(bb))
	chunksize := len(bb[0])
	chunksper := int(math.Sqrt(float64(len(bb))))
	chunkline := 0
	mod := 3
	if size%2 == 0 {
		mod = 2
	}

	for i := 0; i < count; i++ {
		if i%mod == 0 && i > 0 {
			chunk++
			if i/chunksper > 0 {
				chunk = chunk % chunksper
			}

			if i%(chunksize*chunksper) == 0 {
				chunkline += chunksper
			}

			chunk += chunkline
		}

		img = append(img, bb[chunk][bc[chunk]])
		bc[chunk] = bc[chunk] + 1
	}

	return img
}

func split(img []bool) [][]bool {
	size := int(math.Sqrt(float64(len(img))))
	chunk := 0

	mod := 3
	if size%2 == 0 {
		mod = 2
	}

	chunksize := mod * mod
	chunks := len(img) / chunksize
	chunksper := int(math.Sqrt(float64(chunks)))
	chunkline := 0

	bb := make([][]bool, chunks)

	for i := 0; i < len(img); i++ {
		if i%mod == 0 && i > 0 {
			chunk++
			if i/chunksper > 0 {
				chunk = chunk % chunksper
			}

			if i%(chunksize*chunksper) == 0 {
				chunkline += chunksper
			}

			chunk += chunkline
		}

		bb[chunk] = append(bb[chunk], img[i])
	}

	return bb
}

func equals(a []bool, b []bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func matches(r Rule, img []bool) bool {
	size := int(math.Sqrt(float64(len(img))))

	m := false
	b := make([]bool, len(img))
	copy(b, img)

	for i := 0; i < size+1 && !m; i++ {
		m = equals(r.Match, b)
		if !m {
			b = rotate(b)
			m = equals(r.Match, b)
		}

		if !m {
			b = flip(b)
			m = equals(r.Match, b)

			if !m {
				b = flip(b)
			}
		}
	}

	return m
}

func getOn(img []bool) int {
	on := 0
	for _, b := range img {
		if b {
			on++
		}
	}
	return on
}

func getRuleMap(rules []Rule) map[int][]Rule {
	rmap := make(map[int][]Rule)
	for _, r := range rules {
		rmap[r.Cols] = append(rmap[r.Cols], r)
	}

	return rmap
}

func rotate(bs []bool) []bool {
	b := make([]bool, len(bs))
	size := int(math.Sqrt(float64(len(bs))))

	if size == 2 {
		b[0], b[1] = bs[2], bs[0]
		b[2], b[3] = bs[3], bs[1]
	} else if size == 3 {
		b[0], b[1], b[2] = bs[6], bs[3], bs[0]
		b[3], b[4], b[5] = bs[7], bs[4], bs[1]
		b[6], b[7], b[8] = bs[8], bs[5], bs[2]
	} else {
		panic("undefined size")
	}

	return b
}

func flip(bs []bool) []bool {
	b := make([]bool, len(bs))
	size := int(math.Sqrt(float64(len(bs))))

	if size == 3 {
		b[0], b[1], b[2] = bs[2], bs[1], bs[0]
		b[3], b[4], b[5] = bs[5], bs[4], bs[3]
		b[6], b[7], b[8] = bs[8], bs[7], bs[6]
	} else if size == 2 {
		b[0], b[1] = bs[1], bs[0]
		b[2], b[3] = bs[3], bs[2]
	}
	return b
}

func getRule(line string) *Rule {
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		match := []bool{}
		m := strings.Replace(groups[1], "/", "", -1)
		for _, c := range m {
			b := true
			if c == '.' {
				b = false
			}

			match = append(match, b)
		}

		enh := []bool{}
		e := strings.Replace(groups[2], "/", "", -1)
		for _, c := range e {
			b := true
			if c == '.' {
				b = false
			}

			enh = append(enh, b)
		}
		cols := 2
		if len(match) == 9 {
			cols = 3
		}

		return &Rule{Match: match, On: getOn(match), Enhancement: enh, Cols: cols}
	}
	return nil
}
