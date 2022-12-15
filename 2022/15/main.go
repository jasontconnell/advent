package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

type sensor struct {
	pt       xy
	beacon   xy
	strength int
}

var isexample bool

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}
	isexample = common.InputFilename(os.Args) != "input.txt"

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 15 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	sensors := parseInput(in)
	row := 2000000
	if isexample {
		row = 10
	}
	return findVoids(sensors, row)
}

func part2(in input) output {
	sensors := parseInput(in)
	low := 0
	high := 4000000
	if isexample {
		high = 20
	}
	return findBeacon(sensors, low, high, 4000000)
}

func findBeacon(sensors []sensor, searchlow, searchhigh, mult int) int {
	pt := xy{}
	found := false
	for y := searchlow; y <= searchhigh && !found; y++ {
		for x := searchlow; x <= searchhigh && !found; x++ {
			inrange := false
			for _, s := range sensors {
				d := dist(s.pt, xy{x, y})
				if d <= s.strength {
					inrange = true
					skip := s.strength - d

					x += skip
					break
				}
			}

			if !inrange {
				pt = xy{x, y}
				found = true
			}
		}
	}

	return pt.x*mult + pt.y
}

func findVoids(sensors []sensor, row int) int {
	min, max := minmax(sensors)
	total := 0
	for i := min.x; i < max.x+1; i++ {
		pt := xy{i, row}

		inrange := false
		for _, s := range sensors {
			if dist(s.pt, pt) <= s.strength && pt != s.beacon {
				inrange = true
				break
			}
		}

		if inrange {
			total++
		}
	}
	return total
}

func minmax(sensors []sensor) (xy, xy) {
	min, max := xy{math.MaxInt32, math.MaxInt32}, xy{math.MinInt32, math.MinInt32}
	for _, s := range sensors {
		if s.pt.x-s.strength < min.x {
			min.x = s.pt.x - s.strength
		}

		if s.pt.x+s.strength > max.x {
			max.x = s.pt.x + s.strength
		}

		if s.pt.y-s.strength < min.y {
			min.y = s.pt.y - s.strength
		}

		if s.pt.y+s.strength > max.y {
			max.y = s.pt.y + s.strength
		}
	}
	return min, max
}

func dist(p1, p2 xy) int {
	dx := math.Abs(float64(p2.x - p1.x))
	dy := math.Abs(float64(p2.y - p1.y))
	return int(dx + dy)
}

func parseInput(in input) []sensor {
	reg := regexp.MustCompile(`Sensor at x=([0-9\-]+), y=([0-9\-]+): closest beacon is at x=([0-9\-]+), y=([0-9\-]+)`)
	sensors := []sensor{}
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) == 5 {
			sx, _ := strconv.Atoi(m[1])
			sy, _ := strconv.Atoi(m[2])
			bx, _ := strconv.Atoi(m[3])
			by, _ := strconv.Atoi(m[4])

			sxy := xy{sx, sy}
			bxy := xy{bx, by}
			s := sensor{pt: sxy, beacon: bxy, strength: dist(sxy, bxy)}
			sensors = append(sensors, s)
		}
	}
	return sensors
}
