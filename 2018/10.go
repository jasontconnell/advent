package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "10.txt"

var reg *regexp.Regexp = regexp.MustCompile("position=< ?([0-9\\-]+),  ?([0-9\\-]+)> velocity=< ?([0-9\\-]+),  ?([0-9\\-]+)>")

type point struct {
	x, y int
}

type vector struct {
	start    point
	position point
	velocity point
}

func (v *vector) transform(s int) {
	x := v.start.x + (v.velocity.x * s)
	y := v.start.y + (v.velocity.y * s)

	v.position.x = x
	v.position.y = y
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	vectors := []*vector{}
	for scanner.Scan() {
		var txt = scanner.Text()
		v := parseVector(txt)
		if v == nil {
			fmt.Println("parsed wrong.", txt)
			os.Exit(1)
		}
		vectors = append(vectors, v)
	}

	part1(vectors)
	fmt.Println("Time", time.Since(startTime))
}

func clone(vectors []*vector) []*vector {
	c := []*vector{}
	for _, v := range vectors {
		c = append(c, &vector{position: v.position, velocity: v.velocity})
	}
	return c
}

func part1(vectors []*vector) {
	itr := 1
	done := false
	minsize := int64(200000000000000)

	for !done {
		processOne(vectors, itr)
		s := size(vectors)
		if s < minsize {
			minsize = s
		} else {
			fmt.Println(s, minsize)
			done = true
			processOne(vectors, itr-1)
		}
		itr++
	}

	sky := getSky(vectors)
	printSky(sky)

}

func processOne(vectors []*vector, scale int) {
	for i := 0; i < len(vectors); i++ {
		v := vectors[i]
		v.transform(scale)
	}
}

func printSky(sky [][]bool) {
	for i := 0; i < len(sky); i++ {
		ln := ""
		for j := 0; j < len(sky[i]); j++ {
			c := " "
			if sky[i][j] {
				c = "#"
			}
			ln += c
		}
		fmt.Println(ln)
	}
}

func getSky(vectors []*vector) [][]bool {
	_, x1, _, y1, shx, shy := getNormalized(vectors)

	sky := make([][]bool, y1+1)
	for i := 0; i < y1+1; i++ {
		sky[i] = make([]bool, x1+1)
	}
	for _, v := range vectors {
		x := v.position.x - shx
		y := v.position.y - shy

		if y < len(sky) && x < len(sky[y]) {
			sky[y][x] = true
		} else {
			panic(fmt.Sprintf("out of range (%d, %d)  lengths: (%d, %d)", x, y, len(sky), len(sky[x])))
		}
	}

	return sky
}

func getNormalized(vectors []*vector) (int, int, int, int, int, int) {
	x0, x1, y0, y1 := minMax(vectors)
	shx, shy := getShift(x0, x1, y0, y1)

	if x0 < 0 {
		x1 = (-x0) + x1
	} else {
		x1 = x1 - x0
	}
	x0 = 0

	if y1 < 0 {
		y1 = (-y0) + y1
	} else {
		y1 = y1 - y0
	}
	y0 = 0

	return x0, x1, y0, y1, shx, shy
}

func getShift(x0, x1, y0, y1 int) (int, int) {
	return x0, y0
}

func size(vectors []*vector) int64 {
	x0, x1, y0, y1, _, _ := getNormalized(vectors)

	sx := int64(x1 - x0)
	sy := int64(y1 - y0)

	return sx * sy
}

func minMax(vectors []*vector) (int, int, int, int) {
	b := 1000000
	minx, miny := b, b
	maxx, maxy := -b, -b
	for _, v := range vectors {
		if v.position.x < minx {
			minx = v.position.x
		}

		if v.position.y < miny {
			miny = v.position.y
		}

		if v.position.x > maxx {
			maxx = v.position.x
		}

		if v.position.y > maxy {
			maxy = v.position.y
		}
	}
	return minx, maxx, miny, maxy
}

func parseVector(txt string) *vector {
	if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
		px, _ := strconv.Atoi(groups[1])
		py, _ := strconv.Atoi(groups[2])
		vx, _ := strconv.Atoi(groups[3])
		vy, _ := strconv.Atoi(groups[4])

		v := &vector{start: point{x: px, y: py}, position: point{x: px, y: py}, velocity: point{x: vx, y: vy}}
		return v
	}
	return nil
}
