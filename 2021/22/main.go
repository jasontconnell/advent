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
type output = int64
type xyz struct {
	x, y, z int64
}
type xyzr struct {
	x1, x2, y1, y2, z1, z2 int64
}

func (r xyzr) String() string {
	return fmt.Sprintf("(%d,%d,%d) -> (%d, %d, %d) - volume %d", r.x1, r.y1, r.z1, r.x2, r.y2, r.z2, volume(r))
}

type instruction struct {
	xyzr
	status bool
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 22 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	list := parseInput(in)
	grid := map[xyzr]int64{}
	reboot(list, xyzr{-50, 50, -50, 50, -50, 50}, grid)

	return sum(grid)
}

func part2(in input) output {
	// list := parseInput(in)
	grid := map[xyzr]int64{}
	// reboot(list, xyzr{math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32}, grid)

	return sum(grid)
}

func sum(grid map[xyzr]int64) int64 {
	var total int64
	for _, v := range grid {
		total += v
	}
	return total
}

func countOn(grid map[xyzr]bool) int64 {
	var cnt int64
	for k, v := range grid {
		if !v {
			continue
		}
		val := volume(k)
		cnt += val

		fmt.Println(k, v)
	}
	return cnt
}

func reboot(list []instruction, bounds xyzr, grid map[xyzr]int64) int64 {
	for i := 0; i < len(list); i++ {
		inst := list[i]
		if !inbound(inst.xyzr, bounds) {
			continue
		}
		flipBits3(inst.xyzr, inst.status, grid)

	}
	// fmt.Println(grid)
	return 0
}

func volume(rng xyzr) int64 {
	xd := math.Abs(float64(rng.x2 - rng.x1))
	yd := math.Abs(float64(rng.y2 - rng.y1))
	zd := math.Abs(float64(rng.z2 - rng.z1))
	return int64(xd * yd * zd)
}

func flipBits3(rng xyzr, s bool, grid map[xyzr]int64) {
	if len(grid) == 0 && s {
		grid[rng] = volume(rng)
		return
	}
	fmt.Println(fmt.Sprintf("------------ %v %v -----------", rng, s))
	subs := map[xyzr]int64{}
	adds := map[xyzr]int64{}
	remove := []xyzr{}
	for k := range grid {
		list := split(rng, k)
		fmt.Println("splitting", k)

		remove = append(remove, k)
		for _, region := range list {
			fmt.Println("split", region)
			adds[region] = volume(region)
		}

		if !s {
			region, ok := overlapRegion(rng, k)
			if ok {
				delete(adds, region)
			}
		}
		os.Exit(0)
	}

	for _, region := range remove {
		delete(grid, region)
	}

	for region, v := range adds {
		grid[region] = v
	}

	for region, v := range subs {
		grid[region] = v
	}
}

func flipBits2(rng xyzr, s bool, grid map[xyzr]int64) {
	vol := volume(rng)

	subs := map[xyzr]int64{}
	adds := map[xyzr]int64{}
	fmt.Println(fmt.Sprintf("------------ %v %v -----------", rng, s))

	for k, v := range grid {
		region, overlaps := overlapRegion(rng, k)
		if overlaps {
			if s {
				vol -= volume(region)
				subs[k] = volume(region)
				adds[region] = volume(region)
				// fmt.Println("overlapped region, subtracted volume of region", region, "added region, sub from", k, volume(region))
			} else {
				outer := outerRegion(rng, k)
				fmt.Println("remove region", k, "unflip", outer)
				subs[k] = volume(outer)
			}
			fmt.Println("what to do with", region, v)
		}
	}

	for k, v := range subs {
		grid[k] -= v
	}

	for k, v := range adds {
		grid[k] += v
	}

	if s {
		grid[rng] = vol
	}
}

func max(i, j int64) int64 {
	m := i
	if j > m {
		m = j
	}
	return m
}

func min(i, j int64) int64 {
	m := i
	if j < m {
		m = j
	}
	return m
}

func outerPoints(p1, p2 int64) (int64, int64) {
	if p1 <= p2 {
		return p1, p2
	} else {
		return p2, p1
	}
}

func split(rng1, rng2 xyzr) []xyzr {
	region, overlap := overlapRegion(rng1, rng2)
	if !overlap {
		return []xyzr{rng1, rng2}
	}

	left, right := rng1, rng2
	if rng1.x1 < rng2.x1 {
		left, right = rng2, rng1
	}

	// top left
	s1 := xyzr{
		x1: left.x1, x2: right.x1, y1: left.y1, y2: right.y1, z1: left.z1, z2: right.z1,
	}

	// left's top right
	s2 := xyzr{
		x1: left.x1, x2: right.x1, y1: right.y1, y2: left.y1, z1: right.z1, z2: left.z1,
	}

	// left's bottom left
	s3 := xyzr{
		x1: left.x1, x2: right.x1, y1: right.y1, y2: left.y2, z1: right.z1, z2: left.z2,
	}

	// right's top right
	s4 := xyzr{
		x1: left.x2, x2: right.x2, y1: right.y1, y2: left.y2, z1: right.z1, z2: left.z2,
	}

	// right's bottom left
	s5 := xyzr{
		x1: right.x1, x2: left.x2, y1: left.y2, y2: right.y2, z1: right.z1, z2: left.z2,
	}

	// right's bottom right
	s6 := xyzr{
		x1: left.x2, x2: right.x2, y1: left.y2, y2: right.y2, z1: left.z2, z2: right.z2,
	}

	return []xyzr{
		s1, s2, s3, region, s4, s5, s6,
	}
}

