package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var input = "23.txt"

type xyz struct {
	x, y, z int
}

var origin xyz = xyz{0, 0, 0}

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

type state struct {
	point        xyz
	inrangeof    int
	distfromorig int
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

	var minsolvedist int
	var maxsolveinrng int
	maxdiv := 8 // 10^maxdiv
	itr := 1
	maxdev := 100
	optimalfound := false

	_, scalex, _, scaley, _, scalez := minmax(bots)
	maxpt := int(math.Max(float64(scalex), float64(scaley)))
	maxpt = int(math.Max(float64(maxpt), float64(scalez)))
	var optimal xyz

	for itr <= maxdiv {
		scale := 1
		if maxpt > 100 {
			scale = int(math.Pow10(maxdiv - itr))
		} else {
			itr = 9
		}
		clone := scalePoints(bots, scale)

		if !optimalfound {
			maxinrange := 0
			for _, b := range clone {
				inrng := getInRangePoint(b.xyz, clone)
				if inrng > maxinrange {
					optimal = b.xyz
					maxinrange = inrng
					optimalfound = true
				}
			}
		} else {
			optimal = scalePointUp(optimal, 10)
		}

		start := optimal

		solves := search(start, clone, maxdev)

		minsolvedist = math.MaxInt32
		maxsolveinrng = math.MinInt32

		for _, s := range solves {
			if s.inrangeof > maxsolveinrng || (s.inrangeof == maxsolveinrng && s.distfromorig < minsolvedist) {
				minsolvedist = s.distfromorig
				maxsolveinrng = s.inrangeof
				optimal = s.point
			}
		}
		itr++
	}

	return minsolvedist
}

func search(point xyz, bots []nanobot, maxdev int) []state {
	curInRange := getInRangePoint(point, bots)
	distorig := distance(origin, point)
	queue := []state{state{point: point, inrangeof: curInRange, distfromorig: distorig}}
	visited := make(map[xyz]bool)
	solves := []state{queue[0]}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]

		mvs := getMoves(st.point)
		for _, mv := range mvs {
			mvstate := state{point: mv, inrangeof: getInRangePoint(mv, bots), distfromorig: distance(origin, mv)}

			if mvstate.inrangeof > curInRange || (mvstate.inrangeof == curInRange && mvstate.distfromorig < distorig) {
				solves = append(solves, mvstate)
				curInRange = mvstate.inrangeof
			}

			if _, ok := visited[mv]; !ok {
				visited[mv] = true

				dev := distance(point, mv)
				if mvstate.inrangeof >= curInRange && dev < maxdev {
					queue = append(queue, mvstate)
				}
			}
		}
	}

	return solves
}

func scalePointUp(point xyz, scale int) xyz {
	point.x = point.x * scale
	point.y = point.y * scale
	point.z = point.z * scale
	return point
}

func scalePoints(points []nanobot, scale int) []nanobot {
	clone := append([]nanobot{}, points...)
	for i := 0; i < len(clone); i++ {
		pt := clone[i]
		pt.x = pt.x / scale
		pt.y = pt.y / scale
		pt.z = pt.z / scale
		pt.radius = pt.radius / scale
		clone[i] = pt
	}
	return clone
}

func getMoves(pt xyz) []xyz {
	vars := []xyz{}
	for _, p := range []xyz{
		xyz{x: 1, y: 0, z: 0},
		xyz{x: 0, y: 1, z: 0},
		xyz{x: 0, y: 0, z: 1},
		xyz{x: -1, y: 0, z: 0},
		xyz{x: 0, y: -1, z: 0},
		xyz{x: 0, y: 0, z: -1},
	} {
		vars = append(vars, xyz{x: pt.x + p.x, y: pt.y + p.y, z: pt.z + p.z})
	}
	return vars
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

func minmax(bots []nanobot) (minx, maxx, miny, maxy, minz, maxz int) {
	minx, maxx = math.MaxInt32, math.MinInt32
	miny, maxy = math.MaxInt32, math.MinInt32
	minz, maxz = math.MaxInt32, math.MinInt32

	for _, b := range bots {
		if b.x < minx {
			minx = b.x
		}

		if b.x > maxx {
			maxx = b.x
		}

		if b.y < miny {
			miny = b.y
		}

		if b.y > maxy {
			maxy = b.y
		}

		if b.z < minz {
			minz = b.z
		}

		if b.z > maxz {
			maxz = b.z
		}
	}

	return
}
