package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type xypair struct {
	x1, x2 int
	y1, y2 int
}

type xy struct {
	x, y int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	area := parseInput(in)
	searchvec := xypair{20, 300, -80, 80}
	high, _ := findInitialVelocities(searchvec, area, false)
	return high.y
}

func part2(in input) output {
	area := parseInput(in)
	searchvec := xypair{20, 300, -80, 80}
	_, count := findInitialVelocities(searchvec, area, true)
	return count
}

func findInitialVelocities(search, target xypair, fullsearch bool) (xy, int) {
	start := xy{0, 0}
	count := 0
	maxheight := xy{0, 0}
	for y := search.y1; y < search.y2; y++ {
		for x := search.x1; x < search.x2; x++ {
			iv := xy{x, y}
			_, high, found := attemptInitialVelocity(start, iv, target, fullsearch)

			if found {
				count++
				if high.y > maxheight.y {
					maxheight = high
				}
			}
		}
	}
	return maxheight, count
}

// returns point of max height, point within area, true if it found the area
func attemptInitialVelocity(start, iv xy, area xypair, fullsearch bool) (xy, xy, bool) {
	found := false
	pos := start
	maxheight := start
	ivitr := iv
	var end xy
	for {
		pos.x += ivitr.x
		pos.y += ivitr.y

		if pos.x >= area.x1 && pos.x <= area.x2 && pos.y >= area.y1 && pos.y <= area.y2 {
			found = true
			end = pos
			break
		}

		if fullsearch && ivitr.x == 0 && (pos.x <= area.x1 || pos.x >= area.x2 || pos.y < area.y1) {
			found = false
			break
		}

		if pos.y > maxheight.y {
			maxheight = pos
		}

		if pos.y < area.y2 && !fullsearch {
			found = false
			break
		}

		if ivitr.x > 0 {
			ivitr.x--
		} else if ivitr.x < 0 {
			ivitr.x++
		}

		ivitr.y--
	}

	return end, maxheight, found
}

func parseInput(in input) xypair {
	reg := regexp.MustCompile(`target area: x=([0-9\-]+)\.\.([0-9\-]+), y=([0-9\-]+)\.\.([0-9\-]+)`)
	g := reg.FindStringSubmatch(in)

	ta := xypair{}
	if len(g) == 5 {
		x1, _ := strconv.Atoi(g[1])
		x2, _ := strconv.Atoi(g[2])
		y1, _ := strconv.Atoi(g[3])
		y2, _ := strconv.Atoi(g[4])
		ta.x1, ta.x2 = x1, x2
		ta.y1, ta.y2 = y1, y2
	}

	return ta
}