func outerRegion(rng1, rng2 xyzr) xyzr {
	overlap, ok := overlapRegion(rng1, rng2)

	if !ok {
		return rng1
	}

	var right xyzr

	right.x1, right.x2 = outerPoints(max(rng1.x2, rng2.x2), max(overlap.x1, overlap.x2))
	right.y1, right.y2 = outerPoints(max(rng1.x2, rng2.x2), max(overlap.x1, overlap.x2))
	right.z1, right.z2 = outerPoints(max(rng1.x2, rng2.x2), max(overlap.x1, overlap.x2))

	return right
}

func overlapRegion(rng1, rng2 xyzr) (xyzr, bool) {
	region := xyzr{}
	if !hasOverlap(rng1, rng2) {
		return region, false
	}

	left, right := rng1, rng2
	if checkOverlap(rng2, rng1) {
		left, right = rng2, rng1
	}

	region.x1, region.x2 = regionPoints(left.x1, left.x2, right.x1, right.x2)
	region.y1, region.y2 = regionPoints(left.y1, left.y2, right.y1, right.y2)
	region.z1, region.z2 = regionPoints(left.z1, left.z2, right.z1, right.z2)

	return region, true
}

func regionPoints(left1, right1, left2, right2 int64) (int64, int64) {
	p1, p2 := left1, right2

	if between(left1, right1, left2) {
		p1 = left2
	} else if between(left2, right2, left1) {
		p1 = left1
	}

	if between(left1, right1, right2) {
		p2 = right2
	} else if between(left2, right2, right1) {
		p2 = right1
	}

	return p1, p2
}

func completeOverlap(rng1, rng2 xyzr) bool {
	return checkCompleteOverlap(rng1, rng2) || checkCompleteOverlap(rng2, rng1)
}

func checkCompleteOverlap(rng1, rng2 xyzr) bool {
	return between(rng1.x1, rng1.x2, rng2.x1) && between(rng1.x1, rng1.x2, rng2.x2) &&
		between(rng1.y1, rng1.y2, rng2.y1) && between(rng1.y1, rng1.y2, rng2.y2) &&
		between(rng1.z1, rng1.z2, rng2.z1) && between(rng1.z1, rng1.z2, rng2.z2)
}

func hasOverlap(rng1, rng2 xyzr) bool {
	return checkOverlap(rng1, rng2) || checkOverlap(rng2, rng1)
}

func checkOverlap(rng1, rng2 xyzr) bool {
	return (between(rng1.x1, rng1.x2, rng2.x1) || between(rng1.x1, rng1.x2, rng2.x2)) &&
		(between(rng1.y1, rng1.y2, rng2.y1) || between(rng1.y1, rng1.y2, rng2.y2)) &&
		(between(rng1.z1, rng1.z2, rng2.z1) || between(rng1.z1, rng1.z2, rng2.z2))

}

func between(p1, p2, v int64) bool {
	return v > p1 && v < p2
}

func inbound(coord, bounds xyzr) bool {
	return coord.x1 >= bounds.x1 && coord.x2 <= bounds.x2 &&
		coord.y1 >= bounds.y1 && coord.y2 <= bounds.y2 &&
		coord.z1 >= bounds.z1 && coord.z2 <= bounds.z2
}

func parseInput(in input) []instruction {
	reg := regexp.MustCompile(`(on|off) x=([0-9\-]+)\.\.([0-9\-]+),y=([0-9\-]+)\.\.([0-9\-]+),z=([0-9\-]+)\.\.([0-9\-]+)`)

	list := []instruction{}
	for _, line := range in {
		g := reg.FindStringSubmatch(line)
		if len(g) != 8 {
			continue
		}

		s := g[1] == "on"
		c := instruction{status: s}
		c.x1, c.x2 = atoip(g[2], g[3])
		c.y1, c.y2 = atoip(g[4], g[5])
		c.z1, c.z2 = atoip(g[6], g[7])
		list = append(list, c)
	}
	return list
}

func atoip(s1, s2 string) (int64, int64) {
	p1, _ := strconv.ParseInt(s1, 10, 64)
	p2, _ := strconv.ParseInt(s2, 10, 64)
	return p1, p2
}
