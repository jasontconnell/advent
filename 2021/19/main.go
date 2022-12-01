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

type xyz struct {
	x, y, z int
}

func (pt xyz) add(pt2 xyz) xyz {
	res := xyz{pt.x + pt2.x, pt.y + pt2.y, pt.z + pt2.z}
	return res
}

func (pt xyz) delta(pt2 xyz) xyz {
	res := xyz{pt.x - pt2.x, pt.y - pt2.y, pt.z - pt2.z}
	return res
}

type transform int

const (
	x0 transform = iota
	x90
	x180
	x270
	y0
	y90
	y180
	y270
	z0
	z90
	z180
	z270
)

func (t transform) String() string {
	s := ""
	switch t {
	case x0:
		s = "x0"
	case x90:
		s = "x90"
	case x180:
		s = "x180"
	case x270:
		s = "x270"
	case y0:
		s = "y0"
	case y90:
		s = "y90"
	case y180:
		s = "y180"
	case y270:
		s = "y270"
	case z0:
		s = "z0"
	case z90:
		s = "z90"
	case z180:
		s = "z180"
	case z270:
		s = "z270"
	}
	return s
}

type matrix struct {
	x, y, z transform
}

var alltransforms []matrix = []matrix{
	{x0, y180, z0}, {x270, y270, z270}, {x0, y270, z270}, {x270, y180, z90}, {x180, y0, z0},
	{x270, y0, z0}, {x180, y180, z0}, {x180, y180, z270}, {x0, y180, z270}, {x0, y180, z90},
	{x90, y180, z90}, {x180, y180, z90}, {x0, y90, z180}, {x90, y180, z0}, {x270, y180, z270},
	{x90, y180, z270}, {x0, y0, z0}, {x270, y180, z180}, {x90, y90, z180}, {x270, y180, z0},
	{x90, y270, z270}, {x180, y270, z270}, {x180, y90, z180}, {x270, y90, z180},
}

func (m matrix) String() string {
	return m.x.String() + "," + m.y.String() + "," + m.z.String()
}

type rotatedelta struct {
	transform matrix
	delta     xyz
}

type rotation struct {
	result    xyz
	transform matrix
}

type scanner struct {
	id        int
	points    []xyz
	position  xyz
	located   bool
	transform matrix
	delta     xyz
}

