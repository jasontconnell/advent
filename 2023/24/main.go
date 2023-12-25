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

type xyz struct {
	x, y, z float64
}

func (p xyz) add(p2 xyz) xyz {
	return xyz{p.x + p2.x, p.y + p2.y, p.z + p2.z}
}

func (p xyz) mult(n int) xyz {
	ns := float64(n)
	return xyz{p.x * ns, p.y * ns, p.z * ns}
}

type hail struct {
	pos xyz
	vel xyz
}

var ex bool

func main() {
	file := common.InputFilename(os.Args)
	ex = file != "input.txt"
	in, err := common.ReadStrings(file)
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 24 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	list := parseInput(in)
	var start, end float64 = 200000000000000, 400000000000000
	if ex {
		start = 7
		end = 27
	}
	return testIntersections(list, start, end)
}

func part2(in input) output {
	return 0
}

func testIntersections(list []hail, testmin, testmax float64) int {
	total := 0
	for i, h := range list[:len(list)-1] {
		for _, h2 := range list[i+1:] {
			ipt := intersectionPoint2D(h.pos, h.vel, h2.pos, h2.vel)

			future := false

			if (ipt.x-h.pos.x)*h.vel.x >= 0 && (ipt.y-h.pos.y)*h.vel.y >= 0 &&
				(ipt.x-h2.pos.x)*h2.vel.x >= 0 && (ipt.y-h2.pos.y)*h2.vel.y >= 0 {
				future = true
			}

			if future && ipt.x >= testmin && ipt.x <= testmax && ipt.y >= testmin && ipt.y <= testmax {
				total++
			}
		}
	}

	return total
}

func intersectionPoint2D(p1, p1vel, p3, p3vel xyz) xyz {
	p2, p4 := p1.add(p1vel), p3.add(p3vel)
	a1 := p2.y - p1.y
	b1 := p1.x - p2.x
	c1 := a1*p1.x + b1*p1.y

	a2 := p4.y - p3.y
	b2 := p3.x - p4.x
	c2 := a2*p3.x + b2*p3.y

	det := a1*b2 - a2*b1
	if det == 0 {
		return xyz{0, 0, 0} // no intersection
	}

	nx := (b2*c1 - b1*c2) / det
	ny := (a1*c2 - a2*c1) / det

	return xyz{nx, ny, 0}
}

func slope(pt, vel xyz) float64 {
	var m float64
	fp := pt.add(vel)
	if fp.x == pt.x {
		m = 0
	} else {
		m = float64(fp.y-pt.y) / float64(fp.x-pt.x)
	}
	return m
}

func parseInput(in input) []hail {
	list := []hail{}
	for _, line := range in {
		sp := strings.Split(line, "@")
		posflds := strings.Fields(sp[0])
		pos := xyz{
			x: parseFloat(posflds[0]),
			y: parseFloat(posflds[1]),
			z: parseFloat(posflds[2]),
		}
		velflds := strings.Fields(sp[1])
		vel := xyz{
			x: parseFloat(velflds[0]),
			y: parseFloat(velflds[1]),
			z: parseFloat(velflds[2]),
		}

		list = append(list, hail{pos: pos, vel: vel})
	}
	return list
}

func parseFloat(str string) float64 {
	i, _ := strconv.Atoi(strings.Replace(str, ",", "", -1))
	return float64(i)
}
