package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

type grid map[xy]bool

type fold struct {
	ord rune
	val int
}

func maxcoords(g grid) (xy, xy) {
	minx, maxx := math.MaxInt32, math.MinInt32
	miny, maxy := math.MaxInt32, math.MinInt32

	for k, v := range g {
		if !v {
			continue
		}
		if k.x < minx {
			minx = k.x
		}

		if k.x > maxx {
			maxx = k.x
		}

		if k.y < miny {
			miny = k.y
		}

		if k.y > maxy {
			maxy = k.y
		}
	}

	return xy{minx, miny}, xy{maxx, maxy}
}
func printGrid(w io.Writer, g grid) {
	min, max := maxcoords(g)

	for y := min.y; y < max.y+1; y++ {
		for x := min.x; x < max.x+1; x++ {
			c := ' '
			if _, ok := g[xy{x, y}]; ok {
				c = '#'
			}

			fmt.Fprint(w, string(c))
		}
		fmt.Fprintln(w)
	}
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 13 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:")
	part2(w, in)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	coords, folds := parseInput(in)
	g := getGrid(coords)
	g = doFolds(g, folds[:1])
	return countDots(g)
}

func part2(w io.Writer, in input) {
	coords, folds := parseInput(in)
	g := getGrid(coords)
	g = doFolds(g, folds)
	printGrid(w, g)
}

func doFolds(g grid, folds []fold) grid {
	cp := make(grid, len(g))
	for k, v := range g {
		cp[k] = v
	}

	for _, f := range folds {
		switch f.ord {
		case 'x':
			cp = foldX(cp, f.val)
		case 'y':
			cp = foldY(cp, f.val)
		}
	}

	return cp
}

func foldX(g grid, x int) grid {
	cp := make(grid, len(g))
	for k, v := range g {
		cp[k] = v
	}

	for pt := range cp {
		if pt.x > x {
			delete(cp, pt)
			diff := pt.x - x
			np := x - diff

			cp[xy{np, pt.y}] = true
		}
	}
	return cp
}

func foldY(g grid, y int) grid {
	cp := make(grid, len(g))
	for k, v := range g {
		cp[k] = v
	}

	for pt := range cp {
		if pt.y > y {
			delete(cp, pt)
			diff := pt.y - y
			np := y - diff

			cp[xy{pt.x, np}] = true
		}
	}
	return cp
}

func countDots(g grid) output {
	count := 0
	for _, v := range g {
		if v {
			count++
		}
	}
	return count
}

func getGrid(coords []xy) grid {
	g := make(grid)
	for _, c := range coords {
		g[c] = true
	}
	return g
}

func parseInput(in input) ([]xy, []fold) {
	coords := []xy{}
	folds := []fold{}
	foldreg := regexp.MustCompile("fold along ([x|y])=([0-9]+)")

	coordmode := true
	for _, line := range in {
		if line == "" {
			coordmode = false
			continue
		}

		if coordmode {
			sp := strings.Split(line, ",")

			x, _ := strconv.Atoi(sp[0])
			y, _ := strconv.Atoi(sp[1])

			coords = append(coords, xy{x, y})
		} else {
			g := foldreg.FindStringSubmatch(line)
			if g != nil && len(g) == 3 {
				v, _ := strconv.Atoi(g[2])
				f := fold{ord: rune(g[1][0]), val: v}
				folds = append(folds, f)
			}
		}
	}
	return coords, folds
}
