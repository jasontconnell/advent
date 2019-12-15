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

var input = "12.txt"

type moon struct {
	x, y, z    int
	vx, vy, vz int
}

type energy int64

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	reg := regexp.MustCompile("([x-z])=(-?[0-9]*)")
	moons := []moon{}
	for scanner.Scan() {
		var txt = scanner.Text()

		m := moon{}
		if groups := reg.FindAllStringSubmatch(txt, -1); groups != nil && len(groups) > 1 {
			for _, g := range groups {
				val, err := strconv.Atoi(g[2])
				if err != nil {
					panic("parsing " + g[2])
				}
				switch g[1] {
				case "x":
					m.x = val
				case "y":
					m.y = val
				case "z":
					m.z = val
				}
			}
		}
		moons = append(moons, m)
	}

	mcp := make([]moon, len(moons))
	copy(mcp, moons)

	mcp = simulate(mcp, 1000)
	p1 := getEnergy(mcp)
	fmt.Println("Part 1: ", p1)

	p2 := part2(moons)
	fmt.Println("Part 2: ", p2)

	fmt.Println("Time", time.Since(startTime))
}

func pot(m moon) energy {
	return energy(math.Abs(float64(m.x)) + math.Abs(float64(m.y)) + math.Abs(float64(m.z)))
}

func kin(m moon) energy {
	return energy(math.Abs(float64(m.vx)) + math.Abs(float64(m.vy)) + math.Abs(float64(m.vz)))
}

func getEnergy(moons []moon) energy {
	var e energy
	for _, m := range moons {

		e += pot(m) * kin(m)
	}
	return e
}

type quad struct {
	a, b, c, d int
}

func part2(moons []moon) int64 {
	var i int64
	init := make([]moon, len(moons))
	copy(init, moons)

	var xr, yr, zr int64

	done := false

	for !done {
		updateV(moons)
		updateP(moons)
		i++

		xs, ys, zs := true, true, true
		xvs, yvs, zvs := true, true, true
		for mi, m := range moons {
			xs = xs && m.x == init[mi].x
			ys = ys && m.y == init[mi].y
			zs = zs && m.z == init[mi].z

			xvs = xvs && m.vx == init[mi].vx
			yvs = yvs && m.vy == init[mi].vy
			zvs = zvs && m.vz == init[mi].vz
		}

		if xs && xvs && xr == 0 {
			xr = i
		}

		if ys && yvs && yr == 0 {
			yr = i
		}

		if zs && zvs && zr == 0 {
			zr = i
		}

		done = xr != 0 && yr != 0 && zr != 0
	}

	return lcm(xr, yr, zr)
}

func updateV(moons []moon) {
	for mi := 0; mi < len(moons); mi++ {
		m := moons[mi]
		for mj := range moons {
			m2 := moons[mj]
			if mi == mj {
				continue
			}

			m.vx += comp(m.x, m2.x)
			m.vy += comp(m.y, m2.y)
			m.vz += comp(m.z, m2.z)
		}

		moons[mi] = m
	}
}

func updateP(moons []moon) {
	for mi := 0; mi < len(moons); mi++ {
		m := moons[mi]
		m.x += m.vx
		m.y += m.vy
		m.z += m.vz
		moons[mi] = m
	}
}

func simulate(moons []moon, steps int64) []moon {
	var i int64
	for i < steps {
		updateV(moons)
		updateP(moons)
		i++
	}
	return moons
}

func printMoons(moons []moon, step int) {
	fmt.Println("Step ", step)
	for i, m := range moons {
		fmt.Printf("%d pos:<x=%d, y=%d, z=%d> vel:<x=%d, y=%d, z=%d>\n", i+1, m.x, m.y, m.z, m.vx, m.vy, m.vz)
	}
}

func comp(a, b int) int {
	if a > b {
		return -1
	} else if a < b {
		return 1
	}
	return 0
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int64, integers ...int64) int64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
