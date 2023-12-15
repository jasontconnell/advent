package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type lens struct {
	orig     string
	label    string
	val      int
	focallen int
}

type box struct {
	lenses []lens
}

func main() {
	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	steps := parseInput(in)
	return sum(compute(steps))
}

func part2(in input) output {
	steps := parseInput(in)
	lenses := compute(steps)
	m := initialize(lenses)
	return calcFocalLen(m)
}

func calcFocalLen(m map[int]box) int {
	total := 0
	for i := 0; i < 256; i++ {
		if _, ok := m[i]; !ok {
			continue
		}
		c := i + 1
		list := m[i].lenses
		for j := 0; j < len(list); j++ {
			x := c * ((j + 1) * list[j].focallen)
			total += x
		}
	}
	return total
}

func initialize(lenses []lens) map[int]box {
	m := make(map[int]box)
	for _, litr := range lenses {
		dash := strings.Index(litr.orig, "-")
		eq := strings.Index(litr.orig, "=")

		if dash > 0 {
			lbl := litr.orig[:dash]
			litr.label = lbl
			hash := computeOne(lbl).val
			if _, ok := m[hash]; ok && len(m[hash].lenses) > 0 {
				box := m[hash]
				for i := len(box.lenses) - 1; i >= 0; i-- {
					if box.lenses[i].label == lbl {
						box.lenses = append(box.lenses[:i], box.lenses[i+1:]...)
					}
				}
				m[hash] = box
			}
		} else if eq > 0 {
			sp := strings.Split(litr.orig, "=")
			lbl := sp[0]
			length, _ := strconv.Atoi(sp[1])
			hash := computeOne(lbl).val

			litr.label = lbl
			litr.focallen = length
			if _, ok := m[hash]; ok {
				box := m[hash]
				found := false
				for i := 0; i < len(box.lenses); i++ {
					if box.lenses[i].label == lbl {
						box.lenses[i] = litr
						found = true
					}
				}
				if !found {
					box.lenses = append(box.lenses, litr)
				}
				m[hash] = box
			} else {
				b := box{lenses: []lens{litr}}
				m[hash] = b
			}
		}
	}
	return m
}

func sum(results []lens) int {
	s := 0
	for _, r := range results {
		s += r.val
	}
	return s
}

func compute(steps []string) []lens {
	lenses := []lens{}
	for _, step := range steps {
		r := computeOne(step)
		lenses = append(lenses, r)
	}
	return lenses
}

func computeOne(step string) lens {
	r := lens{orig: step}
	for _, c := range step {
		r.val += int(c)
		r.val *= 17
		r.val = r.val % 256
	}
	return r
}

func parseInput(in input) []string {
	sp := strings.Split(in, ",")
	return sp
}