type beacon struct {
	pt            xyz
	delta         xyz
	relativePoint xyz
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
	fmt.Fprintln(w, "--2021 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	s := parseInput(in)
	return getBeacons(s, 12)
}

func part2(in input) output {
	return 0
}

func getBeacons(scanners []*scanner, minmatch int) int {
	m := map[xyz]beacon{}

	rel := scanners[0]
	if locateAllScanners(scanners, minmatch) {
		for _, s := range scanners {
			for _, p := range s.points {
				rp := p
				rp = rotateX(rp, s.transform.x)
				rp = rotateY(rp, s.transform.y)
				rp = rotateZ(rp, s.transform.z)

				rp = rel.position.add(s.delta)

				m[rp] = beacon{}
			}
		}
	}

	return len(m)
}

func locateAllScanners(scanners []*scanner, minmatch int) bool {
	done := false
	base := scanners[0]
	base.delta = xyz{0, 0, 0}
	base.located = true
	base.position = xyz{0, 0, 0}

	for i := 0; i < len(scanners[1:]); i++ {
		rel := scanners[i]

		sub := append(scanners[:i], scanners[i+1:]...)

		sc, rd := findNextScanner(sub[1:], rel, minmatch)

		if sc != nil {
			sc.located = true
			sc.transform = rd.transform
			sc.delta = rd.delta
		} else {
			fmt.Println("couldn't find match", rel.id)
			done = true
		}
	}

	return done
}

// func locateAllScanners(scanners []*scanner, minmatch int) bool {
// 	notLocated := make(map[int]*scanner)
// 	for _, s := range scanners {
// 		notLocated[s.id] = s
// 	}

// 	// we know scanner 0 is at 0,0,0
// 	lastLocated := notLocated[0]
// 	delete(notLocated, 0)
// 	for len(notLocated) > 0 {
// 		list := []*scanner{}
// 		for _, v := range notLocated {
// 			list = append(list, v)
// 		}

// 		sc, rd := findNextScanner(list, lastLocated, minmatch)
// 		if rd == nil {
// 			// problem. couldn't find a scanner
// 			log.Println("can't find next scanner", list)
// 			break
// 		}

// 		newpos := lastLocated.position.add(rd.delta)

// 		lastLocated = sc
// 		lastLocated.transform = rd.transform
// 		lastLocated.position = newpos

// 		// for i := 0; i < len(lastLocated.points); i++ {
// 		// 	r := lastLocated.points[i]
// 		// 	r = rotateX(r, rd.transform.x)
// 		// 	r = rotateY(r, rd.transform.y)
// 		// 	r = rotateZ(r, rd.transform.z)

// 		// 	lastLocated.points[i] = r.add(newpos)
// 		// }
// 		delete(notLocated, lastLocated.id)
// 	}

// 	return len(notLocated) == 0
// }

func findNextScanner(scanners []*scanner, relativeto *scanner, minmatch int) (*scanner, *rotatedelta) {
	var sc *scanner
	var rd *rotatedelta
	for _, s := range scanners {
		if s.id == relativeto.id {
			continue
		}
		fmt.Println("comparing", s.id, relativeto.id)
		rd = getRotationDelta(s, relativeto, minmatch)
		if rd != nil {
			sc = s
			break
		}
	}
	return sc, rd
}

func getRotationDelta(s, relativeto *scanner, minmatch int) *rotatedelta {
	deltam := map[xyz]*rotatedelta{}
	deltac := map[xyz]int{}
	for _, rel := range relativeto.points {
		for _, p := range s.points {
			rpts := rotate(p)
			for _, rpt := range rpts {
				delta := rel.delta(rpt.result)
				deltac[delta]++

				if _, ok := deltam[delta]; ok {
					fmt.Println("delta already added", delta, deltac[delta], p, "->", rel)
				}
				deltam[delta] = &rotatedelta{transform: rpt.transform, delta: delta}
			}
		}
	}
	// fmt.Println(deltac)
	var rd *rotatedelta
	max := 0
	for k, v := range deltac {
		if v >= minmatch && v > max {
			max = v
			fmt.Println("found minmatch", v, k, "relative to", relativeto.id, s.id)
			trd, ok := deltam[k]
			if ok {
				rd = trd
			}
		} else if v > 10 {
			fmt.Println("nearly found minmatch", v, k, "relative to", relativeto.id, s.id)

		}
	}
	return rd
}

// func getMatrices() []matrix {
// 	mx := []matrix{
// 		{x0, y180, z0}, {x270, y270, z270}, {x0, y270, z270}, {x270, y180, z90}, {x180, y0, z0},
// 		{x270, y0, z0}, {x180, y180, z0}, {x180, y180, z270}, {x0, y180, z270}, {x0, y180, z90},
// 		{x90, y180, z90}, {x180, y180, z90}, {x0, y90, z180}, {x90, y180, z0}, {x270, y180, z270},
// 		{x90, y180, z270}, {x0, y0, z0}, {x270, y180, z180}, {x90, y90, z180}, {x270, y180, z0},
// 		{x90, y270, z270}, {x180, y270, z270}, {x180, y90, z180}, {x270, y90, z180},
// 	}
// 	return mx
// }

func rotate2(p xyz) []xyz {
	x, y, z := p.x, p.y, p.z
	pts := []xyz{
		{x, y, z}, {x, z, y}, {y, x, z}, {y, z, x}, {z, x, y}, {z, y, x},
		{-x, y, z}, {-x, z, y}, {y, -x, z}, {y, z, -x}, {z, -x, y}, {z, y, -x},
		{x, -y, z}, {x, z, -y}, {-y, x, z}, {-y, z, x}, {z, x, -y}, {z, -y, x},
		{x, y, -z}, {x, -z, y}, {y, x, -z}, {y, -z, x}, {-z, x, y}, {-z, y, x},
	}
	return pts
}

func rotate(p xyz) []rotation {
	rots := []rotation{}
	for _, mx := range alltransforms {
		r := p
		r = rotateX(r, mx.x)
		r = rotateY(r, mx.y)
		r = rotateZ(r, mx.z)

		rot := rotation{result: r, transform: mx}
		rots = append(rots, rot)
	}

	if len(rots) != 24 {
		fmt.Println("rots", len(rots))
	}
	return rots
}

func rotateX(p xyz, t transform) xyz {
	np := p
	switch t {
	case x90:
		np.y, np.z = -np.z, np.y
	case x180:
		np.y, np.z = -np.y, -np.z
	case x270:
		np.y, np.z = np.z, -np.y
	}
	return np
}

func rotateY(p xyz, t transform) xyz {
	np := p
	switch t {
	case y90:
		np.x, np.z = np.z, -np.x
	case y180:
		np.x, np.z = -np.x, -np.z
	case y270:
		np.x, np.z = -np.z, np.x
	}
	return np
}

func rotateZ(p xyz, t transform) xyz {
	np := p
	switch t {
	case z90:
		np.x, np.y = np.y, -np.x
	case z180:
		np.x, np.y = -np.x, -np.y
	case z270:
		np.x, np.y = -np.y, np.x
	}
	return np
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
