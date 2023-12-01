package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type scanner struct {
	id     int
	points []xyz
}

type xyz struct {
	x, y, z int
}

var rotateX [][]int = [][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, -1}}
var rotateY [][]int = [][]int{{-1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
var rotateZ [][]int = [][]int{{1, 0, 0}, {0, -1, 0}, {0, 0, 1}}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	s := parseInput(in)
	return findBeacons(s, 12)
}

func part2(in input) output {
	return 0
}

func findBeacons(scanners []*scanner, matchBeacons int) int {
	for _, s := range scanners {
		for _, b := range s.points {
			fmt.Println(b, rotateX, b.x*rotateX[0][0], b.y*rotateX[1][1], b.z*rotateX[2][2])
			break
		}
	}
	return 0
}

func parseInput(in input) []*scanner {
	cur := &scanner{}
	scanners := []*scanner{}
	sreg := regexp.MustCompile("--- scanner ([0-9]+) ---")
	pmode := false
	for i, line := range in {
		g := sreg.FindStringSubmatch(line)
		if len(g) == 2 {
			id, _ := strconv.Atoi(g[1])
			cur = &scanner{id: id}
			pmode = true
			continue
		}

		if line == "" || i == len(in)-1 {
			scanners = append(scanners, cur)
			pmode = false
			continue
		}

		if pmode {
			flds := strings.Split(line, ",")

			x, _ := strconv.Atoi(flds[0])
			y, _ := strconv.Atoi(flds[1])
			z, _ := strconv.Atoi(flds[2])

			cur.points = append(cur.points, xyz{x, y, z})
		}
	}
	return scanners
}
