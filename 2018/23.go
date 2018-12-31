package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var input = "23.txt"

type xyz struct {
	x, y, z int
}

func (p xyz) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.x, p.y, p.z)
}

type nanobot struct {
	xyz
	radius int
}

type nanobotdist struct {
	nanobot
	distance int
	inrange  bool
}

type result struct {
	xyz
	num int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	bots := []nanobot{}
	for scanner.Scan() {
		var txt = scanner.Text()
		bot := getBot(txt)
		if bot != nil {
			bots = append(bots, *bot)
		}
	}

	p1 := getInRangeOfStrongest(bots)
	p2 := getOptimal(bots)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func getStrongest(bots []nanobot) nanobot {
	strongest := nanobot{}
	maxr := 0
	for _, b := range bots {
		if b.radius > maxr {
			maxr = b.radius
			strongest = b
		}
	}

	return strongest
}

func getInRangeOfStrongest(bots []nanobot) int {
	strongest := getStrongest(bots)
	return getInRange(strongest.xyz, strongest.radius, bots)
}

func getOptimal(bots []nanobot) int {
	max := 0
	optimal := xyz{}
	r := 0

	for _, b := range bots {
		inrng := getInRangePoint(b.xyz, bots)
		if inrng > max {
			max = inrng
			optimal = b.xyz
			r = b.radius
		}
	}
	fmt.Println("Max in range is", max, optimal)

	dists := []nanobotdist{}
	for _, b := range bots {
		d := distance(optimal, b.xyz)
		inrange := d <= r
		dists = append(dists, nanobotdist{b, d, inrange})
	}

	sort.Slice(dists, func(i, j int) bool {
		return dists[i].distance < dists[j].distance
	})

	// xmv := optimal.x / 2
	// ymv := optimal.y / 2
	// zmv := optimal.z / 2
	// found := false
	fmt.Println("optimal is", optimal)

	// for !found {
	// 	pt := xyz{xmv, ymv, zmv}
	// 	num := getInRangePoint(pt, bots)
	// 	fmt.Println(num)
	// 	if num >= max {
	// 		max = num
	// 		optimal = pt
	// 	} else {
	// 		found = true
	// 	}

	// 	xmv--
	// 	ymv--
	// 	zmv--
	// }

	return distance(optimal, xyz{0, 0, 0})
}

func getInRangePoint(pt xyz, bots []nanobot) int {
	num := 0
	for _, b := range bots {
		if isInRange(pt, b.xyz, b.radius) {
			num++
		}
	}
	return num
}

func isInRange(pt, p2 xyz, radius int) bool {
	return distance(pt, p2) <= radius
}

func getInRange(pt xyz, radius int, bots []nanobot) int {
	num := 0
	for _, b := range bots {
		if isInRange(pt, b.xyz, radius) {
			num++
		}
	}

	return num
}

var reg *regexp.Regexp = regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)`)

func getBot(line string) *nanobot {
	var x, y, z, r int
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		x, _ = strconv.Atoi(groups[1])
		y, _ = strconv.Atoi(groups[2])
		z, _ = strconv.Atoi(groups[3])
		r, _ = strconv.Atoi(groups[4])
		return &nanobot{xyz{x, y, z}, r}
	}
	return nil
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func distance(p1, p2 xyz) int {
	dx := abs(p1.x - p2.x)
	dy := abs(p1.y - p2.y)
	dz := abs(p1.z - p2.z)
	return dx + dy + dz
}
