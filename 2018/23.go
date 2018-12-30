package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	//"strings"
	"math"
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
	
	p1 := getStrongest(bots)

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func getStrongest(bots []nanobot) int {
	strongest := nanobot{}
	maxr := 0
	for _, b := range bots {
		if b.radius > maxr {
			maxr = b.radius
			strongest = b
		}
	}

	num := 0
	for _, b := range bots {
		d := distance(strongest.xyz, b.xyz)

		if d <= strongest.radius {
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
		return &nanobot{xyz{x,y,z}, r}
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
