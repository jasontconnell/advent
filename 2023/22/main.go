package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xyz struct {
	x, y, z int
}

type brick struct {
	name   string
	start  xyz
	end    xyz
	z      int
	height int
}

func (b brick) String() string {
	return fmt.Sprintf("%s (%d,%d,%d)~(%d,%d,%d)", b.name, b.start.x, b.start.y, b.start.z, b.end.x, b.end.y, b.end.z)
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 22 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	bricks := parseInput(in)
	fallen := fall(bricks)
	return determineDisintegrate(fallen)
}

func part2(in input) output {
	return 0
}

func determineDisintegrate(fallen map[int][]brick) int {
	maxz := 0
	levels := []int{}
	for k, v := range fallen {
		levels = append(levels, k)
		if k > maxz && len(v) > 0 {
			maxz = k
		}
	}

	sort.Ints(levels)
	blocks := 0
	for _, i := range levels {
		// fmt.Println(i, fallen[i])
		// if len(fallen[i]) > 0 {
		// 	fmt.Println(i, len(fallen[i]))
		// }
		blocks += len(fallen[i])
	}

	targets := 0
	for _, k := range levels {
		list := fallen[k]
		above, aok := fallen[k+1]

		fmt.Println("     ", k, "    ")

		if !aok || len(above) == 0 {
			targets += len(list)
		} else {
			for _, c := range list {
				counted := false
				for ai := 0; ai < len(above) && !counted; ai++ {
					a := above[ai]
					if c == a {
						continue
					}
					if collides(a, c) {
						fmt.Println(a, "collides with", c)
						cnt := 0
						for _, d := range list {
							if c == d {
								continue
							}
							if collides(a, d) {
								fmt.Println(a, "then collides with", d)
								cnt++
							}
						}
						// if a is supported by another block, we can remove this one
						if cnt > 0 {
							fmt.Println("counting target", c, a)
							targets++
							counted = true
						}
					}
				}
			}
		}
	}
	return targets
}

func fall(bricks []brick) map[int][]brick {
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].z < bricks[j].z
	})

	minz, maxz := math.MaxInt32, 0
	byz := make(map[int][]brick)
	for _, b := range bricks {
		byz[b.z] = append(byz[b.z], b)
		if b.z > maxz {
			maxz = b.z
		}
		if b.z < minz {
			minz = b.z
		}
	}

	fallen := make(map[int][]brick)

	level := 1
	done := false
	for !done {
		done = level > maxz
		if _, ok := fallen[level]; !ok {
			fallen[level] = append(fallen[level], byz[level]...)
		}
		if level == 1 {
			level++
			continue
		}

		cur := fallen[level]
		destlevel := level

		emptylevel := 0
		for j := len(cur) - 1; j >= 0; j-- {
			c := cur[j]
			cdone := false
			move := false
			for z := level - 1; z >= 1 && !cdone; z-- {
				below := fallen[z]
				if len(below) == 0 {
					emptylevel = z
					continue
				}
				collany := false
				for _, br := range below {
					if collides(c, br) {
						collany = true
					}
				}
				if collany {
					destlevel = z + 1
					move = destlevel != level
					cdone = true
					break
				}
			}

			if emptylevel != 0 && !move {
				move = true
				destlevel = emptylevel
			}

			if move && destlevel != level {
				fallen[destlevel] = append(fallen[destlevel], c)
				fallen[level] = append(fallen[level][:j], fallen[level][j+1:]...)

				if c.height > 0 {
					for z := destlevel + 1; z < destlevel+1+c.height; z++ {
						fallen[z] = append(fallen[z], c) // represent it in above levels
					}
				}
			}
		}

		level++
	}

	return fallen
}

func collides(b1, b2 brick) bool {
	xcollides := pointsCollide(b1.start.x, b1.end.x, b2.start.x, b2.end.x)
	ycollides := pointsCollide(b1.start.y, b1.end.y, b2.start.y, b2.end.y)

	return xcollides && ycollides
}

func pointsCollide(a1, a2, b1, b2 int) bool {
	// f1 := between(b1, a1, a2)
	// f2 := between(a1, b1, b2)
	// f3 := between(a2, b1, b2)
	// f4 := between(b2, a1, a2)
	ov := overlaps(a1, a2, b1, b2)

	// anyequal := a1 == b1 || a1 == b2 || a2 == b1 || a2 == b2

	return ov
	// return f1 || f2 || f3 || f4 || ov || anyequal
}

func between(a, b1, b2 int) bool {
	lo, hi := int(math.Min(float64(b1), float64(b2))), int(math.Min(float64(b1), float64(b2)))
	return lo < a && hi > a //|| (lo == a && hi > a) || (lo < a && hi == a)
}

func overlaps(a1, a2, b1, b2 int) bool {
	return (a1 <= b1 && a2 >= b2) || (b1 <= a1 && b2 >= a2)
}

func parseInput(in input) []brick {
	bricks := []brick{}
	a := int('A')
	zero := int('0')
	mod := 0
	nstart := 0
	for _, line := range in {
		lr := strings.Split(line, "~")
		left := strings.Split(lr[0], ",")
		lx, _ := strconv.Atoi(left[0])
		ly, _ := strconv.Atoi(left[1])
		lz, _ := strconv.Atoi(left[2])

		lcoord := xyz{lx, ly, lz}

		right := strings.Split(lr[1], ",")
		rx, _ := strconv.Atoi(right[0])
		ry, _ := strconv.Atoi(right[1])
		rz, _ := strconv.Atoi(right[2])

		rcoord := xyz{rx, ry, rz}

		h := int(math.Abs(float64(lz - rz)))
		z := int(math.Min(float64(lz), float64(rz)))

		bricks = append(bricks, brick{string(rune(nstart+a)) + string(rune(mod+zero)), lcoord, rcoord, z, h})
		if nstart+1 == 26 {
			mod = (mod + 1) % 10
		}
		nstart = (nstart + 1) % 26
	}
	return bricks
}
